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
		serverURL = fmt.Sprintf("%s-%s", target, serverURL)
	}

	fmt.Println("Test ", serverURL)
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
