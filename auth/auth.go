package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
)

type storedCredential struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpireIn    int    `json:"expire_in"`
}

var (
	credFilePath string
)

func init() {
	credFilePath = os.Getenv("INOVIBE_APPLICATION_CREDENTIALS")
}

// IssueCredentials requests new credentials to Auth0.
func IssueCredentials(id, secret, audience string) (*oauth2.Token, error) {
	return nil, nil
}

// LoadCredentials loads credential file from local storage.
func LoadCredentials() (*oauth2.Token, error) {
	storedCred := storedCredential{}

	credData, err := ioutil.ReadFile(credFilePath)
	if err != nil {
		log.Println("LoadCredentials", err)
		return nil, err
	}

	err = json.Unmarshal(credData, &storedCred)
	if err != nil {
		log.Println("LoadCredentials", err)
		return nil, err
	}

	oauthToken := &oauth2.Token{
		AccessToken: storedCred.AccessToken,
		TokenType:   storedCred.TokenType,
	}

	return oauthToken, nil
}
