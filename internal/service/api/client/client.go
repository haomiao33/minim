package client

import (
	"context"
	"im/internal/service/api/config"
	"im/pkg/kafkaclient"
	"im/pkg/redisclient"
)

var RedisClient *redisclient.RedisClient
var KafkaProductClient *kafkaclient.KafkaProductClient

func Init(ctx context.Context) {
	RedisClient = redisclient.NewRedisClient(ctx,
		config.Config.Redis.Addr,
		config.Config.Redis.Password)

	KafkaProductClient = kafkaclient.NewKafkaProductClient(
		ctx,
		config.Config.Kafka.Addresses)
}
