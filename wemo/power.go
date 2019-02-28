package main

import (
	"github.com/codegangsta/cli"
	"github.com/savaki/go.wemo"
)

var onCommand = cli.Command{
	Name: "on",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host", Value: "", Usage: "device host and ip e.g. 10.0.1.2:49128"},
	},
	Action: onAction,
}

func onAction(c *cli.Context) {
	host := c.String("host")
	device := &wemo.Device{
		Host: host,
	}
	device.On()
}

var offCommand = cli.Command{
	Name: "off",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host", Value: "", Usage: "device host and ip e.g. 10.0.1.2:49128"},
	},
	Action: offAction,
}

func offAction(c *cli.Context) {
	host := c.String("host")
	device := &wemo.Device{
		Host: host,
	}
	device.Off()
}

var toggleCommand = cli.Command{
	Name: "toggle",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host", Value: "", Usage: "device host and ip e.g. 10.0.1.2:49128"},
	},
	Action: toggleAction,
}

func toggleAction(c *cli.Context) {
	host := c.String("host")
	device := &wemo.Device{
		Host: host,
	}
	device.Toggle()
}
