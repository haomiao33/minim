package logger

import (
	"context"
	"gorm.io/gorm/logger"
	"time"
)

// GormLogger 是 GORM 兼容的日志实现
type GormLogger struct{}

// NewGormLogger 创建 GORM 兼容的 logger
func NewGormLogger() *GormLogger {
	return &GormLogger{}
}

// LogMode 设置日志级别
func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

// Info 记录信息级别的日志
func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Infof(msg, data...)
}

// Warn 记录警告级别的日志
func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Warnf(msg, data...)
}

// Error 记录错误级别的日志
func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Errorf(msg, data...)
}

// Debug 记录错误级别的日志
func (g *GormLogger) Debug(ctx context.Context, msg string, data ...interface{}) {
	log.Debugf(msg, data...)
}

// Info 记录信息级别的日志
func (g *GormLogger) Infof(msg string, data ...interface{}) {
	log.Infof(msg, data...)
}

// Warn 记录警告级别的日志
func (g *GormLogger) Warnf(msg string, data ...interface{}) {
	log.Warnf(msg, data...)
}

// Error 记录错误级别的日志
func (g *GormLogger) Errorf(msg string, data ...interface{}) {
	log.Errorf(msg, data...)
}

// Trace 记录 SQL 查询信息
func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	if err != nil {
		log.Errorf("SQL Error: %s, Error: %v", sql, err)
	} else {
		log.Infof("SQL Trace: %s, Rows: %d", sql, rows)
	}
}

// Printf 用于记录格式化日志
func (g *GormLogger) Printf(format string, v ...interface{}) {
	log.Infof(format, v...)
}
