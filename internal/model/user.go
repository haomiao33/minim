package model

// User 用户表
type User struct {
	UserID      int64  `gorm:"primaryKey;column:user_id" json:"userId"` // 用户Id
	Status      int    `gorm:"column:status" json:"status"`             // 用户状态:0=正常；1=锁定；
	UniqueID    string `gorm:"column:unique_id" json:"uniqueId"`        // 唯一id
	UserRegion  int    `gorm:"column:user_region" json:"userRegion"`    // 用户区域；0=中国；1=美国；2=日本
	UserName    string `gorm:"column:user_name" json:"userName"`        // 用户名字
	NickName    string `gorm:"column:nick_name" json:"nickName"`        // 用户昵称
	Phone       string `gorm:"column:phone" json:"phone"`               // 手机号
	Anonymous   int    `gorm:"column:anonymous" json:"anonymous"`       // 是否匿名，1=预注册；2=真正注册（手机号等）
	Sex         int    `gorm:"column:sex" json:"sex"`                   // 性别;0=未知；1=男；2=女
	Avatar      string `gorm:"column:avatar" json:"avatar"`             // 头像地址
	Domain      string `gorm:"column:domain" json:"domain"`             // 所属域名
	Password    string `gorm:"column:password" json:"password"`         // 密码
	UserType    int    `gorm:"column:user_type" json:"userType"`        // 用户类型:user=0;admin=10000;机构=300
	Deleted     int    `gorm:"column:deleted" json:"deleted"`           // 是否被删除；0=否；1=是
	UserFlg     string `gorm:"column:user_flg" json:"userFlg"`          // 暂无用处
	Birthday    MyTime `gorm:"column:birthday" json:"birthday"`         // 生日
	Description string `gorm:"column:description" json:"description"`   // 用户描述
	IP          string `gorm:"column:ip" json:"ip"`                     // ip
	Location    string `gorm:"column:location" json:"location"`         // 家庭位置
	Career      string `gorm:"column:career" json:"career"`             // 职业
	School      string `gorm:"column:school" json:"school"`             // 学校
	Email       string `gorm:"column:email" json:"email"`               // 电子邮件
	Qq          string `gorm:"column:qq" json:"qq"`                     // qq
	Weixin      string `gorm:"column:weixin" json:"weixin"`             // weixin
	Weibo       string `gorm:"column:weibo" json:"weibo"`               // weibo
	Stars       int    `gorm:"column:stars" json:"stars"`               // 被点赞数
	MyFollowers int    `gorm:"column:my_followers" json:"myFollowers"`  // 我的关注数
	Followers   int    `gorm:"column:followers" json:"followers"`       // 粉丝数
	RegType     string `gorm:"column:reg_type" json:"regType"`          // 注册来源
	CreateTime  MyTime `gorm:"column:create_time" json:"createTime"`    // 创建时间
	UpdateBy    int64  `gorm:"column:update_by" json:"updateBy"`        // 修改者
	UpdateTime  MyTime `gorm:"column:update_time" json:"updateTime"`    // 修改时间
	CreateBy    int64  `gorm:"column:create_by" json:"createBy"`        // 创建者
}

// TableName get sql table name.获取数据库表名
func (m *User) TableName() string {
	return "user"
}
