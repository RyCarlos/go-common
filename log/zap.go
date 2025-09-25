// Package log -----------------------------
// @file      : zap.go
// @author    : Carlos
// @contact   : 534994749@qq.com
// @time      : 2025/6/17 16:17
// -------------------------------------------
package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

const (
	LevelFatal = iota
	LevelPanic
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelDebugWithSQL
)

var (
	logger      Logger
	logLevelMap = map[int]zapcore.Level{
		LevelDebugWithSQL: zapcore.DebugLevel,
		LevelDebug:        zapcore.DebugLevel,
		LevelInfo:         zapcore.InfoLevel,
		LevelWarn:         zapcore.WarnLevel,
		LevelError:        zapcore.ErrorLevel,
		LevelPanic:        zapcore.PanicLevel,
		LevelFatal:        zapcore.FatalLevel,
	}
)

const (
	callDepth int    = 2
	logPath   string = "./logs/"
)

type Config struct {
	Level      int    `mapstructure:"level"`
	IsJson     bool   `mapstructure:"isJson"`
	IsStdout   bool   `mapstructure:"isStdout"`
	LogPath    string `mapstructure:"logPath"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge"`
	Compress   bool   `mapstructure:"compress"`
}

type ZapLogger struct {
	logLevel zapcore.Level
	zap      *zap.SugaredLogger
}

func init() {
	InitLoggerFromConfig(&Config{
		Level:      LevelDebugWithSQL,
		IsJson:     false,
		IsStdout:   false,
		LogPath:    logPath,
		MaxSize:    0,
		MaxBackups: 0,
		MaxAge:     0,
		Compress:   false,
	})
}

func InitLoggerFromConfig(config *Config) error {
	_, err := NewZapLogger(config)
	if err != nil {
		return err
	}
	return nil
}

func Debug(msg string, keysAndValues ...any) {
	logger.Debug(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...any) {
	logger.Info(msg, keysAndValues...)
}

func Warn(msg string, err error, keysAndValues ...any) {
	logger.Warn(msg, err, keysAndValues...)
}

func Error(msg string, err error, keysAndValues ...any) {
	logger.Error(msg, err, keysAndValues...)
}

func Panic(msg string, err error, keysAndValues ...any) {
	logger.Panic(msg, err, keysAndValues...)
}

func NewZapLogger(config *Config) (*ZapLogger, error) {
	zConfig := zap.Config{
		Level:    zap.NewAtomicLevelAt(logLevelMap[config.Level]),
		Encoding: getEncoding(config.IsJson),
	}
	zLogger := &ZapLogger{logLevel: logLevelMap[config.Level]}
	cores, err := zLogger.cores(config)
	l, err := zConfig.Build(cores)
	if err != nil {
		return nil, err
	}
	zLogger.zap = l.Sugar()
	logger = zLogger.WithCallDepth(callDepth)
	return zLogger, nil
}

func (l *ZapLogger) cores(config *Config) (zap.Option, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = l.timeEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.MessageKey = "msg"
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.CallerKey = "caller"
	encoderConfig.NameKey = "logger"

	var encoder zapcore.Encoder
	if config.IsJson {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
		encoder.AddInt("PId", os.Getpid())
	} else {
		encoderConfig.EncodeLevel = l.capitalColorLevelEncoder
		encoderConfig.EncodeCaller = l.customCallerEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var cores []zapcore.Core
	if config.LogPath != "" {
		writer := l.getWriter(config)
		cores = append(cores, zapcore.NewCore(encoder, writer, zap.NewAtomicLevelAt(l.logLevel)))
	}
	if config.IsStdout {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zap.NewAtomicLevelAt(l.logLevel)))
	}

	// 确保至少有一个核心
	if len(cores) == 0 {
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			zap.NewAtomicLevelAt(l.logLevel),
		))
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	}), nil
}

func (l *ZapLogger) capitalColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s, ok := _levelToCapitalColorString[level]
	if !ok {
		s = _unknownLevelColor[zapcore.ErrorLevel]
	}
	pid := formatString(fmt.Sprintf("[PId:%d]", os.Getpid()), 15, true)
	color := _levelToColor[level]
	enc.AppendString(s)
	enc.AppendString(color.Add(pid))
}

func (l *ZapLogger) timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	layout := "2006-01-02 15:04:05.000"
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}
	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}
	enc.AppendString(t.Format(layout))
}

func (l *ZapLogger) customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString(fmt.Sprintf("%s", caller.FullPath())) //使用这种可以直接点击跳转到对应的代码行
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

// // 自定义的WriteSyncer
func (l *ZapLogger) getWriter(config *Config) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.LogPath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func (l *ZapLogger) Debug(msg string, keysAndValues ...any) {
	if l.logLevel > zapcore.DebugLevel {
		return
	}
	l.zap.Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) Info(msg string, keysAndValues ...any) {
	if l.logLevel > zapcore.InfoLevel {
		return
	}
	l.zap.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(msg string, err error, keysAndValues ...any) {
	if l.logLevel > zapcore.WarnLevel {
		return
	}
	keysAndValues = appendError(keysAndValues, err)
	l.zap.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(msg string, err error, keysAndValues ...any) {
	if l.logLevel > zapcore.ErrorLevel {
		return
	}
	keysAndValues = appendError(keysAndValues, err)
	l.zap.Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) Panic(msg string, err error, keysAndValues ...any) {
	if l.logLevel > zapcore.PanicLevel {
		return
	}
	l.zap.Panicw(msg, keysAndValues...)
}

func (l *ZapLogger) WithCallDepth(depth int) Logger {
	dup := *l
	dup.zap = l.zap.WithOptions(zap.AddCallerSkip(depth))
	return &dup
}

// 根据配置返回编码器类型
func getEncoding(isJson bool) string {
	if isJson {
		return "json"
	}
	return "console"
}

func appendError(keysAndValues []any, err error) []any {
	if err != nil {
		keysAndValues = append(keysAndValues, "error", err.Error())
	}
	return keysAndValues
}

func formatString(text string, length int, alignLeft bool) string {
	if len(text) > length {
		// Truncate the string if it's longer than the desired length
		return text[:length]
	}

	// Create a format string based on alignment preference
	var formatStr string
	if alignLeft {
		formatStr = fmt.Sprintf("%%-%ds", length) // Left align
	} else {
		formatStr = fmt.Sprintf("%%%ds", length) // Right align
	}

	// Use the format string to format the text
	return fmt.Sprintf(formatStr, text)
}
