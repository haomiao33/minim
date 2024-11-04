# mini-im

## 1、说明：

### 1.1、项目介绍：
  简单的go写的im服务，流程简单清晰,大部分接口使用的是http，方便流程控制。login服务目前只是用来做服务端推送消息通知到客户端。本项目采用golang编写，分为login、api、msg-push、online等服务，这些服务都能集群部署和多个实例扩展。用户可以扩充其他协议和服务。
    
  目前打算golang编写im，支持单聊、群聊、推送； 然后客户端目前只打算做个flutter im chat 版本的就行。方便大家集成。

  文档放在doc里面了

  sql里面是数据库，自己创建一个就行

  测试客户端：目前是写到testclient目录里面的，nodejs的

  目前有如下服务（每个服务都可以多实例，方便用户量上来扩展和分布式，目前服务注册到consul里面的）：
    
    login: 用户登录服务，目前只用来接收服务端下发的消息（主要是消息通知），用户和im的websocket,使用 gnet，后续可以扩展到tcp、udp等，长连接都连接到login服务。不同服务交互使用的grpc，这里面没有写用户认证什么的，大家可以根据自己的业务需求来完成。
    
    api: 消息接口服务，采用fiber编写，主要处理客户端的接口请求：消息发送、消息同步、会话管理等等这些。
        
    msg-push: 消息推送服务，单聊消息推送
        
    online: 在线状态服务，用户在线状态放在这里，内部使用redis存放。login服务
## 1.2、项目特点：
    1、go编写
    2、使用gnet高并发库作为底层支持，目前支持websocket、后续可以很轻松扩展到gnet tcp等协议
    3、使用fiber作为api接口底层，支持高并发
    4、消息表目前支持分表，可以在创建的时候创建n个消息表im_msg_0是默认的
    5、协议简单，目前通信通道websocket只来做消息下发通知，其他的都通过api接口实现，稳定简单可靠，扩展性强
    6、使用kafka作为消息mq，高负载，后续添加群聊能很好支持
    7、有flutter客户端sdk和demo



## 1、3 项目部署：
查看 [部署方法](doc/部署.md) 文件

## 2、登陆
<img width="740" alt="image" src="https://github.com/user-attachments/assets/bd8024fa-f838-43ac-b4be-ee0066ed5a5e">



## 3、单聊消息：流程图

### 3.1、单聊流程图
<img width="1086" alt="image" src="https://github.com/user-attachments/assets/7233dda4-79d2-4357-8239-85d687a92334">




