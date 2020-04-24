package loadtest

import (
	"log"
	"os"
	"testing"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	iv_auth "github.com/rootwarp/ino-vibe-go-sdk/auth"
)

func TestLoadUser(t *testing.T) {
	token, err := iv_auth.LoadCredentials()

	report, err := runner.Run(
		"inovibe.api.v3.UserService.RegisterDeviceToken",
		"user.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"username": "dummy@ino-on.com", "user_id": "dummy", "device_token": "asdf"}`),
		runner.WithMetadata(&map[string]string{"authorization": "bearer " + token.AccessToken}),
		runner.WithConcurrency(5),
		runner.WithQPS(20),
	)

	if err != nil {
		log.Fatalln(err)
	}

	printer := printer.ReportPrinter{
		Out:    os.Stdout,
		Report: report,
	}

	printer.Print("summary")
}
