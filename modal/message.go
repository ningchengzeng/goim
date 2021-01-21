package modal

import (
	"context"
	"time"

	log "github.com/go-kratos/kratos/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageType int32

// 消息类型
const (
	BroadcastMessage MessageType = 0 // 广播消息
	RoomMessage      MessageType = 1 // 房间消息
	UserMessage      MessageType = 2 // 用户消息
	KeyMessage       MessageType = 3 // Key消息
)

// Message 消息
type Message struct {
	Id            string      `json:"id" bson:"_id"` // 消息ID
	UserId        string      //消息发送者
	ObjectUserId  string      //消息发送对象
	RoomId        string      //消息发送房间
	Type          MessageType //消息类型
	Keys          []string    //发送标签
	Body          string      //消息内容
	Violation     bool        //违规消息
	ViolationTime time.Time   //违规消息设置时间
	Del           bool        //删除消息
	DelTime       time.Time   //消息删除时间
	Send          bool        //消息是否已经发送
	SendTime      time.Time   // 消息发送时间
	Read          bool        //消息是否已读
	ReadTime      time.Time   //消息发送时间
}

// MessageCollectionName 消息表定义
const MessageCollectionName = "message"

// Insert 插入消息
func (msg *Message) Insert(database *mongo.Database) primitive.ObjectID {
	collection := database.Collection(MessageCollectionName)
	result, err := collection.InsertOne(context.TODO(), msg)

	if err != nil {
		log.Error("Message.Insert(%s) Error(%v)", msg, err)
	}
	doc := result.InsertedID.(primitive.ObjectID)
	return doc
}

// FindByBroadcast 获取广播消息
func (msg *Message) FindByBroadcast(database *mongo.Database) (results []Message) {
	collection := database.Collection(MessageCollectionName)

	opts := options.Find()
	opts.SetSort(bson.M{"SendTime": -1})

	dd, _ := time.ParseDuration("24h")
	cursor, err := collection.Find(context.TODO(), bson.M{
		"Type": BroadcastMessage,
		"Del":  false,
		"Send": true,
		"SendTime": bson.M{
			"$gt": time.Now().Add(dd * -90),
		},
	}, opts)
	if err != nil {
		log.Error("Message.FindByUserId(%s) Error(%v)", msg, err)
	}
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Error("Message.FindByUserId(%s) Error(%v)", msg, err)
	}
	return
}

// FindByUserId 根据用户获取消息
func (msg *Message) FindByUserId(database *mongo.Database) (results []Message) {
	collection := database.Collection(MessageCollectionName)

	opts := options.Find()
	opts.SetSort(bson.M{"SendTime": -1})

	dd, _ := time.ParseDuration("24h")
	cursor, err := collection.Find(context.TODO(), bson.M{
		"$or": bson.M{
			"UserId":       msg.UserId,
			"ObjectUserId": msg.ObjectUserId,
		},
		"Type": UserMessage,
		"Del":  false,
		"Send": true,
		"SendTime": bson.M{
			"$gt": time.Now().Add(dd * -90),
		},
	}, opts)
	if err != nil {
		log.Error("Message.FindByUserId(%s) Error(%v)", msg, err)
	}
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Error("Message.FindByUserId(%s) Error(%v)", msg, err)
	}
	return
}

// FindByKeys 根据用户获取消息
func (msg *Message) FindByKeys(database *mongo.Database) (results []Message) {
	collection := database.Collection(MessageCollectionName)

	opts := options.Find()
	opts.SetSort(bson.M{"SendTime": -1})

	dd, _ := time.ParseDuration("24h")
	cursor, err := collection.Find(context.TODO(), bson.M{
		"$in":  msg.Keys,
		"Type": KeyMessage,
		"Del":  false,
		"Send": true,
		"SendTime": bson.M{
			"$gt": time.Now().Add(dd * -90),
		},
	}, opts)
	if err != nil {
		log.Error("Message.FindByUserId(%s) Error(%v)", msg, err)
	}
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Error("Message.FindByUserId(%s) Error(%v)", msg, err)
	}
	return
}

// FindByRomeId 根据房间获取消息
func (msg *Message) FindByRomeId(database *mongo.Database) (results []Message) {
	collection := database.Collection(MessageCollectionName)

	opts := options.Find()
	opts.SetSort(bson.M{"SendTime": -1})

	dd, _ := time.ParseDuration("24h")
	cursor, err := collection.Find(context.TODO(), bson.M{
		"RoomId": msg.UserId,
		"Type":   RoomMessage,
		"Del":    false,
		"Send":   true,
		"SendTime": bson.M{
			"$gt": time.Now().Add(dd * -90),
		},
	}, opts)
	if err != nil {
		log.Error("Message.FindByRomeId(%s) Error(%v)", msg, err)
	}
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Error("Message.FindByRomeId(%s) Error(%v)", msg, err)
	}
	return
}

// SetViolation 设置为违规消息
func (msg *Message) SetViolation(database *mongo.Database) {
	objID, _ := primitive.ObjectIDFromHex(msg.Id)
	collection := database.Collection(MessageCollectionName)
	collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"Violation": true, "ViolationTime": time.Now()})
}

// SetRead 设置为已读
func (msg *Message) SetRead(database *mongo.Database) {
	objID, _ := primitive.ObjectIDFromHex(msg.Id)
	collection := database.Collection(MessageCollectionName)
	collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"Read": true, "ReadTime": time.Now()})
}

// SetSend 设置为已经发送
func (msg *Message) SetSend(database *mongo.Database) {
	objID, _ := primitive.ObjectIDFromHex(msg.Id)
	collection := database.Collection(MessageCollectionName)
	collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"Send": true, "SendTime": time.Now()})
}

// SetDel 设置为已删除
func (msg *Message) SetDel(database *mongo.Database) {
	objID, _ := primitive.ObjectIDFromHex(msg.Id)
	collection := database.Collection(MessageCollectionName)
	collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"Del": true, "DelTime": time.Now()})
}
