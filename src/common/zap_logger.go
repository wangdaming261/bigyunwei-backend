package common

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(logLevel string) *zap.Logger {
	atomicLevel := zap.NewAtomicLevel()
	switch logLevel {
	case "DEBUG":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "INFO":
		atomicLevel.SetLevel(zap.InfoLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time", //时间字段
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "line",                                             //调用者
		MessageKey:     "msg",                                              //内容
		FunctionKey:    "func",                                             //函数名
		StacktraceKey:  "stacktrace",                                       //堆栈
		LineEnding:     zapcore.DefaultLineEnding,                          //换行字符
		EncodeLevel:    zapcore.LowercaseLevelEncoder,                      //小写
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"), //时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	zapCoreConsole := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		atomicLevel,
	)
	core := zapcore.NewTee(
		zapCoreConsole,
	)
	return zap.New(core, zap.AddCaller())

}
