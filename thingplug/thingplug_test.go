package thingplug

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDevID = "00000125d02544fffefe143a"
	cli       Client
)

func init() {
	target := os.Getenv("TEST_TARGET")
	if target != "" {
		serverURL = target + "-" + serverURL
	}
	fmt.Println(serverURL)

	cli, _ = NewClient()
}

func TestThingplugPowerOff(t *testing.T) {
	ctx := context.Background()
	err := cli.PowerOff(ctx, testDevID)

	assert.Nil(t, err)
}

func TestThingplugReset(t *testing.T) {
	ctx := context.Background()
	err := cli.Reset(ctx, testDevID)

	assert.Nil(t, err)
}

func TestThingplugBaseReset(t *testing.T) {
	ctx := context.Background()
	err := cli.BaseReset(ctx, testDevID)

	assert.Nil(t, err)
}

func TestRegisterSubscription(t *testing.T) {

	tests := []struct {
		Desc      string
		DeviceID  string
		ExpectErr error
	}{
		{
			Desc:      "Success",
			DeviceID:  testDevID,
			ExpectErr: nil,
		},

		{
			Desc:      "Not exist",
			DeviceID:  "non-exist-device-id",
			ExpectErr: ErrNotExistDevice,
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		err := cli.RegisterSubscription(ctx, test.DeviceID)

		assert.Equal(t, test.ExpectErr, err)
	}
}
