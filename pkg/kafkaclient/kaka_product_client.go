package kafkaclient

import (
	"context"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"im/internal/logger"
	"time"
)

type KafkaProductClient struct {
	client   *kafka.Producer
	ctx      context.Context
	callBack func(data string) error
}

func NewKafkaProductClient(ctx context.Context,
	address string) *KafkaProductClient {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": address,
		//"security.protocol": "SSL",
	})
	if err != nil {
		logger.Fatalf("Failed to create producer: %s\n", err)
	}

	return &KafkaProductClient{
		client: p,
		ctx:    ctx,
	}
}

func (k *KafkaProductClient) ProductMessage(topic string, partition int32, data []byte) error {
	i := 0
	for i < 3 {
		err := k.client.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
			Value:          data,
		}, nil)

		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrQueueFull {
				// Producer queue is full, wait 1s for messages
				// to be delivered then try again.
				logger.Warnf("send message queue is full, retry :%d", i)
				time.Sleep(time.Second)
				i++
				continue
			}
			logger.Errorf("Failed to produce message: %v\n, data:%s", err, string(data))
			return err
		}
		return nil
	}
	logger.Errorf("send to kafka max try error， data:%s", string(data))
	return errors.New("send to kafka max try error")
}
