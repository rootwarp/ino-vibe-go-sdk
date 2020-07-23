package auth

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthNoCredentialFile(t *testing.T) {
	credFilePath = "./non-exist/credential.json"

	oauthToken, err := LoadCredentials()

	assert.Nil(t, oauthToken)
	assert.NotNil(t, err)
}

func TestAuthInvalidCredentialFormat(t *testing.T) {
	credFilePath = "./credentials.json"
	f, _ := os.Create("./credentials.json")
	_, err := f.Write([]byte("{hello world}"))
	if err != nil {
		panic(err)
	}
	f.Close()

	defer os.Remove(credFilePath)

	oauthToken, err := LoadCredentials()

	assert.Nil(t, oauthToken)
	assert.NotNil(t, err)
}

func TestAuthLoadSuccess(t *testing.T) {
	credFilePath = os.Getenv("INOVIBE_APPLICATION_CREDENTIALS")
	oauthToken, err := LoadCredentials()

	assert.NotNil(t, oauthToken)
	assert.Nil(t, err)
}

func TestAuthDefaultFilePath(t *testing.T) {
	credFilePath = ""
	_, _ = LoadCredentials()

	home := os.Getenv("HOME")
	assert.Equal(t, home+"/.inovibe/credentials.json", credFilePath)
}

func TestAuthIssueToken(t *testing.T) {
	clientID := "O9so4gOpXmnC6pUHc5rOeslkUA2bXgLK"
	clientSecret := "AV5xpMVFt93uvvPHsBHvB8nJERnamLYxOkBreqWptRSEDcS8QDUmflgMPQVVR5Hv"
	audience := "https://grpc.ino-vibe.ino-on.dev"

	token, err := IssueToken(clientID, clientSecret, audience)

	assert.NotNil(t, token)
	assert.Nil(t, err)
}

func TestAuthCheckTokenValid(t *testing.T) {
	token, _ := LoadCredentials()

	fmt.Println(token)

}
