package user

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	target := os.Getenv("TEST_TARGET")
	if target != "" {
		serverURL = target + "-" + serverURL
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

//
// Test Auth token
//
func TestInvalidToken(t *testing.T) {
	dummyTokens := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		"",
		"hello",
	}

	for _, dummyToken := range dummyTokens {
		oauthToken := &oauth2.Token{
			AccessToken: dummyToken,
			TokenType:   "bearer",
			Expiry:      time.Now().Add(24 * time.Hour),
		}
		cli := &client{oauthToken: oauthToken}

		_, err := cli.GetDeviceToken("jkkim@ino-on.com")

		assert.NotEqual(t, codes.Unavailable, status.Code(err))
		assert.NotNil(t, err)
	}
}
