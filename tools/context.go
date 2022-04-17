package tool

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type contextKey string

const (
	debug contextKey = "debug"
	logger   contextKey = "logger"
)

func SignalContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("listening for shutdown signal...")
		<-sigs // 値を受け取る必要はない
		log.Printf("shutdown signal received!")
		signal.Stop(sigs)
		close(sigs)
		cancel()
	}()
	return ctx
}

func SetDebug(ctx context.Context, v bool) context.Context{
	return context.WithValue(ctx, debug, v)
}

func SetZap(ctx context.Context, v zap.Logger) context.Context{
	return context.WithValue(ctx, logger, v)
}

func GetZap(ctx context.Context) (zap.Logger, error) {
    v := ctx.Value(logger)

    logger, ok := v.(zap.Logger)
    if !ok {
        return zap.Logger{}, fmt.Errorf("key:logger not set")
    }

    return logger, nil
}

func GetDebug(ctx context.Context) (bool, error) {
	v := ctx.Value(debug)

	debug, ok := v.(bool)
	if !ok {
        return false, fmt.Errorf("key:debug not set")
    }

    return debug, nil
}
