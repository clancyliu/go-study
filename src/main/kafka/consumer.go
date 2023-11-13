package kafka

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

const (
	ADDR    = ""
	VERSION = ""
)

func NewConsumerGroup(group string) sarama.ConsumerGroup {
	cfg := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(VERSION)
	if err != nil {
		log.Fatal("NewConsumerGroup Parse kafka version failed", err.Error())
		return nil
	}

	cfg.Version = version
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Offsets.Retry.Max = 3
	cfg.Consumer.Offsets.AutoCommit.Enable = true              // 开启自动提交，需要手动调用MarkMessage才有效
	cfg.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 间隔
	client, err := sarama.NewConsumerGroup([]string{""}, group, cfg)
	if err != nil {
		log.Fatal("NewConsumerGroup failed", err.Error())
	}
	return client
}

type EventHandler struct{}

func (e EventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var data sarama.Message
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			return errors.New("failed to unmarshal message err is " + err.Error())
		}
		// 操作数据，改用打印
		log.Print("consumerClaim data is ")

		// 处理消息成功后标记为处理, 然后会自动提交
		session.MarkMessage(msg, "")
	}
	return nil
}

func (EventHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (EventHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
