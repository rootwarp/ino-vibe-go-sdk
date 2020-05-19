package device

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

const (
	serverURL = "device.ino-vibe.ino-on.dev:443"
)

// Client is client for device instance.
type Client interface {
	List(context.Context, pb.InstallStatus) (*pb.DeviceListResponse, error)
	Detail(context.Context, string) (*pb.DeviceResponse, error)

	UpdateInfo(context.Context, *pb.DeviceInfoUpdateRequest) (*pb.DeviceResponse, error)
	UpdateStatus(context.Context, *pb.DeviceStatusUpdateRequest) (*pb.DeviceResponse, error)
	UpdateConfig(context.Context, *pb.DeviceConfigUpdateRequest) (*pb.DeviceResponse, error)

	StatusLog(context.Context, *pb.StatusLogRequest) (*pb.StatusLogResponse, error)
}

type client struct {
	oauthToken   *oauth2.Token
	deviceClient pb.DeviceServiceClient
}

func (c *client) getDeviceClient() pb.DeviceServiceClient {
	if c.deviceClient == nil {
		if c.oauthToken == nil {
			log.Panicln(errors.New("No credentials"))
		}

		certPool, err := x509.SystemCertPool()
		if err != nil {
			log.Panicln(err)
		}

		creds := credentials.NewClientTLSFromCert(certPool, "")
		conn, _ := grpc.Dial(
			serverURL,
			grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(oauth.NewOauthAccess(c.oauthToken)),
		)
		c.deviceClient = pb.NewDeviceServiceClient(conn)
	}

	return c.deviceClient
}

// List returns slice of devices.
func (c *client) List(ctx context.Context, installStatus pb.InstallStatus) (*pb.DeviceListResponse, error) {
	cli := c.getDeviceClient()

	req := pb.DeviceListRequest{
		InstallStatus: installStatus,
	}
	resp, err := cli.List(context.Background(), &req)

	return resp, err
}

// Detail returns detail information of selected device.
func (c *client) Detail(ctx context.Context, devid string) (*pb.DeviceResponse, error) {
	cli := c.getDeviceClient()

	req := pb.DeviceRequest{
		Devid: devid,
	}
	resp, err := cli.Detail(context.Background(), &req)
	return resp, err
}

// UpdateInfo update basic information of device.
func (c *client) UpdateInfo(ctx context.Context, req *pb.DeviceInfoUpdateRequest) (*pb.DeviceResponse, error) {
	cli := c.getDeviceClient()
	return cli.UpdateInfo(context.Background(), req)
}

// UpdateStatus updates status information of device.
func (c *client) UpdateStatus(ctx context.Context, req *pb.DeviceStatusUpdateRequest) (*pb.DeviceResponse, error) {
	cli := c.getDeviceClient()
	return cli.UpdateStatus(context.Background(), req)
}

// UpdateConfig updates device configs.
func (c *client) UpdateConfig(ctx context.Context, req *pb.DeviceConfigUpdateRequest) (*pb.DeviceResponse, error) {
	cli := c.getDeviceClient()
	return cli.UpdateConfig(context.Background(), req)
}

func (c *client) StatusLog(ctx context.Context, req *pb.StatusLogRequest) (*pb.StatusLogResponse, error) {
	cli := c.getDeviceClient()
	return cli.StatusLog(ctx, req)
}

// NewClient create client.
func NewClient() (Client, error) {
	token, err := iv_auth.LoadCredentials()
	if err != nil {
		log.Panicln(err)
	}

	return &client{oauthToken: token}, nil
}
