package logger

// GNetLogger 是 gnet 兼容的日志实现
type GNetLogger struct{}

// NewGNetLogger 创建 gnet 兼容的 logger
func NewGNetLogger() *GNetLogger {
	return &GNetLogger{}
}

// Info 记录信息级别的日志
func (g *GNetLogger) Infof(fmt string, data ...interface{}) {
	log.Infof(fmt, data...)
}

func (g *GNetLogger) Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func (g *GNetLogger) Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func (g *GNetLogger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func (g *GNetLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
