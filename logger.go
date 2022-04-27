package go_logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

const (
	logTmFmtWithMS = "2006-01-02 15:04:05.000"
)

func InitLogger(writer *RedisWriter, Level zapcore.Level) *zap.Logger {
	return newLogger(writer, Level)

}
func newLogger(syncWriter io.Writer, Level zapcore.Level) *zap.Logger {
	lowPriority := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= Level
	})
	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format(logTmFmtWithMS) + "]")
	}
	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	// 自定义文件：行号输出项
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}
	encoderConf := zapcore.EncoderConfig{
		CallerKey:      "caller", // 打印文件名和行数
		LevelKey:       "level",
		MessageKey:     "msg",
		TimeKey:        "ts",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customTimeEncoder,   // 自定义时间格式
		EncodeLevel:    customLevelEncoder,  // 小写编码器
		EncodeCaller:   customCallerEncoder, // 全路径编码器
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 使用 JSON 格式日志
	jsonEnc := zapcore.NewJSONEncoder(encoderConf)
	stdCore := zapcore.NewCore(jsonEnc, zapcore.Lock(os.Stdout), lowPriority)
	syncer := zapcore.AddSync(syncWriter)
	syncCore := zapcore.NewCore(jsonEnc, syncer, lowPriority)
	if syncWriter != nil {
		return zap.New(stdCore).WithOptions(zap.AddCaller())
	} else {
		core := zapcore.NewTee(stdCore, syncCore)
		return zap.New(core).WithOptions(zap.AddCaller())
	}
}
