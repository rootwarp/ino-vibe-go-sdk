package loadtest

import (
	"log"
	"os"
	"testing"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	iv_auth "github.com/rootwarp/ino-vibe-go-sdk/auth"
)

func TestDeviceList(t *testing.T) {
	token, err := iv_auth.LoadCredentials()

	report, err := runner.Run(
		"inovibe.api.v3.DeviceService.List",
		"device.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"installStatus": 3}`),
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

func TestDeviceDetail(t *testing.T) {
	token, err := iv_auth.LoadCredentials()

	report, err := runner.Run(
		"inovibe.api.v3.DeviceService.Detail",
		"device.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"devid": "000000030000000000000001"}`),
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

func TestDeviceUpdateInfo(t *testing.T) {
	token, err := iv_auth.LoadCredentials()

	report, err := runner.Run(
		"inovibe.api.v3.DeviceService.UpdateInfo",
		"device.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"devid": "000000030000000000000001", "aliasValue": "dummy"}`),
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

func TestDeviceUpdateStatus(t *testing.T) {
	token, err := iv_auth.LoadCredentials()

	report, err := runner.Run(
		"inovibe.api.v3.DeviceService.UpdateStatus",
		"device.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"devid": "000000030000000000000001", "batteryValue": 100}`),
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

func TestDeviceUpdateConfig(t *testing.T) {
	token, err := iv_auth.LoadCredentials()

	report, err := runner.Run(
		"inovibe.api.v3.DeviceService.UpdateConfig",
		"device.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"devid": "000000030000000000000001", "intThresholdValue": 1050}`),
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
