package main

import (
	//"code.google.com/p/go.net/context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/codegangsta/cli"
	"github.com/danward79/go.wemo"
	"golang.org/x/net/context"
)

var discoverCommand = cli.Command{
	Name:        "discover",
	Usage:       "find devices in the local network",
	Description: "search for devices in the local network",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "interface", Value: "en0", Usage: "search by interface"}, //, ""},
		cli.StringFlag{Name: "ip", Value: "", Usage: "discovery wemo by ip"},          //,""},
		cli.IntFlag{Name: "timeout", Value: 3, Usage: "timeout"},                      //, ""},
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

	format := "%-20s %-20s %-21s %-20s\n"
	fmt.Printf(format,
		"Host",
		"Friendly Name",
		"Firmware Version",
		"Serial Number",
	)
	fmt.Printf(format,
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
			deviceInfo.FriendlyName,
			deviceInfo.FirmwareVersion,
			deviceInfo.SerialNumber)
	}
}
