package msg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/gofiber/fiber/v2"
	"im/internal/dao"
	"im/internal/db"
	"im/internal/logger"
	"im/internal/response"
	"im/internal/service/api/client"
	"im/internal/service/api/config"
	"im/internal/service/api/req"
	resp2 "im/internal/service/api/resp"
	"time"
)

type MsgHandler struct {
}

func NewMsgHandler(router fiber.Router) *MsgHandler {
	handler := &MsgHandler{}
	// 单聊消息发送
	router.Post("/msg/send", handler.SendMsg)
	// 单聊消息同步
	router.Post("/msg/sync", handler.SyncMsg)
	return handler
}

func (u *MsgHandler) SendMsg(c *fiber.Ctx) error {
	var msg req.ImMsgCommandReq
	if err := c.BodyParser(&msg); err != nil {
		return errors.New("参数错误")
	}

	getMsg, err := dao.MsgDao.GetMsg(db.Db, msg.MsgId)
	if err != nil {
		logger.Errorf("--- get msg failed:%s ---", msg.MsgId)
		return err
	}

	if getMsg != nil {
		logger.Warnf("--- msg already exist:%s ---", msg.MsgId)
		resp := resp2.MsgSendResp{
			Id:             msg.MsgId,
			Sequence:       getMsg.Sequence,
			Status:         getMsg.Status,
			ConversationId: getMsg.ConversationID,
			CreatedTime:    getMsg.CreatedTime,
			UpdatedTime:    getMsg.UpdatedTime,
			DeletedTime:    getMsg.DeletedTime,
			RevokedTime:    getMsg.RevokedTime,
		}
		return c.JSON(response.Success(resp))
	}

	bigId := msg.FromId
	if msg.ToId > bigId {
		bigId = msg.ToId
	}
	smallId := msg.ToId
	if smallId > msg.FromId {
		smallId = msg.FromId
	}

	//lock
	key := fmt.Sprintf("lock:msg:%d:%d:%d", msg.MsgType, bigId, smallId)
	mutex := client.RedisClient.RedisLock.NewMutex(key, redsync.WithExpiry(5*time.Second))
	if err := mutex.Lock(); err != nil {
		logger.Errorf("--- lock failed:%s ---", msg.MsgId)
	}

	//事务开启
	tx := db.Db.Begin()

	defer func() {
		if err != nil {
			logger.Errorf("--- rollback:%s ---", msg.MsgId)
			tx.Rollback()
		}
		if ok, err := mutex.Unlock(); !ok || err != nil {
			logger.Errorf("--- unlock failed:%s ---", msg.MsgId)
		}
	}()

	//查询会话关系
	conversation, err := dao.ConversationDao.GetConversation(tx, msg.ChatType, bigId, smallId)
	if err != nil {
		logger.Errorf("--- get conversation failed:%s ---", msg.MsgId)
		return err
	}
	if conversation == nil {
		conversation, err = dao.ConversationDao.AddConversation(tx, msg.ChatType, bigId, smallId)
		if err != nil {
			logger.Errorf("--- create conversation failed:%s ---", msg.MsgId)
			return err
		}
	}

	//获取序号
	seqKey := fmt.Sprintf("sequence:%d", conversation.ID)
	sequence, err := client.RedisClient.GetSequence(seqKey, conversation.Sequence)
	if err != nil {
		logger.Errorf("--- get sequence failed:%s ---", msg.MsgId)
		return err
	}
	//添加消息
	addMsg, err := dao.MsgDao.AddMsg(tx, conversation.ID, sequence, msg.MsgId,
		int(msg.ChatType), int(msg.MsgType),
		msg.FromId, msg.ToId, msg.Content, msg.Ts, 0)
	if err != nil {
		logger.Errorf("--- add msg failed:%s ---", msg.MsgId)
		return err
	}

	//更新会话序号
	err = dao.ConversationDao.UpdateConversation(tx, conversation.ID, sequence)
	if err != nil {
		logger.Errorf("--- update conversation failed:%s ---", msg.MsgId)
		return err
	}

	//发送者session
	senderSession, err := dao.RecentSessionDao.Get(tx, msg.ChatType, msg.FromId, msg.ToId)
	if err != nil {
		logger.Errorf("--- get recent session failed:%s ---", msg.MsgId)
		return err
	}
	if senderSession == nil {
		err := dao.RecentSessionDao.Add(tx,
			conversation.ID,
			msg.ChatType, msg.FromId, msg.ToId,
			addMsg.ID, addMsg.Content, time.UnixMilli(msg.Ts))
		if err != nil {
			logger.Errorf("--- add recent session failed:%s ---", msg.MsgId)
			return err
		}
	} else {
		err := dao.RecentSessionDao.Update(tx,
			msg.ChatType, msg.FromId, msg.ToId,
			addMsg.ID, addMsg.Content, time.UnixMilli(msg.Ts))
		if err != nil {
			logger.Errorf("--- update recent session failed:%s ---", msg.MsgId)
			return err
		}
	}

	//添加接收者session
	receiverSession, err := dao.RecentSessionDao.Get(tx, msg.ChatType, msg.ToId, msg.FromId)
	if err != nil {
		logger.Errorf("--- get recent session failed:%s ---", msg.MsgId)
		return err
	}

	if receiverSession == nil {
		err := dao.RecentSessionDao.Add(tx,
			conversation.ID,
			msg.ChatType, msg.ToId, msg.FromId,
			addMsg.ID, addMsg.Content, time.UnixMilli(msg.Ts))
		if err != nil {
			logger.Errorf("--- add recent session failed:%s ---", msg.MsgId)
			return err
		}
	} else {
		err := dao.RecentSessionDao.Update(tx,
			msg.ChatType, msg.ToId, msg.FromId,
			addMsg.ID, addMsg.Content, time.UnixMilli(msg.Ts))
		if err != nil {
			logger.Errorf("--- update recent session failed:%s ---", msg.MsgId)
			return err
		}
	}

	logger.Infof("send msg success  msgId:%s,chatType:%d, msgType:%d, sequence:%d",
		msg.MsgId, msg.ChatType, msg.MsgType, sequence)
	//提交事务
	tx.Commit()

	//发送消息到kafka，进行后续推送
	//相同会话发送到同一个分区
	partition := int32(conversation.ID % int64(config.Config.Kafka.MsgPartitionCount))

	marshal, _ := json.Marshal(addMsg)
	err = client.KafkaProductClient.ProductMessage(
		config.Config.Kafka.MsgTopic,
		partition,
		marshal)
	if err != nil {
		logger.Errorf("--- send message to kafka error: %v， msgId:%s ---", err, msg.MsgId)
		//这里后面有其他兜底（客户端主动pull），这里忽略
		err = nil
	} else {
		logger.Infof("add to push mq  success:%s", msg.MsgId)
	}

	resp := resp2.MsgSendResp{
		Id:             msg.MsgId,
		Sequence:       sequence,
		ConversationId: conversation.ID,
		Status:         addMsg.Status,
		CreatedTime:    addMsg.CreatedTime,
		UpdatedTime:    addMsg.UpdatedTime,
		DeletedTime:    addMsg.DeletedTime,
		RevokedTime:    addMsg.RevokedTime,
	}
	return c.JSON(response.Success(resp))

}

func (u *MsgHandler) SyncMsg(c *fiber.Ctx) error {
	var req req.ImMsgSyncCommandReq
	if err := c.BodyParser(&req); err != nil {
		return errors.New("参数错误")
	}

	list, err := dao.MsgDao.GetMsgList(db.Db, req.ConversationId, req.Sequence)
	if err != nil {
		logger.Errorf("--- get msg list failed:%s ---", req.ConversationId)
		return err
	}

	resp := resp2.ImMsgSyncCommandResp{
		Items:     list,
		OtherInfo: nil,
	}
	if req.OtherId > 0 {
		//获取用户信息附带返回
		info, err := dao.UserDao.GetUserByFiled(db.Db, req.OtherId, []string{
			"user_id", "avatar", "nick_name", "user_type",
		})
		if err != nil {
			logger.Errorf("--- get user failed:%s ---", req.OtherId)
			return err
		}
		resp.OtherInfo = &resp2.ImUserInfoResp{
			UserID:   info.UserID,
			NickName: info.NickName,
			Avatar:   info.Avatar,
			UserType: info.UserType,
		}
	}

	return c.JSON(response.Success(resp))
}
