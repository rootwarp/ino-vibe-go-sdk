package device

import (
	"context"
	"crypto/x509"
	"errors"
	"log"
	"math"
	"time"

	"cloud.google.com/go/datastore"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	pb "bitbucket.org/ino-on/ino-vibe-api"
	iv_auth "github.com/rootwarp/ino-vibe-go-sdk/auth"
)

const (
	statusLogKind   = "DevStatusLog"
	inclinationKind = "inclination-log"
)

var (
	serverURL = "grpc.ino-vibe.ino-on.dev:443"
)

// Errors
var (
	ErrInvalidParameter        = errors.New("Invalid parameter value")
	ErrNonExistDevice          = errors.New("Request on non-exist device")
	ErrForbiddenInstallStatus  = errors.New("Request is not permitted on current install status")
	ErrInvalidInclinationValue = errors.New("Requested inclination is NaN or Inf")
	ErrNoEntities              = errors.New("Device has no valid entity")
)

// Client is client for device instance.
type Client interface {
	List(context.Context, pb.InstallStatus) (*pb.DeviceListResponse, error)
	Detail(context.Context, string) (*pb.DeviceResponse, error)

	UpdateInfo(context.Context, *pb.DeviceInfoUpdateRequest) (*pb.DeviceResponse, error)
	UpdateStatus(context.Context, *pb.DeviceStatusUpdateRequest) (*pb.DeviceResponse, error)
	UpdateConfig(context.Context, *pb.DeviceConfigUpdateRequest) (*pb.DeviceResponse, error)

	StatusLog(ctx context.Context, devid, installKey string, timeFrom, timeTo time.Time, offset, limit int) ([]StatusLog, error)
	StoreStatusLog(ctx context.Context, devid string, battery, temperature, RSSI int) error

	LastInclinationLog(context.Context, string) (*InclinationLog, error)
	StoreInclinationLog(context.Context, string, int, int, int) error

	PrepareInstall(context.Context, *pb.PrepareInstallRequest) (*pb.PrepareInstallResponse, error)
	CompleteInstall(context.Context, *pb.CompleteInstallRequest) (*pb.CompleteInstallResponse, error)
	WaitCompleteInstall(context.Context, *pb.WaitCompleteInstallRequest) (*pb.WaitCompleteInstallResponse, error)
	Uninstalling(context.Context, *pb.UninstallingRequest) (*pb.UninstallingResponse, error)
	Uninstall(context.Context, *pb.UninstallRequest) (*pb.UninstallResponse, error)
	Discard(context.Context, *pb.DiscardRequest) (*pb.DiscardResponse, error)
}

type client struct {
	oauthToken   *oauth2.Token
	deviceClient pb.DeviceServiceClient
	dsClient     *datastore.Client
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

func (c *client) getDatastoreClient() *datastore.Client {
	if c.dsClient == nil {
		var err error

		ctx := context.Background()
		cred, err := google.FindDefaultCredentials(ctx)
		if err != nil {
			log.Fatal(err)
		}

		c.dsClient, err = datastore.NewClient(ctx, cred.ProjectID)
		if err != nil {
			log.Fatal(err)
		}
	}

	return c.dsClient
}

// List returns slice of devices.
func (c *client) List(ctx context.Context, installStatus pb.InstallStatus) (*pb.DeviceListResponse, error) {
	cli := c.getDeviceClient()

	req := pb.DeviceListRequest{
		InstallStatus: installStatus,
	}
	resp, err := cli.List(ctx, &req)

	return resp, err
}

// Detail returns detail information of selected device.
func (c *client) Detail(ctx context.Context, devid string) (*pb.DeviceResponse, error) {
	cli := c.getDeviceClient()

	req := pb.DeviceRequest{
		Devid: devid,
	}
	resp, err := cli.Detail(ctx, &req)
	return resp, err
}

// UpdateInfo update basic information of device.
func (c *client) UpdateInfo(ctx context.Context, req *pb.DeviceInfoUpdateRequest) (*pb.DeviceResponse, error) {
	cli := c.getDeviceClient()
	return cli.UpdateInfo(ctx, req)
}

// UpdateStatus updates status information of device.
func (c *client) UpdateStatus(ctx context.Context, req *pb.DeviceStatusUpdateRequest) (*pb.DeviceResponse, error) {
	cli := c.getDeviceClient()
	return cli.UpdateStatus(ctx, req)
}

// UpdateConfig updates device configs.
func (c *client) UpdateConfig(ctx context.Context, req *pb.DeviceConfigUpdateRequest) (*pb.DeviceResponse, error) {
	cli := c.getDeviceClient()
	return cli.UpdateConfig(ctx, req)
}

// StatusLog returns slice of status log of selected device within time range.
// ErrNonExistDevice
// ErrNoEntities
// ErrInvalidParameter
func (c *client) StatusLog(ctx context.Context, devid, installSession string, timeFrom, timeTo time.Time, offset, limit int) ([]StatusLog, error) {
	if timeFrom.After(timeTo) {
		return []StatusLog{}, ErrInvalidParameter
	}

	if offset < 0 || limit <= 0 {
		return []StatusLog{}, ErrInvalidParameter
	}

	_, err := c.getDevice(ctx, devid)
	if err != nil {
		return []StatusLog{}, err
	}

	dsCli := c.getDatastoreClient()

	q := datastore.NewQuery(statusLogKind).
		Filter("Devid =", devid).
		Filter("InstallSessionKey =", installSession).
		Order("-Time").
		Offset(offset).
		Limit(limit)

	iter := dsCli.Run(ctx, q)

	logs := make([]StatusLog, 0, limit)
	idx := 0

	for {
		newLog := StatusLog{}
		_, err := iter.Next(&newLog)
		if err == iterator.Done {
			break
		}

		if err, ok := err.(*datastore.ErrFieldMismatch); ok {
			log.Println("StatusLog", err)
		} else if err != nil {
			return []StatusLog{}, err
		}

		logs = append(logs, newLog)
		idx++
	}

	if idx == 0 {
		return logs, ErrNoEntities
	}

	return logs, nil
}

func (c *client) getDevice(ctx context.Context, devid string) (*pb.Device, error) {
	cli := c.getDeviceClient()

	resp, err := cli.Detail(ctx, &pb.DeviceRequest{Devid: devid})
	if err != nil {
		return nil, err
	}

	if resp.ResultCode != pb.ResponseCode_SUCCESS {
		return nil, ErrNonExistDevice
	}

	device := resp.Devices[0]

	return device, nil
}

// StoreStatusLog stores requested values into StatusLog entity.
// New status log can be stored onto installed device or return ErrForbeddenInstallStatus error.
//
// ErrNonExistDevice returns if selected device is not exist.
// ErrForbiddenInstallStatus if selected device is not installed.
func (c *client) StoreStatusLog(ctx context.Context, devid string, battery, temperature, RSSI int) error {
	device, err := c.getDevice(ctx, devid)
	if err != nil {
		return err
	}

	if device.InstallStatus != pb.InstallStatus_Installed {
		return ErrForbiddenInstallStatus
	}

	newLog := StatusLog{
		Devid:             devid,
		Time:              time.Now(),
		Temperature:       temperature,
		Battery:           battery,
		RSSI:              RSSI,
		InstallSessionKey: device.InstallSessionKey,
	}

	dsCli := c.getDatastoreClient()
	newKey := datastore.IncompleteKey(statusLogKind, nil)
	_, err = dsCli.Put(ctx, newKey, &newLog)

	return err
}

// LastInclinationLog try to get latest inclination log of selected device.
// InstallSessionKey value from log should be same to current status of device.
//
// ErrNonExistDevice returns if requested device is not exist.
// ErrNoEntities returns if inclination log with current install session key does not exist.
func (c *client) LastInclinationLog(ctx context.Context, devid string) (*InclinationLog, error) {
	device, err := c.getDevice(ctx, devid)
	if err != nil {
		return nil, err
	}

	dsCli := c.getDatastoreClient()
	q := datastore.NewQuery(inclinationKind).
		Filter("devid =", devid).
		Order("-time_created").
		Limit(1)

	iter := dsCli.Run(ctx, q)

	latestInclination := InclinationLog{}
	_, err = iter.Next(&latestInclination)
	switch {
	case err == nil && latestInclination.InstallSessionKey != device.InstallSessionKey:
		return nil, ErrNoEntities
	case err == iterator.Done:
		return nil, ErrNoEntities
	case err != nil:
		return nil, err
	}

	return &latestInclination, nil
}

// StoreInclinationLog creates new inclination log.
//
// ErrNonExistDevice returns if requested device is not exist.
// ErrForbiddenInstallStatus returns if requested device is not installed.
// ErrInvalidInclinationValue returns if calculated angle is NaN or Inf.
func (c *client) StoreInclinationLog(ctx context.Context, devid string, rawX, rawY, rawZ int) error {
	device, err := c.getDevice(ctx, devid)
	if err != nil {
		return err
	}

	if device.InstallStatus != pb.InstallStatus_Installed {
		return ErrForbiddenInstallStatus
	}

	var unit float64
	if device.DevType == pb.DeviceType_InoVibe {
		unit = 3.9
	} else {
		unit = 0.244
	}

	x, y, z := float64(rawX)*unit, float64(rawY)*unit, float64(rawZ)*unit
	angleZ := angle(x, y, z, unit)

	if math.IsNaN(angleZ) || math.IsInf(angleZ, 0) {
		return ErrInvalidInclinationValue
	}

	newLog := InclinationLog{
		Devid:             devid,
		Time:              time.Now(),
		InstallSessionKey: device.InstallSessionKey,
		AccXMg:            x,
		AccYMg:            y,
		AccZMg:            z,
		AngleZ:            angleZ,
	}

	dsCli := c.getDatastoreClient()
	newKey := datastore.IncompleteKey(inclinationKind, nil)
	_, err = dsCli.Put(ctx, newKey, &newLog)

	return err
}

func angle(x, y, z, unit float64) float64 {
	t := math.Sqrt(math.Pow(x, 2)+math.Pow(y, 2)) / z
	rad := math.Atan(t)
	return rad * (180 / math.Pi)
}

func (c *client) PrepareInstall(ctx context.Context, in *pb.PrepareInstallRequest) (*pb.PrepareInstallResponse, error) {
	cli := c.getDeviceClient()
	return cli.PrepareInstall(ctx, in)
}

func (c *client) CompleteInstall(ctx context.Context, in *pb.CompleteInstallRequest) (*pb.CompleteInstallResponse, error) {
	cli := c.getDeviceClient()
	return cli.CompleteInstall(ctx, in)
}

func (c *client) Uninstalling(ctx context.Context, in *pb.UninstallingRequest) (*pb.UninstallingResponse, error) {
	cli := c.getDeviceClient()
	return cli.Uninstalling(ctx, in)
}

func (c *client) Uninstall(ctx context.Context, in *pb.UninstallRequest) (*pb.UninstallResponse, error) {
	cli := c.getDeviceClient()
	return cli.Uninstall(ctx, in)
}

func (c *client) Discard(ctx context.Context, in *pb.DiscardRequest) (*pb.DiscardResponse, error) {
	cli := c.getDeviceClient()
	return cli.Discard(ctx, in)
}

func (c *client) WaitCompleteInstall(ctx context.Context, in *pb.WaitCompleteInstallRequest) (*pb.WaitCompleteInstallResponse, error) {
	cli := c.getDeviceClient()
	return cli.WaitCompleteInstall(ctx, in)
}

// NewClient create client.
func NewClient() (Client, error) {
	token, err := iv_auth.LoadCredentials()
	if err != nil {
		log.Panicln(err)
	}

	return &client{oauthToken: token}, nil
}
