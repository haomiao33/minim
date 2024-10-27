package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// 全局 logger 实例
var log *zap.SugaredLogger

// Init 初始化全局日志
func Init(pathFile string, level string) error {
	var err error
	// 创建文件输出
	file, err := os.OpenFile(pathFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// 创建日志编码器为文本格式
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	})

	// 创建日志级别
	atomicLevel := zap.NewAtomicLevel()
	if level == "debug" {
		atomicLevel.SetLevel(zap.DebugLevel)
	} else if level == "info" {
		atomicLevel.SetLevel(zap.InfoLevel)
	} else if level == "error" {
		atomicLevel.SetLevel(zap.ErrorLevel)
	} else if level == "warn" {
		atomicLevel.SetLevel(zap.WarnLevel)
	} else if level == "fatal" {
		atomicLevel.SetLevel(zap.FatalLevel)
	}

	// 创建多个日志核心（同时写入文件和控制台）
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(file), atomicLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), atomicLevel),
	)

	// 创建 Logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	log = logger.Sugar()
	return nil
}

// 日志记录方法
func Debugf(msg string, data ...interface{}) {
	log.Debugf(msg, data...)
}

func Errorf(msg string, data ...interface{}) {
	log.Errorf(msg, data...)
}

func Infof(msg string, data ...interface{}) {
	log.Infof(msg, data...)
}

func Warnf(msg string, data ...interface{}) {
	log.Warnf(msg, data...)
}

func Fatalf(msg string, data ...interface{}) {
	log.Fatalf(msg, data...)
}

// 日志记录方法
func Debug(msg string) {
	log.Debug(msg)
}

func Error(msg string) {
	log.Error(msg)
}

func Info(msg string) {
	log.Info(msg)
}

func Warn(msg string) {
	log.Warn(msg)
}

func Fatal(msg string) {
	log.Fatal(msg)
}

// Sync 确保日志写入
func Sync() error {
	return log.Sync()
}
