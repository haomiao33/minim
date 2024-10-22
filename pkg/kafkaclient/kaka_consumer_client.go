package kafkaclient

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"log"
)

type KafkaConsumerClient struct {
	client   *kafka.Consumer
	ctx      context.Context
	callBack func(data []byte) error
}

func NewKafkaConsumerClient(ctx context.Context, address string, groupId string) *KafkaConsumerClient {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        address,
		"group.id":                 groupId,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
		"session.timeout.ms":       6000,
	})
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	return &KafkaConsumerClient{
		client: c,
		ctx:    ctx,
	}
}

func (k *KafkaConsumerClient) StartConsume(topics []string, callBack func(data []byte) error) error {
	k.callBack = callBack

	err := k.client.SubscribeTopics(topics, nil)
	if err != nil {
		logging.Errorf("consumer subscribe failed: %v", err)
		return err
	}
	return k.run()
}

func (k *KafkaConsumerClient) run() error {
	logging.Infof("kafka client run")

	for {
		select {
		case <-k.ctx.Done():
			logging.Infof("kafka client exit")
			return nil
		default:
			ev := k.client.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
				err := k.callBack(e.Value)
				if err != nil {
					logging.Errorf("call back ret error:%v", err)
				}
				_, err = k.client.StoreMessage(e)
				if err != nil {
					logging.Errorf("--- Error storing offset after message %s:\n, err:%v",
						e.TopicPartition, err)
				}
				break
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				logging.Errorf("%% Error: %v: %v\n", e.Code(), e)
			default:
			}
		}
	}
}
