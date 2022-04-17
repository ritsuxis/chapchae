package chat

import "context"

type client struct {
	Host     string
	UserName string
}

func Client(host, username string) *client {
	return &client{
		Host:     host,
		UserName: username,
	}
}

func (*client) Run(ctx context.Context) error {
	return nil
}
