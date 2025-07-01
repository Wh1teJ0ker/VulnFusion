// log/logger.go

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var logger *zap.SugaredLogger

// InitLogger 初始化日志系统，配置 zap + lumberjack 实现日志切割与格式化输出
func InitLogger(logFile string, level string) error {
	core, err := buildLoggerCore(logFile, level)
	if err != nil {
		return err
	}
	// 创建 sugared logger，添加调用信息
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	return nil
}

// buildLoggerCore 构建 zapcore.Core，包含日志编码器、日志等级与输出目的地
func buildLoggerCore(logFile string, level string) (zapcore.Core, error) {
	atomicLevel := zap.NewAtomicLevel()
	if err := atomicLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}

	writerSyncer := getLogWriter(logFile) // 获取文件写入器（支持切割）
	encoder := getConsoleEncoder()        // 获取控制台格式编码器

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writerSyncer),
		atomicLevel,
	)
	return core, nil
}

// getLogWriter 返回基于 lumberjack 的写入器，用于日志自动切割
func getLogWriter(logFile string) zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile, // 日志文件路径
		MaxSize:    10,      // 每个日志文件最大 10MB
		MaxBackups: 7,       // 最多保留 7 个备份
		MaxAge:     30,      // 日志保留 30 天
		Compress:   true,    // 启用 gzip 压缩
	}
	return zapcore.AddSync(lumberjackLogger)
}

// getConsoleEncoder 返回 zap 使用的控制台编码器（可读性强）
func getConsoleEncoder() zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",       // 时间字段
		LevelKey:       "level",      // 日志级别字段
		NameKey:        "logger",     // logger 名称
		CallerKey:      "caller",     // 调用位置字段
		MessageKey:     "msg",        // 日志内容字段
		StacktraceKey:  "stacktrace", // 堆栈信息字段
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 彩色大写等级
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,   // 时长格式
		EncodeCaller:   zapcore.ShortCallerEncoder,       // 短文件路径
	}
	return zapcore.NewConsoleEncoder(encoderCfg)
}

// Debug 输出 DEBUG 等级日志
func Debug(format string, args ...interface{}) { logger.Debugf(format, args...) }

// Info 输出 INFO 等级日志
func Info(format string, args ...interface{}) { logger.Infof(format, args...) }

// Warn 输出 WARN 等级日志
func Warn(format string, args ...interface{}) { logger.Warnf(format, args...) }

// Error 输出 ERROR 等级日志
func Error(format string, args ...interface{}) { logger.Errorf(format, args...) }

// Fatal 输出 FATAL 等级日志并退出程序
func Fatal(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
	os.Exit(1)
}
