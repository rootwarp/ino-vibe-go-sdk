package loadtest

import (
	"log"
	"os"
	"testing"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
)

func TestDeviceList(t *testing.T) {
	report, err := runner.Run(
		"inovibe.api.v3.DeviceService.List",
		"device.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"installStatus": 3}`),
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
	report, err := runner.Run(
		"inovibe.api.v3.DeviceService.Detail",
		"device.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"devid": "000000030000000000000001"}`),
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
	report, err := runner.Run(
		"inovibe.api.v3.DeviceService.UpdateInfo",
		"device.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"devid": "000000030000000000000001", "aliasValue": "dummy"}`),
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
	report, err := runner.Run(
		"inovibe.api.v3.DeviceService.UpdateStatus",
		"device.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"devid": "000000030000000000000001", "batteryValue": 100}`),
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
	report, err := runner.Run(
		"inovibe.api.v3.DeviceService.UpdateConfig",
		"device.ino-vibe.ino-on.dev:443",
		runner.WithDataFromJSON(`{"devid": "000000030000000000000001", "intThresholdValue": 1050}`),
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
