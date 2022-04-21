package tool

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

const timeFormat = "2006/01/02 15:04:05"

func ServerLogf(format string, args...interface{}) {
	zap.L().Sugar().Infow("<<Server>>",
		"message", fmt.Sprintf(format, args...),
		"Time", time.Now().Format(timeFormat),
	)
}

func ClientLogf(format string, args ...interface{}) {
	zap.L().Sugar().Infow("<<Client>>",
		"message", fmt.Sprintf(format, args...),
		"Time", time.Now().Format(timeFormat),
	)
}

func MessageLogf(name string, msg string, args ...interface{}) {
	zap.L().Sugar().Infow("<<"+name+">>",
		"message", fmt.Sprintf(msg, args...),
		"Time", time.Now().Format(timeFormat),
	)
}

func DebugLogf(format string, args ...interface{}) {
	zap.L().Sugar().Infow("<<Debug>>",
		"message", fmt.Sprintf(format, args...),
		"Time", time.Now().Format(timeFormat),
	)
}
