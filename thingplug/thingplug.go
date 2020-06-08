package thingplug

import (
	"context"
	"crypto/x509"
	"fmt"

	pb "bitbucket.org/ino-on/ino-vibe-api"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	iv_auth "github.com/rootwarp/ino-vibe-go-sdk/auth"
)

var (
	serverURL = "thingplug.ino-vibe.ino-on.dev:443"
)

// Client provides control interfaces for Ino-Vibe.
type Client interface {
	PowerOff(ctx context.Context, devid string) error
	BaseReset(ctx context.Context, devid string) error
	Reset(ctx context.Context, devid string) error
	Close()
}

type client struct {
	oauthToken      *oauth2.Token
	conn            *grpc.ClientConn
	thingplugClient pb.ThingplugServiceClient
}

func (c *client) PowerOff(ctx context.Context, devid string) error {
	resp, err := c.thingplugClient.PowerOff(ctx, &pb.ThingplugDeviceRequest{Devid: devid})
	fmt.Println(resp)

	return err
}

func (c *client) BaseReset(ctx context.Context, devid string) error {
	resp, err := c.thingplugClient.BaseReset(ctx, &pb.ThingplugDeviceRequest{Devid: devid})
	fmt.Println(resp)

	return err
}

func (c *client) Reset(ctx context.Context, devid string) error {
	resp, err := c.thingplugClient.Reset(ctx, &pb.ThingplugDeviceRequest{Devid: devid})
	fmt.Println(resp)

	return err
}

func (c *client) Close() {
	c.conn.Close()
}

// NewClient create new client.
func NewClient() (Client, error) {
	token, err := iv_auth.LoadCredentials()
	if err != nil {
		return nil, err
	}

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	creds := credentials.NewClientTLSFromCert(certPool, "")
	conn, err := grpc.Dial(
		serverURL,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(oauth.NewOauthAccess(token)),
	)
	if err != nil {
		return nil, err
	}

	thingplugClient := pb.NewThingplugServiceClient(conn)

	return &client{oauthToken: token, conn: conn, thingplugClient: thingplugClient}, nil
}
