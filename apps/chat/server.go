package chat

import "context"

type server struct {
	Host string
}

func Server(host string) *server {
	return &server{
		Host: host,
	}
}

func (*server) Run(ctx context.Context) error {
	return nil
}
