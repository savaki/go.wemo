package main

import (
	"github.com/codegangsta/cli"
	"github.com/savaki/go.wemo"
	"log"
	"time"
)

var discoverCommand = cli.Command{
	Name:        "discover",
	Usage:       "find devices in the local network",
	Description: "search for devices in the local network",
	Flags: []cli.Flag{
		cli.StringFlag{"interface", "", "search by interface", ""},
		cli.StringFlag{"ip", "", "discovery wemo by ip", ""},
		cli.IntFlag{"timeout", 3, "timeout", ""},
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

	for _, device := range devices {
		log.Printf("%#v\n", device)
	}
}
