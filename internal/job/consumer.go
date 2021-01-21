package job

import (
	"context"
	"time"

	log "github.com/go-kratos/kratos/pkg/log"
	"github.com/golang/protobuf/proto"
	pb "github.com/ningchengzeng/goim/api/logic"
	"github.com/ningchengzeng/goim/internal/job/conf"
	"github.com/nsqio/go-nsq"
)

// NsqHandler 结构
type NsqHandler struct {
	consumer *nsq.Consumer
	topic    string
	channel  string
	j        *Job
}

// HandleMessage 处理消息
func (nh *NsqHandler) HandleMessage(msg *nsq.Message) error {
	pushMsg := new(pb.PushMsg)
	if err := proto.Unmarshal(msg.Body, pushMsg); err != nil {
		log.Error("proto.Unmarshal(%v) error(%v)", msg, err)
		return err
	}
	if err := nh.j.push(context.Background(), pushMsg); err != nil {
		log.Error("j.push(%v) error(%v)", pushMsg, err)
		return err
	}
	log.Info("consume: %s/%d/%d\t%s\t%+v", nh.topic, nh.channel, msg.NSQDAddress, msg.ID, pushMsg)
	msg.Finish()
	return nil
}

// Close 关闭接受
func (nh *NsqHandler) Close() error {
	if nh.consumer != nil {
		nh.consumer.Stop()
	}
	return nil
}

// NewNsqConsumer 创建 nsq consumer
func NewNsqConsumer(cnsq *conf.Nsq, job *Job) (*NsqHandler, error) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = 3 * time.Second
	c, err := nsq.NewConsumer(cnsq.Topic, cnsq.Channel, cfg)
	if err != nil {
		log.Error("init Consumer NewConsumer error:", err)
		return nil, err
	}

	handler := &NsqHandler{consumer: c, j: job}
	c.AddHandler(handler)
	err = c.ConnectToNSQLookupds(cnsq.Address)
	if err != nil {
		log.Error("init Consumer ConnectToNSQLookupd error:", err)
		return nil, err
	}
	return handler, nil
}
