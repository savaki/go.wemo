package main

import (
	"github.com/codegangsta/cli"
	"github.com/savaki/go.wemo"
)

var onCommand = cli.Command{
	Name: "on",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "host",
			Aliases: []string{"h"},
			Value: "",
			Usage: "device host and ip e.g. 10.0.1.2:49128",
			EnvVar: "WEMO_POWER_HOST",
		},
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
		cli.StringFlag{
			Name: "host",
			Aliases: []string{"h"},
			Value: "",
			Usage: "device host and ip e.g. 10.0.1.2:49128",
			EnvVar: "WEMO_POWER_HOST",
		},
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
		cli.StringFlag{
			Name: "host",
			Aliases: []string{"h"},
			Value: "",
			Usage: "device host and ip e.g. 10.0.1.2:49128",
			EnvVar: "WEMO_POWER_HOST",
		},
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
