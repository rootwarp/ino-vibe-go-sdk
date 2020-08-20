package wave

import (
	"context"
	"crypto/x509"

	pb "bitbucket.org/ino-on/ino-vibe-api"
	iv_auth "github.com/rootwarp/ino-vibe-go-sdk/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var (
	serverURL = "grpc.ino-vibe.ino-on.dev:443"
)

// Client provides interfaces.
type Client interface {
	Detail(ctx context.Context, req *pb.WaveDetailRequest) (*pb.WaveDetailResponse, error)
	Close()
}

type client struct {
	conn       *grpc.ClientConn
	waveClient pb.WaveServiceClient
}

func (c *client) Detail(ctx context.Context, req *pb.WaveDetailRequest) (*pb.WaveDetailResponse, error) {
	return c.waveClient.Detail(ctx, req)
}

func (c *client) Close() {
	_ = c.conn.Close()
}

// NewClient creates new client instance.
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

	return &client{conn: conn, waveClient: pb.NewWaveServiceClient(conn)}, nil
}
