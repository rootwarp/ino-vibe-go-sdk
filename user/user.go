package user

import (
	"context"
	"crypto/x509"
	"errors"
	"log"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	pb "bitbucket.org/ino-on/ino-vibe-api"
	iv_auth "github.com/rootwarp/ino-vibe-go-sdk/auth"
)

// Writer provides insterfaces for modify user's information.
type Writer interface {
	RegisterDeviceToken(userID, username, deviceToken string) error
}

// Client is client for write user's information.
type Client struct {
	oauthToken *oauth2.Token
}

const (
	serverURL = "user.ino-vibe.ino-on.dev:443"
)

// RegisterDeviceToken register device token to receive mobile push notification.
func (c *Client) RegisterDeviceToken(userID, username, deviceToken string) error {
	if c.oauthToken == nil {
		log.Panicln(errors.New("No credentials"))
	}

	// TODO: Refactoring. Will be duplicated.
	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Panicln(err)
	}

	creds := credentials.NewClientTLSFromCert(certPool, "")
	conn, err := grpc.Dial(
		serverURL,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(oauth.NewOauthAccess(c.oauthToken)),
	)
	if err != nil {
		return err
	}

	ctx := context.Background()
	req := pb.RegisterDeviceTokenRequest{
		UserId:      userID,
		Username:    username,
		DeviceToken: deviceToken,
	}

	cli := pb.NewUserServiceClient(conn)
	resp, err := cli.RegisterDeviceToken(ctx, &req)

	log.Println("Resp ", resp, err)

	return err
}

// NewClient creates client.
func NewClient() (*Client, error) {
	token, err := iv_auth.LoadCredentials()
	if err != nil {
		log.Panicln(err)
	}

	return &Client{oauthToken: token}, nil
}
