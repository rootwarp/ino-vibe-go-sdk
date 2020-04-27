package device

import (
	"context"
	"crypto/x509"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "bitbucket.org/ino-on/ino-vibe-api"
)

// Reader provides retrieve operations for Device.
type Reader interface {
	List(context.Context, pb.InstallStatus) (*pb.DeviceListResponse, error)
	Detail(context.Context, string) (*pb.DeviceResponse, error)
}

// Writer provides update operations for Device.
type Writer interface {
	UpdateInfo(context.Context, *pb.DeviceInfoUpdateRequest) (*pb.DeviceResponse, error)
	UpdateStatus(context.Context, *pb.DeviceStatusUpdateRequest) (*pb.DeviceResponse, error)
	UpdateConfig(context.Context, *pb.DeviceConfigUpdateRequest) (*pb.DeviceResponse, error)
}

const (
	serverURL = "device.ino-vibe.ino-on.dev:443"
)

// Client is client for device instance.
type Client struct {
	deviceClient pb.DeviceServiceClient
}

func (c *Client) getDeviceClient() pb.DeviceServiceClient {
	if c.deviceClient == nil {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			log.Panicln(err)
		}

		creds := credentials.NewClientTLSFromCert(certPool, "")
		conn, _ := grpc.Dial(
			serverURL,
			grpc.WithTransportCredentials(creds),
		)
		c.deviceClient = pb.NewDeviceServiceClient(conn)
	}

	return c.deviceClient
}

// List returns slice of devices.
func (c *Client) List(ctx context.Context, installStatus pb.InstallStatus) (*pb.DeviceListResponse, error) {
	cli := c.getDeviceClient()

	req := pb.DeviceListRequest{
		InstallStatus: installStatus,
	}
	resp, err := cli.List(context.Background(), &req)

	return resp, err
}

// Detail returns detail information of selected device.
func (c *Client) Detail(ctx context.Context, devid string) (*pb.DeviceResponse, error) {
	cli := c.getDeviceClient()

	req := pb.DeviceRequest{
		Devid: devid,
	}
	resp, err := cli.Detail(context.Background(), &req)
	return resp, err
}

// UpdateInfo update basic information of device.
func (c *Client) UpdateInfo(ctx context.Context, req *pb.DeviceInfoUpdateRequest) (*pb.DeviceResponse, error) {
	cli := c.getDeviceClient()
	return cli.UpdateInfo(context.Background(), req)
}

// UpdateStatus updates status information of device.
func (c *Client) UpdateStatus(ctx context.Context, req *pb.DeviceStatusUpdateRequest) (*pb.DeviceResponse, error) {
	cli := c.getDeviceClient()
	return cli.UpdateStatus(context.Background(), req)
}

// UpdateConfig updates device configs.
func (c *Client) UpdateConfig(ctx context.Context, req *pb.DeviceConfigUpdateRequest) (*pb.DeviceResponse, error) {
	cli := c.getDeviceClient()
	return cli.UpdateConfig(context.Background(), req)
}

// NewClient create client.
func NewClient() (*Client, error) {
	return &Client{}, nil
}
