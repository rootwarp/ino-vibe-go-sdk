package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
)

func main() {
	/*
		r := bufio.NewReader(os.Stdin)
		fmt.Println("Client ID? ")
		clientID := r.ReadString("\n")

		fmt.Println("Client Secret? ")
		clientSecret := r.ReadString("\n")

		fmt.Println("Audience? ")
		aud := r.ReadString("\n")
	*/

	clientID := "O9so4gOpXmnC6pUHc5rOeslkUA2bXgLK"
	clientSecret := "AV5xpMVFt93uvvPHsBHvB8nJERnamLYxOkBreqWptRSEDcS8QDUmflgMPQVVR5Hv"
	aud := "https://grpc.ino-vibe.ino-on.dev"

	credBody := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      aud,
		"grant_type":    "client_credentials",
	}

	credData, err := json.Marshal(credBody)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(
		"https://ino-vibe.auth0.com/oauth/token",
		"application/json",
		bytes.NewReader(credData),
	)

	u, err := user.Current()

	_ = os.Mkdir(u.HomeDir+"/.inovibe", 0744)

	respData, err := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile(u.HomeDir+"/.inovibe/credentials.json", respData, 0644)
	fmt.Println(err)
}
