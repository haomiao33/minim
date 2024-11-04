package sharding

import "fmt"

// 目前分表数量
const shardingCount = 2

func GetTableName(tableName string, conversationId int64) string {
	return fmt.Sprintf("%s_%d", tableName, conversationId%shardingCount)
}
