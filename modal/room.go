package modal

import (
	"context"
	"time"

	log "github.com/go-kratos/kratos/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RoomUser 房间用户
type RoomUser struct {
	UserId   string    //用户ID
	NickName string    //用户昵称
	Top      bool      //房间是否置顶
	Join     bool      //是否加入房间
	JoinTime time.Time //加入房间时间
	OutTime  time.Time //离开房间时间
}

// RoomUser 房间
type Room struct {
	Icon            string     //房间图标
	Name            string     //房间名称
	Type            string     //房间类型
	User            []RoomUser //房间用户
	Dissolution     bool       //房间是否解散
	LeadUserId      string     //房间管理者
	DissolutionTime time.Time  //房间解散时间
	CreateTime      time.Time  //创建时间
	UpdateTime      time.Time  //更新时间
}

// RoomCollectionName 消息表定义
const RoomCollectionName = "message"

// Insert 插入消息
func (msg *Room) Insert(database *mongo.Database) primitive.ObjectID {
	collection := database.Collection(RoomCollectionName)
	msg.CreateTime = time.Now()
	result, err := collection.InsertOne(context.TODO(), msg)

	if err != nil {
		log.Error("Room.Insert(%s) Error(%v)", msg, err)
	}
	doc := result.InsertedID.(primitive.ObjectID)
	return doc
}
