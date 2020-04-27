package device

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"

	pb "bitbucket.org/ino-on/ino-vibe-api"
)

/*
func TestGetDeviceServiceVersion(t *testing.T) {
	version, err := GetDeviceServiceVersion()

	assert.Nil(t, err)
	fmt.Println(version)
	assert.NotEqual(t, version, "")
}
*/

func TestGetDeviceListUnauthorized(t *testing.T) {
	ctx := context.Background()
	cli, _ := NewClient()
	cli.oauthToken.AccessToken = "invalid-token"

	_, err := cli.List(ctx, pb.InstallStatus_Installed)

	assert.NotNil(t, err)
}

func TestGetDeviceListInstallStatus(t *testing.T) {
	ctx := context.Background()
	cli, _ := NewClient()
	resp, err := cli.List(ctx, pb.InstallStatus_Installed)

	assert.Nil(t, err)
	assert.Equal(t, pb.ResponseCode_SUCCESS, resp.ResultCode)

	for _, device := range resp.Devices {
		assert.Equal(t, pb.InstallStatus_Installed, device.InstallStatus)
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
	req := &pb.DeviceInfoUpdateRequest{
		Devid:       testDevID,
		Alias:       &pb.DeviceInfoUpdateRequest_AliasValue{fmt.Sprintf("alias-%d", r.Uint32())},
		GroupId:     &pb.DeviceInfoUpdateRequest_GroupIdValue{fmt.Sprintf("group-%d", r.Uint32())},
		Latitude:    &pb.DeviceInfoUpdateRequest_LatitudeValue{r.Float64()},
		Longitude:   &pb.DeviceInfoUpdateRequest_LongitudeValue{r.Float64()},
		Installer:   &pb.DeviceInfoUpdateRequest_InstallerValue{fmt.Sprintf("installer-%d", r.Uint32())},
		InstallDate: &pb.DeviceInfoUpdateRequest_InstallDateValue{&timestamp.Timestamp{Seconds: current.Unix(), Nanos: 0}},
		DevType:     &pb.DeviceInfoUpdateRequest_DevTypeValue{pb.DeviceType(r.Int() % 3)},
		AppFwVer:    &pb.DeviceInfoUpdateRequest_AppFwVerValue{fmt.Sprintf("app-%d", r.Uint32())},
		LoraFwVer:   &pb.DeviceInfoUpdateRequest_LoraFwVerValue{fmt.Sprintf("lora-%d", r.Uint32())},
		RecogType:   &pb.DeviceInfoUpdateRequest_RecogTypeValue{pb.RecogType(r.Int() % 2)},
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
}

func TestUpdateDeviceStatus(t *testing.T) {
	testDevid := "000000030000000000000001"

	ctx := context.Background()
	cli, _ := NewClient()

	current := time.Now()

	reqs := []*pb.DeviceStatusUpdateRequest{
		&pb.DeviceStatusUpdateRequest{
			Devid:       testDevid,
			Battery:     &pb.DeviceStatusUpdateRequest_BatteryValue{100},
			Temperature: &pb.DeviceStatusUpdateRequest_TemperatureValue{20},
			Rssi:        &pb.DeviceStatusUpdateRequest_RssiValue{-120},
			AccX:        &pb.DeviceStatusUpdateRequest_AccXMgValue{1000.0},
			AccY:        &pb.DeviceStatusUpdateRequest_AccYMgValue{0.0},
			AccZ:        &pb.DeviceStatusUpdateRequest_AccZMgValue{0.0},
			IsDeviceOk:  &pb.DeviceStatusUpdateRequest_IsDeviceOkValue{true},
			AlarmStatus: &pb.DeviceStatusUpdateRequest_IsAlarmedValue{true},
			AlarmDate:   &pb.DeviceStatusUpdateRequest_AlarmDateValue{&timestamp.Timestamp{Seconds: current.Unix()}},
		},
		&pb.DeviceStatusUpdateRequest{
			Devid:       testDevid,
			Battery:     &pb.DeviceStatusUpdateRequest_BatteryValue{10},
			Temperature: &pb.DeviceStatusUpdateRequest_TemperatureValue{-10},
			Rssi:        &pb.DeviceStatusUpdateRequest_RssiValue{-100},
			AccX:        &pb.DeviceStatusUpdateRequest_AccXMgValue{-900.0},
			AccY:        &pb.DeviceStatusUpdateRequest_AccYMgValue{100.0},
			AccZ:        &pb.DeviceStatusUpdateRequest_AccZMgValue{-200.123},
			IsDeviceOk:  &pb.DeviceStatusUpdateRequest_IsDeviceOkValue{false},
			AlarmStatus: &pb.DeviceStatusUpdateRequest_IsAlarmedValue{false},
			AlarmDate:   &pb.DeviceStatusUpdateRequest_AlarmDateValue{&timestamp.Timestamp{Seconds: 0, Nanos: 0}},
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
		&pb.DeviceConfigUpdateRequest{
			Devid:             testDevid,
			SensorRange:       &pb.DeviceConfigUpdateRequest_SensorRangeValue{pb.SensorRangeType_Gravity2},
			IntThreshold:      &pb.DeviceConfigUpdateRequest_IntThresholdValue{1080.0},
			DecisionThreshold: &pb.DeviceConfigUpdateRequest_DecisionThresholdValue{1100.0},
			SampleRate:        &pb.DeviceConfigUpdateRequest_SampleRateValue{100},
			WaveBlocks:        &pb.DeviceConfigUpdateRequest_WaveBlocksValue{12},
			IsNotifEnabled:    &pb.DeviceConfigUpdateRequest_IsNotifEnabledValue{true},
			RecogParam_0:      &pb.DeviceConfigUpdateRequest_RecogParam_0Value{12.0},
			RecogParam_1:      &pb.DeviceConfigUpdateRequest_RecogParam_1Value{0.6},
			RecogParam_2:      &pb.DeviceConfigUpdateRequest_RecogParam_2Value{8.0},
		},
		&pb.DeviceConfigUpdateRequest{
			Devid:             testDevid,
			SensorRange:       &pb.DeviceConfigUpdateRequest_SensorRangeValue{pb.SensorRangeType_Gravity16},
			IntThreshold:      &pb.DeviceConfigUpdateRequest_IntThresholdValue{1900.0},
			DecisionThreshold: &pb.DeviceConfigUpdateRequest_DecisionThresholdValue{1800.0},
			SampleRate:        &pb.DeviceConfigUpdateRequest_SampleRateValue{200},
			WaveBlocks:        &pb.DeviceConfigUpdateRequest_WaveBlocksValue{2},
			IsNotifEnabled:    &pb.DeviceConfigUpdateRequest_IsNotifEnabledValue{false},
			RecogParam_0:      &pb.DeviceConfigUpdateRequest_RecogParam_0Value{10.0},
			RecogParam_1:      &pb.DeviceConfigUpdateRequest_RecogParam_1Value{0.3},
			RecogParam_2:      &pb.DeviceConfigUpdateRequest_RecogParam_2Value{6.0},
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
	}
}
