package log

import "github.com/sirupsen/logrus"

var std = newWithStackLevel(LevelInfo, 3)

func Debug(format string, args ...interface{}) {
	std.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	std.Info(format, args...)
}

func Warning(format string, args ...interface{}) {
	std.Warning(format, args...)
}

func Error(format string, args ...interface{}) {
	std.Error(format, args...)
}

func Panic(format string, args ...interface{}) {
	std.Panic(format, args...)
}

func SetLevel(logLevel Level) {
	std.SetLevel(logLevel)
}

func GetLevel() Level {
	switch std.l.Level {
	case logrus.DebugLevel:
		return LevelDebug
	case logrus.InfoLevel:
		return LevelInfo
	case logrus.WarnLevel:
		return LevelWarn
	case logrus.ErrorLevel:
		return LevelError
	case logrus.PanicLevel:
		return LevelPanic
	}
	return LevelUnknown
}

func AddHook(h Hook) {
	std.AddHook(h)
}
