package redisclient

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

const luaScript = `
        local key = KEYS[1]
        local db_value = tonumber(ARGV[1])  -- 处理传入的数据库值

        -- 检查 Redis 中是否存在序号
        local current_sequence = redis.call('GET', key)

        -- 如果 Redis 中没有序号，用数据库的值初始化
        if not current_sequence then
            redis.call('SET', key, db_value)
            current_sequence = db_value
        else
            current_sequence = tonumber(current_sequence)
        end

        -- 序号自增
        current_sequence = redis.call('INCR', key)

        return current_sequence
    `

type RedisClient struct {
	ctx                context.Context
	Client             *redis.Client
	RedisLock          *redsync.Redsync
	sequenceScriptHash string
}

func NewRedisClient(ctx context.Context, addr string, password string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	//加载lua脚本
	sequenceScriptHash := rdb.ScriptLoad(ctx, luaScript)

	pool := goredis.NewPool(rdb)
	redisLock := redsync.New(pool)
	return &RedisClient{
		ctx:                ctx,
		Client:             rdb,
		RedisLock:          redisLock,
		sequenceScriptHash: sequenceScriptHash.Val(),
	}
}

func (r *RedisClient) GetSequence(key string, seq int64) (int64, error) {
	res := r.Client.EvalSha(r.ctx, r.sequenceScriptHash, []string{key}, seq)
	if res.Err() != nil {
		return -1, res.Err()
	}
	return res.Val().(int64), nil
}
