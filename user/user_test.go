package user

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
}

func TestUserRegisterDeviceToken(t *testing.T) {
	cli, err := NewClient()

	err = cli.RegisterDeviceToken("dummy", "dummy@ino-on.com", "hello world")

	assert.Nil(t, err)
}

func TestUserRegisterDeviceTokenUnauthorized(t *testing.T) {
	cli, err := NewClient()
	(cli.(*client)).oauthToken.AccessToken = "invalid-token"

	err = cli.RegisterDeviceToken("dummy", "dummy@ino-on.com", "hello world")

	assert.NotNil(t, err)
}

func TestUserGetDeviceToken(t *testing.T) {
	cli, err := NewClient()

	tokens, err := cli.GetDeviceToken("jkkim@ino-on.com")

	assert.Nil(t, err)
	assert.True(t, len(tokens) > 0)
}

func TestUserGetDeviceTokenNonExistUser(t *testing.T) {
	cli, err := NewClient()

	tokens, err := cli.GetDeviceToken("non-exist@ino-on.com")

	assert.Nil(t, err)
	assert.True(t, len(tokens) == 0)
}
