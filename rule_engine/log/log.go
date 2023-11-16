package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
)

var (
	lock    sync.Mutex
	_logger *zap.Logger
)

const (
	// FormatText format log text
	FormatText = "text"
	// FormatJSON format log json
	FormatJSON           = "json"
	DefaultLogTimeLayout = "2006-01-02 15:04:05.000"
)

// type Level uint

// 日志配置
type Config struct {
	LogPath    string
	LogLevel   string
	Compress   bool
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Format     string
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func newLogWriter(logPath string, maxSize, maxBackups, maxAge int, compress bool) io.Writer {
	if logPath == "" || logPath == "-" {
		return os.Stdout
	}
	return &lumberjack.Logger{
		Filename:   logPath,    // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}
}

func newZapEncoder() zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "line",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		//EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeTime:     zapcore.TimeEncoderOfLayout(DefaultLogTimeLayout),
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return encoderConfig
}

func newLoggerCore(cfg *Config) zapcore.Core {
	hook := newLogWriter(cfg.LogPath, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Compress)

	encoderConfig := newZapEncoder()

	atomLevel := zap.NewAtomicLevelAt(getZapLevel(cfg.LogLevel))

	var encoder zapcore.Encoder
	if cfg.Format == FormatJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(hook)),
		atomLevel,
	)
	return core
}

func newLoggerOptions() []zap.Option {
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	callerskip := zap.AddCallerSkip(1)
	// 开发者
	development := zap.Development()
	options := []zap.Option{
		caller,
		callerskip,
		development,
		zap.Fields(zap.String("project_name", "rule_engine")),
	}
	return options
}

// fill default config
func (c *Config) fillWithDefault() {
	if c.MaxSize <= 0 {
		c.MaxSize = 20
	}
	if c.MaxAge <= 0 {
		c.MaxAge = 7
	}
	if c.MaxBackups <= 0 {
		c.MaxBackups = 7
	}
	if c.LogLevel == "" {
		c.LogLevel = "debug"
	}
	if c.Format == "" {
		c.Format = FormatText
	}
}

// InitLog config
func Init(cfg *Config) {
	cfg.fillWithDefault()
	core := newLoggerCore(cfg)
	zapOpts := newLoggerOptions()
	_logger = zap.New(core, zapOpts...)
}

func Logger() *zap.Logger {
	return _logger
}

// Debug output log
func Debug(msg string, fields ...zap.Field) {
	_logger.Debug(msg, fields...)
}

// Debugf logs a debug message.
func Debugf(msg string, args ...interface{}) {
	lock.Lock()
	defer lock.Unlock()

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	_logger.Debug(msg)
}

// Info output log
func Info(msg string, fields ...zap.Field) {
	_logger.Info(msg, fields...)
}

// Infof logs an info message.
func Infof(msg string, args ...interface{}) {
	lock.Lock()
	defer lock.Unlock()

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	_logger.Info(msg)
}

// Warn output log
func Warn(msg string, fields ...zap.Field) {
	_logger.Warn(msg, fields...)
}

// Warnf logs an error message.
func Warnf(msg string, args ...interface{}) {
	lock.Lock()
	defer lock.Unlock()

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	_logger.Warn(msg)
}

// Error output log
func Error(msg string, fields ...zap.Field) {
	_logger.Error(msg, fields...)
}

// Errorf logs an error message.
func Errorf(msg string, args ...interface{}) {
	lock.Lock()
	defer lock.Unlock()

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	_logger.Error(msg)
}

// Panic output panic
func Panic(msg string, fields ...zap.Field) {
	_logger.Panic(msg, fields...)
}

// Fatal output log
func Fatal(msg string, fields ...zap.Field) {
	_logger.Fatal(msg, fields...)
}
