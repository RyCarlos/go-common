// Package log -----------------------------
// @file      : logger.go
// @author    : Carlos
// @contact   : 534994749@qq.com
// @time      : 2025/6/17 16:10
// -------------------------------------------
package log

type Logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, err error, keysAndValues ...any)
	Error(msg string, err error, keysAndValues ...any)
	Panic(msg string, err error, keysAndValues ...any)
}
