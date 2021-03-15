package device

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"

	pb "bitbucket.org/ino-on/ino-vibe-api"
)

var (
	db *firestore.Client
)

func init() {
	target := os.Getenv("TEST_TARGET")
	if target != "" {
		serverURL = target + "-" + serverURL
	}
	fmt.Println(serverURL)

	ctx := context.Background()
	option := option.WithCredentialsFile(os.Getenv("FIREBASE_APPLICATION_CREDENTIALS"))
	db, _ = firestore.NewClient(ctx, "crash-detector", option)
}

func TestGetDeviceListUnauthorized(t *testing.T) {
	ctx := context.Background()
	cli, _ := NewClient()
	(cli.(*client)).oauthToken.AccessToken = "invalid-token"

	_, err := cli.List(ctx, pb.InstallStatus_Installed)

	assert.NotNil(t, err)
}

func TestGetInitial(t *testing.T) {
	ctx := context.Background()
	cli, _ := NewClient()

	devs, err := cli.FilterList(ctx, &pb.DeviceFilterListRequest{
		InstallStatus: &pb.DeviceFilterListRequest_InstallStatusValue{
			InstallStatusValue: pb.InstallStatus_Initial,
		},
	})

	fmt.Println(len(devs), err)
	for _, dev := range devs {
		assert.Equal(t, "", dev.InstallSessionKey)
		assert.Equal(t, pb.InstallStatus_Initial, dev.InstallStatus)
	}
}

func TestGetDeviceList(t *testing.T) {
	tests := []pb.InstallStatus{
		pb.InstallStatus_Installed,
		pb.InstallStatus_Uninstalling,
	}

	ctx := context.Background()
	cli, _ := NewClient()

	for _, installStatus := range tests {
		resp, err := cli.List(ctx, installStatus)

		assert.Nil(t, err)
		assert.Equal(t, pb.ResponseCode_SUCCESS, resp.ResultCode)

		for _, device := range resp.Devices {
			assert.Equal(t, installStatus, device.InstallStatus)
		}
	}
}

func TestGetDeviceFilterByGroup(t *testing.T) {
	ctx := context.Background()
	cli, _ := NewClient()

	testGroupID := "0bee7b43-0b57-4b54-9062-430e2bd3fa79"

	devs, err := cli.FilterList(ctx, &pb.DeviceFilterListRequest{
		GroupId: &pb.DeviceFilterListRequest_GroupIdValue{GroupIdValue: testGroupID},
	})

	assert.Nil(t, err)

	for _, dev := range devs {
		assert.Equal(t, testGroupID, dev.GroupId)
	}
}

func TestGetDeviceFilterByInstallStatus(t *testing.T) {
	ctx := context.Background()
	cli, _ := NewClient()

	devs, err := cli.FilterList(ctx, &pb.DeviceFilterListRequest{
		InstallStatus: &pb.DeviceFilterListRequest_InstallStatusValue{
			InstallStatusValue: pb.InstallStatus_Installed,
		},
	})

	assert.Nil(t, err)

	for _, dev := range devs {
		assert.Equal(t, pb.InstallStatus_Installed, dev.InstallStatus)
	}
}

func TestGetDeviceDetailNonExist(t *testing.T) {
	testDevid := "non-exist-device"

	ctx := context.Background()
	cli, _ := NewClient()
	resp, err := cli.Detail(ctx, testDevid)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.ResultCode, pb.ResponseCode_NON_EXIST)
	assert.Equal(t, len(resp.Devices), 0)
}

func TestUpdateDeviceInfo(t *testing.T) {
	testDevID := "000000030000000000000001"
	ctx := context.Background()
	cli, _ := NewClient()

	current := time.Now()
	r := rand.New(rand.NewSource(current.Unix()))

	// InstallSession key should not be random variable because it has foreign key constraints.
	req := &pb.DeviceInfoUpdateRequest{
		Devid:          testDevID,
		Alias:          &pb.DeviceInfoUpdateRequest_AliasValue{AliasValue: fmt.Sprintf("alias-%d", r.Uint32())},
		GroupId:        &pb.DeviceInfoUpdateRequest_GroupIdValue{GroupIdValue: fmt.Sprintf("group-%d", r.Uint32())},
		Latitude:       &pb.DeviceInfoUpdateRequest_LatitudeValue{LatitudeValue: r.Float64()},
		Longitude:      &pb.DeviceInfoUpdateRequest_LongitudeValue{LongitudeValue: r.Float64()},
		Installer:      &pb.DeviceInfoUpdateRequest_InstallerValue{InstallerValue: fmt.Sprintf("installer-%d", r.Uint32())},
		InstallDate:    &pb.DeviceInfoUpdateRequest_InstallDateValue{InstallDateValue: &timestamp.Timestamp{Seconds: current.Unix(), Nanos: 0}},
		DevType:        &pb.DeviceInfoUpdateRequest_DevTypeValue{DevTypeValue: pb.DeviceType(r.Int() % 3)},
		AppFwVer:       &pb.DeviceInfoUpdateRequest_AppFwVerValue{AppFwVerValue: fmt.Sprintf("app-%d", r.Uint32())},
		LoraFwVer:      &pb.DeviceInfoUpdateRequest_LoraFwVerValue{LoraFwVerValue: fmt.Sprintf("lora-%d", r.Uint32())},
		RecogType:      &pb.DeviceInfoUpdateRequest_RecogTypeValue{RecogTypeValue: pb.RecogType(r.Int() % 2)},
		InstallSession: &pb.DeviceInfoUpdateRequest_InstallSessionValue{InstallSessionValue: "fff75f20fc45a6fa0fb8445218b5a700244c93e2a7c7274e5f063e39fef230fd"},
	}

	resp, err := cli.UpdateInfo(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, pb.ResponseCode_SUCCESS, resp.ResultCode)

	resp, err = cli.Detail(ctx, testDevID)

	assert.Nil(t, err)
	assert.Equal(t, req.Devid, resp.Devices[0].Devid)
	assert.Equal(t, req.GetAliasValue(), resp.Devices[0].Alias)
	assert.Equal(t, req.GetGroupIdValue(), resp.Devices[0].GroupId)
	assert.Equal(t, req.GetLatitudeValue(), resp.Devices[0].Latitude)
	assert.Equal(t, req.GetLongitudeValue(), resp.Devices[0].Longitude)
	assert.Equal(t, req.GetInstallerValue(), resp.Devices[0].Installer)
	assert.Equal(t, req.GetInstallDateValue().Seconds, resp.Devices[0].InstallDate.Seconds)
	assert.Equal(t, req.GetDevTypeValue(), resp.Devices[0].DevType)
	assert.Equal(t, req.GetAppFwVerValue(), resp.Devices[0].AppFwVer)
	assert.Equal(t, req.GetLoraFwVerValue(), resp.Devices[0].LoraFwVer)
	assert.Equal(t, req.GetRecogTypeValue(), resp.Devices[0].RecogType)
	assert.Equal(t, req.GetInstallSessionValue(), resp.Devices[0].InstallSessionKey)
}

func TestUpdateDeviceStatus(t *testing.T) {
	testDevid := "000000030000000000000001"

	ctx := context.Background()
	cli, _ := NewClient()

	current := time.Now()

	reqs := []*pb.DeviceStatusUpdateRequest{
		{
			Devid:       testDevid,
			Battery:     &pb.DeviceStatusUpdateRequest_BatteryValue{BatteryValue: 100},
			Temperature: &pb.DeviceStatusUpdateRequest_TemperatureValue{TemperatureValue: 20},
			Rssi:        &pb.DeviceStatusUpdateRequest_RssiValue{RssiValue: -120},
			AccX:        &pb.DeviceStatusUpdateRequest_AccXMgValue{AccXMgValue: 1000.0},
			AccY:        &pb.DeviceStatusUpdateRequest_AccYMgValue{AccYMgValue: 0.0},
			AccZ:        &pb.DeviceStatusUpdateRequest_AccZMgValue{AccZMgValue: 0.0},
			IsDeviceOk:  &pb.DeviceStatusUpdateRequest_IsDeviceOkValue{IsDeviceOkValue: true},
			AlarmStatus: &pb.DeviceStatusUpdateRequest_IsAlarmedValue{IsAlarmedValue: true},
			AlarmDate:   &pb.DeviceStatusUpdateRequest_AlarmDateValue{AlarmDateValue: &timestamp.Timestamp{Seconds: current.Unix()}},
		},
		{
			Devid:       testDevid,
			Battery:     &pb.DeviceStatusUpdateRequest_BatteryValue{BatteryValue: 10},
			Temperature: &pb.DeviceStatusUpdateRequest_TemperatureValue{TemperatureValue: -10},
			Rssi:        &pb.DeviceStatusUpdateRequest_RssiValue{RssiValue: -100},
			AccX:        &pb.DeviceStatusUpdateRequest_AccXMgValue{AccXMgValue: -900.0},
			AccY:        &pb.DeviceStatusUpdateRequest_AccYMgValue{AccYMgValue: 100.0},
			AccZ:        &pb.DeviceStatusUpdateRequest_AccZMgValue{AccZMgValue: -200.123},
			IsDeviceOk:  &pb.DeviceStatusUpdateRequest_IsDeviceOkValue{IsDeviceOkValue: false},
			AlarmStatus: &pb.DeviceStatusUpdateRequest_IsAlarmedValue{IsAlarmedValue: false},
			AlarmDate:   &pb.DeviceStatusUpdateRequest_AlarmDateValue{AlarmDateValue: &timestamp.Timestamp{Seconds: 0, Nanos: 0}},
		},
	}

	for _, req := range reqs {
		resp, err := cli.UpdateStatus(ctx, req)

		assert.Nil(t, err)
		assert.Equal(t, pb.ResponseCode_SUCCESS, resp.ResultCode)
		assert.Equal(t, testDevid, resp.Devices[0].Devid)
		assert.Equal(t, req.GetBatteryValue(), resp.Devices[0].Battery)
		assert.Equal(t, req.GetTemperatureValue(), resp.Devices[0].Temperature)
		assert.Equal(t, req.GetRssiValue(), resp.Devices[0].Rssi)
		assert.Equal(t, req.GetAccXMgValue(), resp.Devices[0].AccXMg)
		assert.Equal(t, req.GetAccYMgValue(), resp.Devices[0].AccYMg)
		assert.Equal(t, req.GetAccZMgValue(), resp.Devices[0].AccZMg)
		assert.Equal(t, req.GetIsDeviceOkValue(), resp.Devices[0].IsDeviceOk)
		assert.Equal(t, req.GetIsAlarmedValue(), resp.Devices[0].IsAlarmed)

		alarmDate := req.GetAlarmDateValue()
		if alarmDate.Seconds == 0 && alarmDate.Nanos == 0 {
			assert.Nil(t, resp.Devices[0].AlarmDate)
		} else {
			assert.Equal(t, alarmDate.Seconds, resp.Devices[0].AlarmDate.Seconds)
		}
	}

}

func TestUpdateDeviceConfig(t *testing.T) {
	testDevid := "000000030000000000000001"
	ctx := context.Background()
	cli, _ := NewClient()

	reqs := []*pb.DeviceConfigUpdateRequest{
		{
			Devid:             testDevid,
			SensorRange:       &pb.DeviceConfigUpdateRequest_SensorRangeValue{SensorRangeValue: pb.SensorRangeType_Gravity2},
			IntThreshold:      &pb.DeviceConfigUpdateRequest_IntThresholdValue{IntThresholdValue: 1080.0},
			DecisionThreshold: &pb.DeviceConfigUpdateRequest_DecisionThresholdValue{DecisionThresholdValue: 1100.0},
			SampleRate:        &pb.DeviceConfigUpdateRequest_SampleRateValue{SampleRateValue: 100},
			WaveBlocks:        &pb.DeviceConfigUpdateRequest_WaveBlocksValue{WaveBlocksValue: 12},
			IsNotifEnabled:    &pb.DeviceConfigUpdateRequest_IsNotifEnabledValue{IsNotifEnabledValue: true},
			RecogParam_0:      &pb.DeviceConfigUpdateRequest_RecogParam_0Value{RecogParam_0Value: 12.0},
			RecogParam_1:      &pb.DeviceConfigUpdateRequest_RecogParam_1Value{RecogParam_1Value: 0.6},
			RecogParam_2:      &pb.DeviceConfigUpdateRequest_RecogParam_2Value{RecogParam_2Value: 8.0},
			MuteDate: &pb.DeviceConfigUpdateRequest_MuteDateValue{
				MuteDateValue: &timestamp.Timestamp{Seconds: time.Now().Unix()}},
		},
		{
			Devid:             testDevid,
			SensorRange:       &pb.DeviceConfigUpdateRequest_SensorRangeValue{SensorRangeValue: pb.SensorRangeType_Gravity16},
			IntThreshold:      &pb.DeviceConfigUpdateRequest_IntThresholdValue{IntThresholdValue: 1900.0},
			DecisionThreshold: &pb.DeviceConfigUpdateRequest_DecisionThresholdValue{DecisionThresholdValue: 1800.0},
			SampleRate:        &pb.DeviceConfigUpdateRequest_SampleRateValue{SampleRateValue: 200},
			WaveBlocks:        &pb.DeviceConfigUpdateRequest_WaveBlocksValue{WaveBlocksValue: 2},
			IsNotifEnabled:    &pb.DeviceConfigUpdateRequest_IsNotifEnabledValue{IsNotifEnabledValue: false},
			RecogParam_0:      &pb.DeviceConfigUpdateRequest_RecogParam_0Value{RecogParam_0Value: 10.0},
			RecogParam_1:      &pb.DeviceConfigUpdateRequest_RecogParam_1Value{RecogParam_1Value: 0.3},
			RecogParam_2:      &pb.DeviceConfigUpdateRequest_RecogParam_2Value{RecogParam_2Value: 6.0},
			MuteDate:          &pb.DeviceConfigUpdateRequest_MuteDateValue{MuteDateValue: nil},
		},
	}

	for _, req := range reqs {
		resp, err := cli.UpdateConfig(ctx, req)

		assert.Nil(t, err)
		assert.Equal(t, pb.ResponseCode_SUCCESS, resp.ResultCode)
		assert.Equal(t, testDevid, resp.Devices[0].Devid)
		assert.Equal(t, req.GetSensorRangeValue(), resp.Devices[0].SensorRange)
		assert.Equal(t, req.GetIntThresholdValue(), resp.Devices[0].IntThresholdMg)

		assert.Equal(t, req.GetSampleRateValue(), resp.Devices[0].SampleRate)
		assert.Equal(t, req.GetWaveBlocksValue(), resp.Devices[0].WaveBlocks)
		assert.Equal(t, req.GetIsNotifEnabledValue(), resp.Devices[0].IsNotifEnabled)
		assert.Equal(t, req.GetRecogParam_0Value(), resp.Devices[0].RecogParam_0)
		assert.Equal(t, req.GetRecogParam_1Value(), resp.Devices[0].RecogParam_1)
		assert.Equal(t, req.GetRecogParam_2Value(), resp.Devices[0].RecogParam_2)

		devResp, _ := cli.Detail(ctx, testDevid)
		dev := devResp.Devices[0]

		if req.GetMuteDateValue() != nil {
			assert.Equal(t, req.GetMuteDateValue().Seconds, dev.MuteDate.Seconds)
		} else {
			assert.Nil(t, dev.GetMuteDate())
		}
	}
}

func TestGetStatusLog(t *testing.T) {
	tests := []struct {
		Desc       string
		DevID      string
		InstallKey string
		ExpectErr  error
	}{
		{
			Desc:       "Successful",
			DevID:      "customer_test_n",
			InstallKey: "d51a5d903d75786a9ea80f4f6d7ce58034feddc6562960cad791b464af419107",
			ExpectErr:  nil,
		},
		{
			Desc:       "Non-exist device.",
			DevID:      "non-exist-device",
			InstallKey: "d51a5d903d75786a9ea80f4f6d7ce58034feddc6562960cad791b464af419107",
			ExpectErr:  ErrNonExistDevice,
		},
		{
			Desc:       "Non-exist install key",
			DevID:      "customer_test_n",
			InstallKey: "dummy",
			ExpectErr:  nil,
		},
	}

	ctx := context.Background()
	cli, _ := NewClient()

	timeTo := time.Now()
	timeFrom := timeTo.AddDate(0, 0, -7)

	for _, test := range tests {
		logs, err := cli.StatusLog(ctx, test.DevID, test.InstallKey, timeFrom, timeTo, 0, 100)

		if test.ExpectErr != nil {
			assert.Equal(t, test.ExpectErr, err)
			continue
		}

		assert.True(t, len(logs) <= 100)
		for _, log := range logs {
			assert.Equal(t, test.DevID, log.Devid)
			assert.Equal(t, test.InstallKey, log.InstallSessionKey)
			assert.True(t, log.Time.After(timeFrom))
			assert.True(t, log.Time.Before(timeTo))
		}
	}
}

func TestStoreStatusLog(t *testing.T) {
	tests := []struct {
		Desc        string
		DevID       string
		Battery     int
		Temperature int
		RSSI        int
		ExpectErr   error
	}{
		{
			DevID:       "customer_test_n",
			Battery:     100,
			Temperature: 20,
			RSSI:        -120,
			ExpectErr:   nil,
		},
		{
			DevID:       "non-exist-device",
			Battery:     10,
			Temperature: -10,
			RSSI:        -100,
			ExpectErr:   ErrNonExistDevice,
		},
		{
			DevID:       "000000030000000000000001",
			Battery:     40,
			Temperature: 2,
			RSSI:        -80,
			ExpectErr:   ErrForbiddenInstallStatus,
		},
	}

	ctx := context.Background()
	cli, _ := NewClient()

	for _, test := range tests {
		err := cli.StoreStatusLog(ctx, test.DevID, test.Battery, test.Temperature, test.RSSI)

		if test.ExpectErr != nil {
			assert.Equal(t, test.ExpectErr, err)
			continue
		}

		time.Sleep(time.Second)

		devResp, _ := cli.Detail(ctx, test.DevID)
		device := devResp.Devices[0]

		installKey := device.InstallSessionKey

		current := time.Now()

		logs, err := cli.StatusLog(ctx, test.DevID, installKey, current.Add(-time.Second*60), current, 0, 1)

		assert.Equal(t, test.DevID, logs[0].Devid)
		assert.InDelta(t, current.Unix(), logs[0].Time.Unix(), 5)
		assert.Equal(t, test.Battery, logs[0].Battery)
		assert.Equal(t, test.Temperature, logs[0].Temperature)
		assert.Equal(t, test.RSSI, logs[0].RSSI)
	}
}

func TestInstall(t *testing.T) {
	testDevid := "000000030000000000000001"
	ctx := context.Background()
	cli, _ := NewClient()

	_, _ = cli.UpdateInfo(ctx, &pb.DeviceInfoUpdateRequest{
		Devid:     testDevid,
		Alias:     &pb.DeviceInfoUpdateRequest_AliasValue{AliasValue: ""},
		Latitude:  &pb.DeviceInfoUpdateRequest_LatitudeValue{LatitudeValue: 0},
		Longitude: &pb.DeviceInfoUpdateRequest_LongitudeValue{LongitudeValue: 0},
		Installer: &pb.DeviceInfoUpdateRequest_InstallerValue{InstallerValue: ""},
	})
	_, _ = cli.UpdateStatus(ctx, &pb.DeviceStatusUpdateRequest{
		Devid:         testDevid,
		InstallStatus: &pb.DeviceStatusUpdateRequest_InstallStatusValue{InstallStatusValue: pb.InstallStatus_Initial},
	})

	req := &pb.PrepareInstallRequest{
		Devid:     testDevid,
		Alias:     "test-alias",
		Latitude:  36.1,
		Longitude: 127.1,
		Installer: "contact@ino-on.com",
		GroupId:   "",
	}

	// First Test PrepareInstall
	current := time.Now()

	resp, err := cli.PrepareInstall(ctx, req)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, pb.ResponseCode_SUCCESS, resp.ResponseCode)
	assert.Equal(t, testDevid, resp.Devid)
	assert.NotEqual(t, "", resp.InstallSessionKey)

	devResp, _ := cli.Detail(ctx, testDevid)
	device := devResp.Devices[0]

	assert.Equal(t, testDevid, device.Devid)
	assert.Equal(t, pb.InstallStatus_Requested, device.InstallStatus)
	assert.Equal(t, req.Alias, device.Alias)
	assert.Equal(t, req.Latitude, device.Latitude)
	assert.Equal(t, req.Longitude, device.Longitude)
	assert.Equal(t, req.Installer, device.Installer)
	assert.InDelta(t, current.Unix(), device.InstallDate.AsTime().Unix(), 3)

	// Force change to WaitCompleteInstall status.
	cli.WaitCompleteInstall(ctx, &pb.WaitCompleteInstallRequest{Devid: testDevid})

	// Change status for rest testing.
	_, err = cli.UpdateStatus(ctx, &pb.DeviceStatusUpdateRequest{
		Devid:       testDevid,
		AlarmStatus: &pb.DeviceStatusUpdateRequest_IsAlarmedValue{IsAlarmedValue: true},
		AlarmDate: &pb.DeviceStatusUpdateRequest_AlarmDateValue{
			AlarmDateValue: &timestamp.Timestamp{Seconds: current.Unix()},
		},
	})

	_, err = cli.UpdateConfig(ctx, &pb.DeviceConfigUpdateRequest{
		Devid: testDevid,
		MuteDate: &pb.DeviceConfigUpdateRequest_MuteDateValue{
			MuteDateValue: &timestamp.Timestamp{Seconds: current.Unix()},
		},
	})

	// Second Test CompleteInstall
	completeResp, err := cli.CompleteInstall(ctx, &pb.CompleteInstallRequest{
		Devid:             testDevid,
		InstallSessionKey: resp.InstallSessionKey,
	})

	// Asserts.
	assert.Equal(t, pb.ResponseCode_SUCCESS, completeResp.ResponseCode)

	devResp, _ = cli.Detail(ctx, testDevid)
	device = devResp.Devices[0]

	assert.Equal(t, testDevid, device.Devid)
	assert.Equal(t, pb.InstallStatus_Installed, device.InstallStatus)
	assert.Equal(t, resp.InstallSessionKey, device.InstallSessionKey)
	assert.False(t, device.IsAlarmed)
	assert.Nil(t, device.AlarmDate)
	assert.Nil(t, device.MuteDate)

	// Check firestore.
	docRef := db.Doc(fmt.Sprintf("device/%s/install/%s", testDevid, resp.InstallSessionKey))
	doc, err := docRef.Get(ctx)

	assert.Nil(t, err)

	data := doc.Data()

	assert.Equal(t, req.Installer, data["installer"].(string))
	assert.Equal(t, req.GroupId, data["group_id"].(string))

	// Clear
	_, _ = cli.UpdateStatus(ctx, &pb.DeviceStatusUpdateRequest{
		Devid:         testDevid,
		InstallStatus: &pb.DeviceStatusUpdateRequest_InstallStatusValue{InstallStatusValue: pb.InstallStatus_Initial},
	})
}

func TestWaitInstallComplete(t *testing.T) {
	// Prepare
	testDevid := "000000030000000000000001"
	ctx := context.Background()
	cli, _ := NewClient()

	prepareReq := &pb.PrepareInstallRequest{
		Devid:     testDevid,
		Alias:     "test-alias",
		Latitude:  36.1,
		Longitude: 127.1,
		Installer: "contact@ino-on.com",
		GroupId:   "",
	}

	_, _ = cli.PrepareInstall(ctx, prepareReq)

	cli.UpdateStatus(ctx, &pb.DeviceStatusUpdateRequest{
		Devid:         testDevid,
		InstallStatus: &pb.DeviceStatusUpdateRequest_InstallStatusValue{InstallStatusValue: pb.InstallStatus_Requested},
	})

	// Test
	_, err := cli.WaitCompleteInstall(ctx, &pb.WaitCompleteInstallRequest{Devid: testDevid})

	// Asserts
	assert.Nil(t, err)

	devResp, _ := cli.Detail(ctx, testDevid)
	device := devResp.Devices[0]

	assert.Equal(t, pb.InstallStatus_WaitInstallComplete, device.InstallStatus)

	// Clear
	_, _ = cli.UpdateStatus(ctx, &pb.DeviceStatusUpdateRequest{
		Devid:         testDevid,
		InstallStatus: &pb.DeviceStatusUpdateRequest_InstallStatusValue{InstallStatusValue: pb.InstallStatus_Initial},
	})
}

func TestInstallCompleteOnOtherStatus(t *testing.T) {
	testDevid := "000000030000000000000001"

	ctx := context.Background()
	cli, _ := NewClient()

	prepareReq := &pb.PrepareInstallRequest{
		Devid:     testDevid,
		Alias:     "test-alias",
		Latitude:  36.1,
		Longitude: 127.1,
		Installer: "contact@ino-on.com",
		GroupId:   "",
	}

	installResp, _ := cli.PrepareInstall(ctx, prepareReq)

	// TODO: Ommit wait complete install sequence.
	resp, err := cli.CompleteInstall(ctx, &pb.CompleteInstallRequest{
		Devid:             testDevid,
		InstallSessionKey: installResp.InstallSessionKey,
	})

	assert.Nil(t, err)
	assert.Equal(t, pb.ResponseCode_SUCCESS, resp.ResponseCode)
}

func TestUninstalling(t *testing.T) {
	// Prepare
	testDevid := "000000030000000000000001"
	ctx := context.Background()
	cli, _ := NewClient()

	prepareReq := &pb.PrepareInstallRequest{
		Devid:     testDevid,
		Alias:     "test-alias",
		Latitude:  36.1,
		Longitude: 127.1,
		Installer: "contact@ino-on.com",
		GroupId:   "",
	}

	resp, _ := cli.PrepareInstall(ctx, prepareReq)
	cli.WaitCompleteInstall(ctx, &pb.WaitCompleteInstallRequest{Devid: testDevid})
	_, _ = cli.CompleteInstall(ctx, &pb.CompleteInstallRequest{Devid: testDevid, InstallSessionKey: resp.InstallSessionKey})

	// Test
	_, err := cli.Uninstalling(ctx, &pb.UninstallingRequest{Devid: testDevid})

	// Asserts
	assert.Nil(t, err)

	devResp, _ := cli.Detail(ctx, testDevid)
	device := devResp.Devices[0]

	assert.Equal(t, pb.InstallStatus_Uninstalling, device.InstallStatus)

	// Clear
	_, _ = cli.UpdateStatus(ctx, &pb.DeviceStatusUpdateRequest{
		Devid:         testDevid,
		InstallStatus: &pb.DeviceStatusUpdateRequest_InstallStatusValue{InstallStatusValue: pb.InstallStatus_Initial},
	})
}

func TestUninstall(t *testing.T) {
	// Prepare
	testDevid := "000000030000000000000001"
	ctx := context.Background()
	cli, _ := NewClient()

	prepareReq := &pb.PrepareInstallRequest{
		Devid:     testDevid,
		Alias:     "test-alias",
		Latitude:  36.1,
		Longitude: 127.1,
		Installer: "contact@ino-on.com",
		GroupId:   "",
	}

	resp, _ := cli.PrepareInstall(ctx, prepareReq)
	cli.WaitCompleteInstall(ctx, &pb.WaitCompleteInstallRequest{Devid: testDevid})
	_, _ = cli.CompleteInstall(ctx, &pb.CompleteInstallRequest{Devid: testDevid, InstallSessionKey: resp.InstallSessionKey})

	// Test
	_, err := cli.Uninstall(ctx, &pb.UninstallRequest{Devid: testDevid})

	// Asserts
	assert.Nil(t, err)

	devResp, _ := cli.Detail(ctx, testDevid)
	device := devResp.Devices[0]

	assert.Equal(t, pb.InstallStatus_Initial, device.InstallStatus)
	assert.Equal(t, "", device.Alias)
	assert.Equal(t, float64(0), device.Latitude)
	assert.Equal(t, float64(0), device.Longitude)
	assert.Equal(t, "", device.Installer)
	assert.Nil(t, device.InstallDate)

	// Clear
	_, _ = cli.UpdateStatus(ctx, &pb.DeviceStatusUpdateRequest{
		Devid:         testDevid,
		InstallStatus: &pb.DeviceStatusUpdateRequest_InstallStatusValue{InstallStatusValue: pb.InstallStatus_Initial},
	})
}

func TestUninstallOnUninstalling(t *testing.T) {
	// Prepare
	testDevid := "000000030000000000000001"
	ctx := context.Background()
	cli, _ := NewClient()

	prepareReq := &pb.PrepareInstallRequest{
		Devid:     testDevid,
		Alias:     "test-alias",
		Latitude:  36.1,
		Longitude: 127.1,
		Installer: "contact@ino-on.com",
		GroupId:   "",
	}

	resp, _ := cli.PrepareInstall(ctx, prepareReq)
	cli.WaitCompleteInstall(ctx, &pb.WaitCompleteInstallRequest{Devid: testDevid})
	_, _ = cli.CompleteInstall(ctx, &pb.CompleteInstallRequest{Devid: testDevid, InstallSessionKey: resp.InstallSessionKey})
	_, _ = cli.Uninstalling(ctx, &pb.UninstallingRequest{Devid: testDevid})

	// Test
	_, err := cli.Uninstall(ctx, &pb.UninstallRequest{Devid: testDevid})

	// Asserts
	assert.Nil(t, err)

	devResp, _ := cli.Detail(ctx, testDevid)
	device := devResp.Devices[0]

	assert.Equal(t, pb.InstallStatus_Initial, device.InstallStatus)
	assert.Equal(t, "", device.Alias)
	assert.Equal(t, float64(0), device.Latitude)
	assert.Equal(t, float64(0), device.Longitude)
	assert.Equal(t, "", device.Installer)
	assert.Nil(t, device.InstallDate)

	// Clear
	_, _ = cli.UpdateStatus(ctx, &pb.DeviceStatusUpdateRequest{
		Devid:         testDevid,
		InstallStatus: &pb.DeviceStatusUpdateRequest_InstallStatusValue{InstallStatusValue: pb.InstallStatus_Initial},
	})
}

func TestDiscard(t *testing.T) {
	// Prepare
	testDevid := "000000030000000000000001"
	ctx := context.Background()
	cli, _ := NewClient()

	prepareReq := &pb.PrepareInstallRequest{
		Devid:     testDevid,
		Alias:     "test-alias",
		Latitude:  36.1,
		Longitude: 127.1,
		Installer: "contact@ino-on.com",
		GroupId:   "",
	}

	resp, _ := cli.PrepareInstall(ctx, prepareReq)
	cli.WaitCompleteInstall(ctx, &pb.WaitCompleteInstallRequest{Devid: testDevid})
	_, _ = cli.CompleteInstall(ctx, &pb.CompleteInstallRequest{Devid: testDevid, InstallSessionKey: resp.InstallSessionKey})

	// Test
	_, err := cli.Discard(ctx, &pb.DiscardRequest{Devid: testDevid})

	// Asserts
	assert.Nil(t, err)

	devResp, _ := cli.Detail(ctx, testDevid)
	device := devResp.Devices[0]

	assert.Equal(t, pb.InstallStatus_Discarded, device.InstallStatus)

	// Clear
	_, _ = cli.UpdateStatus(ctx, &pb.DeviceStatusUpdateRequest{
		Devid:         testDevid,
		InstallStatus: &pb.DeviceStatusUpdateRequest_InstallStatusValue{InstallStatusValue: pb.InstallStatus_Initial},
	})
}

func TestDiscardOnOtherStatus(t *testing.T) {
	// Prepare
	testDevid := "000000030000000000000001"
	ctx := context.Background()
	cli, _ := NewClient()

	prepareReq := &pb.PrepareInstallRequest{
		Devid:     testDevid,
		Alias:     "test-alias",
		Latitude:  36.1,
		Longitude: 127.1,
		Installer: "contact@ino-on.com",
		GroupId:   "",
	}

	// Intensionally call only PrepareInstall for testing on "requested" state.
	_, _ = cli.PrepareInstall(ctx, prepareReq)

	// Test
	_, err := cli.Discard(ctx, &pb.DiscardRequest{Devid: testDevid})

	// Asserts
	assert.Nil(t, err)

	devResp, _ := cli.Detail(ctx, testDevid)
	device := devResp.Devices[0]

	assert.Equal(t, pb.InstallStatus_Discarded, device.InstallStatus)

	// Clear
	_, _ = cli.UpdateStatus(ctx, &pb.DeviceStatusUpdateRequest{
		Devid:         testDevid,
		InstallStatus: &pb.DeviceStatusUpdateRequest_InstallStatusValue{InstallStatusValue: pb.InstallStatus_Initial},
	})
}

func TestLastInclinationLog(t *testing.T) {
	ctx := context.Background()
	cli, _ := NewClient()

	tests := []struct {
		Desc      string
		DevID     string
		ExpectErr error
	}{
		{
			Desc:      "Successful",
			DevID:     "customer_test_n",
			ExpectErr: nil,
		},
		{
			Desc:      "Non exist device",
			DevID:     "non-exist-device",
			ExpectErr: ErrNonExistDevice,
		},
		{
			Desc:      "No inclinations",
			DevID:     "000000030000000000000001",
			ExpectErr: ErrNoEntities,
		},
	}

	for _, test := range tests {
		log, err := cli.LastInclinationLog(ctx, test.DevID)

		assert.Equal(t, test.ExpectErr, err)
		fmt.Println(log, err)

		if err == nil {
			resp, _ := cli.Detail(ctx, test.DevID)
			dev := resp.Devices[0]

			assert.Equal(t, test.DevID, log.Devid)
			assert.Equal(t, dev.InstallSessionKey, log.InstallSessionKey)
		}
	}
}

func TestStoreInclination(t *testing.T) {
	ctx := context.Background()
	cli, _ := NewClient()

	tests := []struct {
		Desc      string
		DevID     string
		RawX      int
		RawY      int
		RawZ      int
		ExpectErr error
	}{
		{
			Desc:      "Acc. values are wrong.",
			DevID:     "customer_test_n",
			RawX:      0,
			RawY:      0,
			RawZ:      0,
			ExpectErr: ErrInvalidInclinationValue,
		},
		{
			Desc:      "Successful",
			DevID:     "customer_test_n",
			RawX:      0,
			RawY:      0,
			RawZ:      998.0,
			ExpectErr: nil,
		},
		{
			Desc:      "Store request on non-exist device",
			DevID:     "non-exist-device",
			RawX:      0,
			RawY:      0,
			RawZ:      998.0,
			ExpectErr: ErrNonExistDevice,
		},
		{
			Desc:      "Store request on not installed device",
			DevID:     "000000030000000000000001",
			RawX:      0,
			RawY:      0,
			RawZ:      998.0,
			ExpectErr: ErrForbiddenInstallStatus,
		},
	}

	for _, test := range tests {
		_, err := cli.StoreInclinationLog(ctx, test.DevID, test.RawX, test.RawY, test.RawZ)

		assert.Equal(t, test.ExpectErr, err)

		if err == nil {
			// Check values.
			log, err := cli.LastInclinationLog(ctx, test.DevID)

			assert.Nil(t, err)
			assert.Equal(t, test.DevID, log.Devid)
			assert.Equal(t, float64(test.RawX)*0.244, log.AccXMg)
			assert.Equal(t, float64(test.RawY)*0.244, log.AccYMg)
			assert.Equal(t, float64(test.RawZ)*0.244, log.AccZMg)
		}
	}
}

func TestAngle(t *testing.T) {
	tests := []struct {
		X float64
		Y float64
		Z float64
	}{
		{
			X: 0.0,
			Y: 0.0,
			Z: 998.0,
		},

		{
			X: 993.812,
			Y: -13.908,
			Z: -1.22,
		},

		{
			X: 993.812,
			Y: -13.908,
			Z: 1.22,
		},

		{
			X: 13.908,
			Y: 993.812,
			Z: 1.22,
		},
		{
			X: -13.908,
			Y: 993.812,
			Z: 1.22,
		},
	}

	for i, test := range tests {
		angleZ := angle(test.X, test.Y, test.Z, 1)
		fmt.Println("Angle", i, angleZ)
	}

}
