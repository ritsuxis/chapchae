package tool

import (
	"context"
	"log"
	"time"
)

const timeFormat = "2006/01/02 15:04:05"

func ServerLogf(ctx context.Context, format string) {
	logger, err := GetZap(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	logger.Sugar().Infow("<<Server>>",
		"message", format,
		"Time", time.Now().Format(timeFormat),
	)
}

func ClientLogf(ctx context.Context, format string) {
	logger, err := GetZap(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	logger.Sugar().Infow("<<Client>>",
		"message", format,
		"Time", time.Now().Format(timeFormat),
	)
}

func MessageLogf(ctx context.Context, name string ,msg string) {
	logger, err := GetZap(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	logger.Sugar().Infow("<<" + name + ">>",
		"message", msg,
		"Time", time.Now().Format(timeFormat),
	)
}

func DebugLogf(ctx context.Context, format string) {
	debug, err := GetDebug(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	if !debug {
		return
	}
	logger, err := GetZap(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	logger.Sugar().Infow("<<Debug>>",
		"message", format,
		"Time", time.Now().Format(timeFormat),
	)
}
