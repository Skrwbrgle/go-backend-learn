package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

func InitLogger(env string) {
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/app.log", // path file log
		MaxSize:    10,               // MB
		MaxBackups: 5,                // berapa file backup
		MaxAge:     30,               // hari
		Compress:   true,             // compress pakai gzip
	})

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.LevelKey = "level"
	encoderCfg.CallerKey = "caller"
	encoderCfg.MessageKey = "msg"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),                              // JSON format
		zapcore.NewMultiWriteSyncer(writer, zapcore.AddSync(os.Stdout)), // ke file + stdout
		zapcore.InfoLevel, // level minimal (Info, bisa diubah ke Debug)
	)

	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	if env != "production" {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)
		// Combine file logger (JSON) + console (readable)
		core = zapcore.NewTee(core, consoleCore)
		Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}
}
