package chat

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	chat "github.com/ritsuxis/chapchae/protoc"
	tool "github.com/ritsuxis/chapchae/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type client struct {
	chat.ChatClient
	Host, Name string
}

func Client(host, username string) *client {
	return &client{
		Host: host,
		Name: username,
	}
}

func (c *client) Run(ctx context.Context) error {
	connectionCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	// blocking dial until connections to be established.
	// NOTE: NewCredentials returns a credentials which disables transport security.
	conn, err := grpc.DialContext(connectionCtx, c.Host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return errors.WithMessage(err, "unable to connect")
	}
	defer conn.Close()

	c.ChatClient = chat.NewChatClient(conn)

	// TODO: login process

	err = c.stream(ctx)

	// TODO: logout

	return errors.WithMessage(err, "stream error")
}

func (c *client) stream(ctx context.Context) error {
	// TODO: token

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := c.ChatClient.Stream(ctx)
	if err != nil {
		return err
	}
	defer client.CloseSend()

	tool.ClientLogf(ctx, "connected to stream")

	go c.send(client)
	return c.receive(client)
}

func (c *client) receive(sc chat.Chat_StreamClient) error {
	for {
		res, err := sc.Recv()

		if s, ok := status.FromError(err); ok && s.Code() == codes.Canceled {
			tool.DebugLogf(sc.Context(), "stream canceled (usually indicates shutdown)")
			return nil
		} else if err == io.EOF {
			tool.DebugLogf(sc.Context(), "stream closed by server")
			return nil
		} else if err != nil {
			return err
		}

		switch evt := res.Event.(type) {
		// case *chat.StreamResponse_ClientLogin:
		// 	tool.ServerLogf(ts, "%s has logged in", evt.ClientLogin.Name)
		// case *chat.StreamResponse_ClientLogout:
		// 	tool.ServerLogf(ts, "%s has logged out", evt.ClientLogout.Name)
		case *chat.StreamResponse_ClientMessage:
			tool.MessageLogf(sc.Context(), evt.ClientMessage.Name, evt.ClientMessage.Message)
			// case *chat.StreamResponse_ServerShutdown:
			// 	tool.ServerLogf(ts, "the server is shutting down")
			// 	c.Shutdown = true
			// 	return nil
			// default:
			// 	ClientLogf(ts, "unexpected event from the server: %T", evt)
			// 	return nil
		}
	}
}

func (c *client) send(client chat.Chat_StreamClient) {
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanLines)

	// send loop
	for {
		select {
		case <-client.Context().Done():
			tool.DebugLogf(client.Context(), "client send loop disconnected")
		default:
			if sc.Scan() {
				// send message
				if err := client.Send(&chat.StreamRequest{Message: sc.Text()}); err != nil {
					tool.ClientLogf(client.Context(), "failed to send message: "+err.Error())
					return
				}
			} else {
				tool.ClientLogf(client.Context(), "input scanner failure: "+sc.Err().Error())
				return
			}
		}
	}
}
