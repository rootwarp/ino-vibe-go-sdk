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
		serverURL = "grpc-dev.ino-vibe.ino-on.dev:443"

		if target == "feature" {
			serverURL = target + "-" + serverURL
		}
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
