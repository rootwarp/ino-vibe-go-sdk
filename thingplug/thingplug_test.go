package thingplug

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDevID = "00000125d02544fffefe143a"
	cli       Client
)

func init() {
	serverURL = "feature-thingplug.ino-vibe.ino-on.dev:443"

	cli, _ = NewClient()
}

func TestThingplugPowrOff(t *testing.T) {
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
