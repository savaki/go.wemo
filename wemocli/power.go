package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/danward79/go.wemo"
)

var onCommand = cli.Command{
	Name: "on",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host", Value: "192.168.1.8:49153", Usage: "device host and ip e.g. 10.0.1.2:49128"},
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

var bulbCommand = cli.Command{
	Name:        "bulb",
	Usage:       "Command a bulb!",
	Description: "bulb --host 192.168.1.25:49153 dim 255",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host", Value: "192.168.1.25:49153", Usage: "device host and ip e.g. 10.0.1.2:49128"},
		cli.StringFlag{Name: "id", Value: "", Usage: "device id"},
	},
	Action: bulbAction,
}

func bulbAction(c *cli.Context) {
	host := c.String("host")
	id := c.String("id")
	args := c.Args()

	var cmd, value string
	if len(args) > 0 {
		cmd = args[0]

		if cmd == "dim" {
			value = args[1]
		}
	}

	device := &wemo.Device{
		Host: host,
	}

	err := device.Bulb(id, cmd, value, false)
	if err != nil {
		log.Println(err)
	}
}

var bulbStatusCommand = cli.Command{
	Name:        "bulbStatus",
	Usage:       "Status of a bulb!",
	Description: "bulb --host 192.168.1.25:49153 --id 94103EF6BF42867F",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host", Value: "192.168.1.25:49153", Usage: "device host and ip e.g. 10.0.1.2:49128"},
		cli.StringFlag{Name: "id", Value: "", Usage: "device id"},
	},
	Action: bulbStatusAction,
}

func bulbStatusAction(c *cli.Context) {
	host := c.String("host")
	id := c.String("id")

	device := &wemo.Device{
		Host: host,
	}

	result, err := device.GetBulbStatus(id)
	if err != nil {
		log.Println(err)
	}

	for k, v := range result {
		fmt.Println("DeviceID:", k, "State:", v)
	}
}
