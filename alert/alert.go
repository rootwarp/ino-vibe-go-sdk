package alert

import (
	"context"
	"crypto/x509"

	pb "bitbucket.org/ino-on/ino-vibe-api"
	iv_auth "github.com/rootwarp/ino-vibe-go-sdk/auth"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var (
	serverURL = "grpc.ino-vibe.ino-on.dev:443"
)

// Client provides API interfaces to access alerts.
type Client interface {
	List(ctx context.Context, request *pb.AlertListRequest) (*pb.AlertListResponse, error)
	Close()
}

type client struct {
	oauthToken  *oauth2.Token
	conn        *grpc.ClientConn
	alertClient pb.AlertServiceClient
}

func (c *client) List(ctx context.Context, request *pb.AlertListRequest) (*pb.AlertListResponse, error) {
	return c.alertClient.List(ctx, request)
}

func (c *client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

// NewClient creates client.
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

	alertClient := pb.NewAlertServiceClient(conn)

	return &client{oauthToken: token, alertClient: alertClient, conn: conn}, nil
}
