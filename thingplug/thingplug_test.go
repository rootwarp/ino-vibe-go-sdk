package thingplug

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThingplugPowrOff(t *testing.T) {
	serverURL = "thingplug.ino-vibe.ino-on.dev:443"

	cli, err := NewClient()

	assert.Nil(t, err)

	ctx := context.Background()
	err = cli.PowerOff(ctx, "00000125d02544fffefe143a")

	assert.Nil(t, err)
}
