package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

	if credFilePath == "" {
		home := os.Getenv("HOME")
		credFilePath = home + "/.inovibe/credentials.json"
	}

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

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(storedCred.AccessToken, claims,
		func(token *jwt.Token) (interface{}, error) {
			return nil, nil
		})

	expTime := time.Unix(int64(claims["exp"].(float64)), 0)

	oauthToken := &oauth2.Token{
		AccessToken: storedCred.AccessToken,
		TokenType:   storedCred.TokenType,
		Expiry:      expTime,
	}

	return oauthToken, nil
}

// IssueToken issues new OAuth2 token.
func IssueToken(clientID, clientSecret, audience string) (*oauth2.Token, error) {
	// Issue token
	cred := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      audience,
		"grant_type":    "client_credentials",
	}

	credData, err := json.Marshal(cred)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://ino-vibe.auth0.com/oauth/token",
		bytes.NewReader(credData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Auth failed")
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	newCred := storedCredential{}
	err = json.Unmarshal(respData, &newCred)
	if err != nil {
		return nil, err
	}

	oauthToken := &oauth2.Token{
		AccessToken: newCred.AccessToken,
		TokenType:   newCred.TokenType,
	}

	// Store
	u, err := user.Current()
	_ = os.Mkdir(u.HomeDir+"/.inovibe", 0744)

	err = ioutil.WriteFile(u.HomeDir+"/.inovibe/credentials.json", respData, 0644)
	if err != nil {
		return nil, err
	}

	return oauthToken, nil
}

// IsValidToken checks wheather received token is valid or not.
func IsValidToken(token *oauth2.Token) bool {
	current := time.Now()
	return token.Expiry.After(current)
}
