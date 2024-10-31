package kafkaclient

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"im/internal/logger"
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
		"broker.address.family":    "v4",
		"enable.auto.offset.store": false,
		"session.timeout.ms":       6000,
	})
	if err != nil {
		logger.Fatalf("failed to dial leader:%v", err)
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
		logger.Errorf("consumer subscribe failed: %v", err)
		return err
	}
	return k.run()
}

func (k *KafkaConsumerClient) run() error {
	logger.Infof("kafka client run")

	for {
		select {
		case <-k.ctx.Done():
			logger.Infof("kafka client exit")
			return nil
		default:
			ev := k.client.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				if e.Headers != nil {

				}
				err := k.callBack(e.Value)
				if err != nil {
					logger.Errorf("call back ret error:%v", err)
				}
				_, err = k.client.StoreMessage(e)
				if err != nil {
					logger.Errorf("--- Error storing offset after message %s:\n, err:%v",
						e.TopicPartition, err)
				}
				break
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				logger.Errorf("%% Error: %v: %v\n", e.Code(), e)
			default:
			}
		}
	}
}
