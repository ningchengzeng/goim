package modal

import "time"

// UserFriend is 用户好友
type UserFriend struct {
	UserId      string    //好友编号
	NickName    string    //好友备注
	Tag         []string  //好友标签
	Top         bool      //是否置顶
	Remove      bool      //是否移除
	Silence     bool      //是否浸没
	SilenceTime bool      //浸没时间
	JoinTime    time.Time //加入时间
	RemoveTime  time.Time //移除时间
}

// User is 用户信息
type User struct {
	Code       string       //好友编号
	Tag        []string     //用户标签
	Name       string       //好友名称
	Phone      string       //好友手机号码
	Icon       string       //好友头像
	Del        bool         //是否删除
	Friends    []UserFriend //好友名单
	Frozen     bool         //冻结好友
	FrozenTime time.Time    //冻结时间
	CreateTime time.Time    //创建时间
	UpdateTime time.Time    //修改时间
}
