package main

import (
	"context"
	"flag"
	"os"

	"github.com/ritsuxis/chapchae/apps/chat"
	tool "github.com/ritsuxis/chapchae/tools"
	"go.uber.org/zap"
)

var (
	serverMode bool
	debug      bool

	host     string // .envとかから持ってくるようにする
	username string
)

func init() {
	flag.BoolVar(&serverMode, "s", false, "run as the server")
	flag.BoolVar(&debug, "v", false, "enable debug logging")
	flag.StringVar(&host, "h", "0.0.0.0:6262", "the chat server's host")
	// flag.StringVar(&password, "p", "", "the chat server's password")
	flag.StringVar(&username, "n", "", "the username for the client")
	flag.Parse()
}

func main() {
	// setting zap
	logger, err := zap.NewDevelopment()
	// logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	zap.ReplaceGlobals(logger)

	// context
	ctx := tool.SignalContext(context.Background())
	ctx = tool.SetZap(ctx, *logger)
	ctx = tool.SetDebug(ctx, debug)

	// run service
	if serverMode {
		err = chat.Server(host).Run(ctx)
	} else {
		err = chat.Client(host, username).Run(ctx)
	}

	if err != nil {
		tool.MessageLogf(ctx, "Process", err.Error())
		os.Exit(1)
	}

	tool.MessageLogf(ctx, "Process", "byebye...")
}
