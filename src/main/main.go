package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

type Person struct {
	Name     string    `form:"name"`
	Address  string    `form:"address"`
	Age      *int      `form:"age"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

func main() {
	consumerGroup, err := sarama.NewConsumerGroup([]string{}, "", sarama.NewConfig())
	if err != nil {
		return
	}
	if err := consumerGroup.Consume(context.Background(), []string{}, EventHandler{}); err != nil {
		return
	}
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
