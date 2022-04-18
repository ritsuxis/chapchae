package chat

import (
	"context"
	"sync"

	chat "github.com/ritsuxis/chapchae/protoc"
	tool "github.com/ritsuxis/chapchae/tools"
	"google.golang.org/grpc"
)

type server struct {
	chat.UnimplementedChatServer
	Host string

	Broadcast chan *chat.StreamResponse

	ClientNames map[string]string
	ClientStreams map[string]chan *chat.StreamResponse

	namesMtx, streamsMtx sync.RWMutex
}

func Server(host string) *server {
	return &server{
		Host: host,

		Broadcast: make(chan *chat.StreamResponse, 1000),

		ClientNames: make(map[string]string),
		ClientStreams: make(map[string]chan *chat.StreamResponse),
	}
}

func (s *server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tool.ServerLogf(ctx, "starting on " + s.Host)

	srv := grpc.NewServer()
	chat.RegisterChatServer(srv, s)

	return nil
}

func (s *server) Stream(srv chat.Chat_StreamServer) error {
	return nil
}
