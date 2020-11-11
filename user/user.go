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

// User is entity structure for describing Auth0 user.
type User struct {
	UserID string
	Email  string
}

// Client is client for write user's information.
type Client interface {
	RegisterDeviceToken(userID, username, deviceToken string) error
	GetDeviceToken(username string) ([]DeviceToken, error)
}

// DeviceToken describes token for FCM.
type DeviceToken struct {
	DeviceName string
	Token      string
}

type client struct {
	oauthToken *oauth2.Token
}

var (
	serverURL = "grpc.ino-vibe.ino-on.dev:443"
)

// RegisterDeviceToken register device token to receive mobile push notification.
func (c *client) RegisterDeviceToken(userID, username, deviceToken string) error {
	if c.oauthToken == nil {
		log.Panicln(errors.New("No credentials"))
	}

	conn, err := c.connection()
	if err != nil {
		return err
	}

	cli := pb.NewUserServiceClient(conn)

	ctx := context.Background()
	req := pb.RegisterDeviceTokenRequest{
		UserId:      userID,
		Username:    username,
		DeviceToken: deviceToken,
	}

	resp, err := cli.RegisterDeviceToken(ctx, &req)

	log.Println("Resp ", resp, err)

	return err
}

func (c *client) connection() (*grpc.ClientConn, error) {
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

	return conn, err
}

func (c *client) GetDeviceToken(username string) ([]DeviceToken, error) {
	if c.oauthToken == nil {
		log.Panicln(errors.New("No credentials"))
	}

	conn, err := c.connection()
	if err != nil {
		return nil, err
	}

	cli := pb.NewUserServiceClient(conn)

	ctx := context.Background()
	req := &pb.GetDeviceTokenRequest{Username: username}
	resp, err := cli.GetDeviceToken(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != pb.ResponseCode_SUCCESS {
		return make([]DeviceToken, 0), nil
	}

	deviceTokens := make([]DeviceToken, len(resp.GetDeviceTokens()))
	for i, token := range resp.GetDeviceTokens() {
		deviceTokens[i] = DeviceToken{
			DeviceName: token.GetDeviceName(),
			Token:      token.GetToken(),
		}
	}

	return deviceTokens, nil
}

// NewClient creates client.
func NewClient() (Client, error) {
	token, err := iv_auth.LoadCredentials()
	if err != nil {
		log.Panicln(err)
	}

	return &client{oauthToken: token}, nil
}
