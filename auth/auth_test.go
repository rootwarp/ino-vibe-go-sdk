package auth

import (
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
	credFilePath = "./credentials.json"
	f, _ := os.Create("./credentials.json")
	_, err := f.Write([]byte(`{"access_token": "hello world", "token_type": "bearer"}`))
	if err != nil {
		panic(err)
	}
	f.Close()

	defer os.Remove(credFilePath)

	oauthToken, err := LoadCredentials()

	assert.Nil(t, err)
	assert.Equal(t, "hello world", oauthToken.AccessToken)
	assert.Equal(t, "bearer", oauthToken.TokenType)
}

// TODO: Issue credentials.
