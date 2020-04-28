package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRegisterDeviceToken(t *testing.T) {
	cli, err := NewClient()

	err = cli.RegisterDeviceToken("dummy", "dummy@ino-on.com", "hello world")

	assert.Nil(t, err)
}

func TestUserRegisterDeviceTokenUnauthorized(t *testing.T) {
	cli, err := NewClient()
	cli.oauthToken.AccessToken = "invalid-token"

	err = cli.RegisterDeviceToken("dummy", "dummy@ino-on.com", "hello world")

	assert.NotNil(t, err)
}
