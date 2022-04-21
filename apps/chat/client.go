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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type client struct {
	chat.ChatClient
	Host, Name, Password, Token string
	Shutdown                    bool
}

func Client(host, pass, username string) *client {
	return &client{
		Host:     host,
		Password: pass,
		Name:     username,
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

	if c.Token, err = c.login(ctx); err != nil {
		return errors.WithMessage(err, "failed to login")
	}
	tool.ClientLogf("logged in successfully")

	err = c.stream(ctx)

	tool.ClientLogf("logging out")
	if err := c.logout(ctx); err != nil {
		tool.ClientLogf("failed to log out: " + err.Error())
	}

	return errors.WithMessage(err, "stream error")
}

func (c *client) stream(ctx context.Context) error {
	md := metadata.New(map[string]string{tokenHeader: c.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := c.ChatClient.Stream(ctx)
	if err != nil {
		return err
	}
	defer client.CloseSend()

	tool.ClientLogf("connected to stream")

	go c.send(client)
	return c.receive(client)
}

func (c *client) receive(sc chat.Chat_StreamClient) error {
	for {
		res, err := sc.Recv()

		if s, ok := status.FromError(err); ok && s.Code() == codes.Canceled {
			tool.DebugLogf("stream canceled (usually indicates shutdown)")
			return nil
		} else if err == io.EOF {
			tool.DebugLogf("stream closed by server")
			return nil
		} else if err != nil {
			return err
		}

		switch evt := res.Event.(type) {
		case *chat.StreamResponse_ClientLogin:
			tool.ServerLogf("%s has logged in", evt.ClientLogin.Name)
		case *chat.StreamResponse_ClientLogout:
			tool.ServerLogf("%s has logged out", evt.ClientLogout.Name)
		case *chat.StreamResponse_ClientMessage:
			tool.MessageLogf(evt.ClientMessage.Name, evt.ClientMessage.Message)
		case *chat.StreamResponse_ServerShutdown:
			tool.ServerLogf("the server is shutting down")
			c.Shutdown = true
			return nil
		default:
			tool.ClientLogf("unexpected event from the server: %T", evt)
			return nil
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
			tool.DebugLogf("client send loop disconnected")
		default:
			if sc.Scan() {
				// send message
				if err := client.Send(&chat.StreamRequest{Message: sc.Text()}); err != nil {
					tool.ClientLogf("failed to send message: " + err.Error())
					return
				}
			} else {
				tool.ClientLogf("input scanner failure: " + sc.Err().Error())
				return
			}
		}
	}
}

func (c *client) login(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	res, err := c.ChatClient.Login(ctx, &chat.LoginRequest{
		Name:     c.Name,
		Password: c.Password,
	})

	if err != nil {
		return "", err
	}

	return res.Token, nil
}

func (c *client) logout(_ context.Context) error {
	if c.Shutdown {
		tool.DebugLogf("unable to logout (server sent shutdown signal)")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.ChatClient.Logout(ctx, &chat.LogoutRequest{Token: c.Token})
	if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
		tool.DebugLogf("unable to logout (connection already closed)")
		return nil
	}

	return err
}
