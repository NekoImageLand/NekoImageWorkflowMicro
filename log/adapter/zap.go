package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"os"
	"strconv"
	"time"
)

type customEncoder struct {
	zapcore.Encoder
}

func (c *customEncoder) Clone() zapcore.Encoder {
	return &customEncoder{c.Encoder.Clone()}
}

func (c *customEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	levelColor := ""
	switch entry.Level {
	case zapcore.DebugLevel:
		levelColor = "36" // Cyan
	case zapcore.InfoLevel:
		levelColor = "32" // Green
	case zapcore.WarnLevel:
		levelColor = "33" // Yellow
	case zapcore.ErrorLevel:
		levelColor = "31" // Red
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		levelColor = "35" // Magenta
	default:
		levelColor = "37" // White
	}
	fileLine := "unknown.go:0"
	if entry.Caller.Defined {
		fileLine = entry.Caller.TrimmedPath() + ":" + strconv.Itoa(entry.Caller.Line)
	}
	log := fmt.Sprintf("\033[32m%s \033[0m[\033[%sm%s\033[0m] \033[34;4m%s\033[0m | %s\n",
		time.Now().Format("01-02 15:04:05"), levelColor, entry.Level.String(), fileLine, entry.Message)
	buf := buffer.Buffer{}
	buf.AppendString(log)
	return &buf, nil
}

func NewCustomLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		&customEncoder{zapcore.NewConsoleEncoder(encoderConfig)},
		zapcore.Lock(os.Stdout),
		zapcore.DebugLevel,
	)

	return zap.New(core)
}
