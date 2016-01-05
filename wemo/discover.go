package main

import (
	"code.google.com/p/go.net/context"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/savaki/go.wemo"
	"log"
	"sort"
	"time"
)

var discoverCommand = cli.Command{
	Name:        "discover",
	Usage:       "find devices in the local network",
	Description: "search for devices in the local network",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "interface",
			Aliases: []string{"i"},
			Value: "",
			Usage: "search by interface",
			EnvVar: "WEMO_IFACE",
		},
		cli.IntFlag{
			Name: "timeout",
			Aliases: []string{"t"},
			Value: 3,
			Usage: "timeout period in seconds",
			EnvVar: "WEMO_TIMEOUT_DISCOVERY",
		},
	},
	Action: commandAction,
}

func commandAction(c *cli.Context) {
	timeout := c.Int("timeout")
	iface := c.String("interface")

	api, err := wemo.NewByInterface(iface)
	if err != nil {
		log.Fatal(err)
	}

	devices, err := api.DiscoverAll(time.Duration(timeout) * time.Second)
	if err != nil {
		log.Fatal(err)
	}

	format := "%-20s %-20s %-20s %-21s %-20s\n"
	fmt.Printf(format,
		"Host",
		"Mac Address",
		"Friendly Name",
		"Firmware Version",
		"Serial Number",
	)
	fmt.Printf(format,
		"----------------",
		"----------------",
		"----------------",
		"----------------",
		"----------------",
	)

	deviceInfos := wemo.DeviceInfos{}
	for _, device := range devices {
		deviceInfo, err := device.FetchDeviceInfo(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		deviceInfos = append(deviceInfos, deviceInfo)
	}

	sort.Sort(deviceInfos)
	for _, deviceInfo := range deviceInfos {
		fmt.Printf(format,
			deviceInfo.Device.Host,
			deviceInfo.MacAddress,
			deviceInfo.FriendlyName,
			deviceInfo.FirmwareVersion,
			deviceInfo.SerialNumber)
	}
}
