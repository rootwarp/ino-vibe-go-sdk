package loadtest

import (
	"log"
	"os"
	"testing"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	iv_auth "github.com/rootwarp/ino-vibe-go-sdk/auth"
)

func TestLoadGroupDetail(t *testing.T) {
	token, err := iv_auth.LoadCredentials()

	report, err := runner.Run(
		"inovibe.api.v3.GroupService.Detail",
		"group.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"groupid": "607f9db4-7eee-4a08-894d-356c8a462ae1"}`),
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

func TestLoadGroupChilds(t *testing.T) {
	token, err := iv_auth.LoadCredentials()

	report, err := runner.Run(
		"inovibe.api.v3.GroupService.Childs",
		"group.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"groupid": "607f9db4-7eee-4a08-894d-356c8a462ae1"}`),
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

func TestLoadGroupNestedUsers(t *testing.T) {
	token, err := iv_auth.LoadCredentials()

	report, err := runner.Run(
		"inovibe.api.v3.GroupService.NestedUsers",
		"group.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"groupid": "607f9db4-7eee-4a08-894d-356c8a462ae1"}`),
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
