package alert

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	pb "bitbucket.org/ino-on/ino-vibe-api"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
)

var (
	testDevID = "00000125d02544fffefe108a"
)

func init() {
	target := os.Getenv("TEST_TARGET")
	if target != "" {
		serverURL = target + "-" + serverURL
	}
	fmt.Println(serverURL)
}

func TestListAlertByDeviceID(t *testing.T) {
	cli, _ := NewClient()
	defer cli.Close()

	// Tests
	requests := []*pb.AlertListRequest{
		{
			Search: &pb.AlertListRequest_Devid{Devid: testDevID},
		},
		{
			Search:   &pb.AlertListRequest_Devid{Devid: testDevID},
			MaxCount: 20,
		},
		{
			Search: &pb.AlertListRequest_Devid{Devid: testDevID},
			DateFrom: &pb.AlertListRequest_DateFromValue{
				DateFromValue: &timestamp.Timestamp{Seconds: (time.Now().Add(-(time.Hour * 24 * 7))).Unix()},
			},
			DateTo: &pb.AlertListRequest_DateToValue{
				DateToValue: &timestamp.Timestamp{Seconds: (time.Now().Add(-(time.Hour * 24 * 2))).Unix()},
			},
		},
	}

	for _, request := range requests {
		// Send request.
		resp, err := cli.List(context.Background(), request)

		assert.Nil(t, err)
		assert.Equal(t, pb.ResponseCode_SUCCESS, resp.ResponseCode)

		if request.MaxCount == 0 {
			// Default value matches to 10.
			assert.True(t, 10 >= len(resp.Alerts))
		} else {
			assert.Equal(t, int(request.MaxCount), len(resp.Alerts))
		}

		for _, alert := range resp.GetAlerts() {
			assert.Equal(t, request.GetDevid(), alert.Devid)

			if request.GetDateFrom() != nil {
				assert.True(t, alert.GetIssued().Seconds >= request.GetDateFromValue().Seconds)
			}

			if request.GetDateTo() != nil {
				assert.True(t, alert.GetIssued().Seconds <= request.GetDateToValue().Seconds)
			}
		}
	}
}

func TestListAlertByGroupID(t *testing.T) {
	cli, _ := NewClient()
	defer cli.Close()

	requests := []*pb.AlertListRequest{
		{
			Search: &pb.AlertListRequest_Groupid{Groupid: "0bee7b43-0b57-4b54-9062-430e2bd3fa79"},
		},
	}

	for _, request := range requests {
		resp, err := cli.List(context.Background(), request)

		assert.Nil(t, err)
		assert.Equal(t, 10, len(resp.Alerts))

		for _, alert := range resp.GetAlerts() {
			assert.Equal(t, request.GetGroupid(), alert.AlarmGroup.GetGroupid())
		}
	}
}
