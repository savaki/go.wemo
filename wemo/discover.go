package main

import (
	"code.google.com/p/go.net/context"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/andrewpurkett/go.wemo"
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
			Value: "",
			Usage: "search by interface",
			EnvVar: "WEMO_IFACE",
		},
		cli.IntFlag{
			Name: "timeout",
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

	format := "%-22s %-18s %-15s %-26s %-35s\n"
	fmt.Printf(format,
		"Host",
		"Mac Address",
		"Serial Number",
		"Friendly Name",
		"Firmware Version",
	)
	fmt.Printf(format,
		"---------------------",
		"-----------------",
		"--------------",
		"-------------------------",
		"-----------------------------------",
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
			fmt.Sprintf("%s:%s:%s:%s:%s:%s", deviceInfo.MacAddress[0:2], deviceInfo.MacAddress[2:4], deviceInfo.MacAddress[4:6], deviceInfo.MacAddress[6:8], deviceInfo.MacAddress[8:10], deviceInfo.MacAddress[10:12]),
			deviceInfo.SerialNumber,
			deviceInfo.FriendlyName,
			deviceInfo.FirmwareVersion)
	}
}
