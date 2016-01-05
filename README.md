go.wemo
=======

[![GoDoc](http://godoc.org/github.com/andrewpurkett/go.wemo?status.png)](http://godoc.org/github.com/andrewpurkett/go.wemo)

Simple package to interface with Belkin wemo devices.

## Standalone Usage Guide

Install the go library and configure it as needed, ensuring `~/bin/go/src/` is in your `GOPATH` environment variable. 

Run the following command from `~/bin/go/src/`:

`go get https://github.com/andrewpurkett/go.wemo`

navigate into the new directory:

`cd ~/bin/go/src/github.com/andrewpurkett/go.wemo`

Run `go get` to retrieve dependencies

navigate into the example usage directory:

`cd ~/bin/go/src/github.com/andrewpurkett/go.wemo/wemo`

Run `go get` again to retrieve dependencies for the example usage directory

Build the example usage tool:

`go build`

Then refer to the command line tool to see sample usage:
 
`~/bin/go/src/github.com/andrewpurkett/go.wemo/wemo/wemo`

`~/bin/go/src/github.com/andrewpurkett/go.wemo/wemo/wemo discover -h`

![Example usage](https://i.imgur.com/UYI2E4F.png)

If you were unable to build the CLI tool, run `go test` (in both `~/bin/go/src/github.com/andrewpurkett/go.wemo/wemo` and `~/bin/go/src/github.com/andrewpurkett/go.wemo/`), check your `GOPATH`, `GOROOT`, and repeat any other golang setup steps required.

## Utilizing the library in projects

Here is some example usage of the various functionality incorporated in this go repository:

### Example - Device discovery

```
package main

import (
	"fmt"
	"github.com/andrewpurkett/go.wemo"
	"time"
)

func main() {
  api, _ := wemo.NewByInterface("en0")
  devices, _ := api.DiscoverAll(3*time.Second)
  for _, device := range devices {
    fmt.Printf("Found %+v\n", device)
  }
}
```

### Example - Control a device

```
package main

import (
  "fmt"
  "github.com/andrewpurkett/go.wemo"
)

func main() {
  // you can either create a device directly OR use the
  // #Discover/#DiscoverAll methods to find devices
  device        := &wemo.Device{Host:"10.0.1.32:49153"}

  // retrieve device info
  deviceInfo, _ := device.FetchDeviceInfo()
  fmt.Printf("Found => %+v\n", deviceInfo)

  // device controls
  device.On()
  device.Off()
  device.Toggle()
  device.BinaryState() // returns 0 or 1
}
```

### Example - Control a named light

As a convenience method, you can control lights through a more generic interface.

```
package main

import (
  "github.com/andrewpurkett/go.wemo"
  "time"
)

func main() {
  api, _ := wemo.NewByInterface("en0")
  api.On("Left Light", 3*time.Second)
  api.Off("Left Light", 3*time.Second)
  api.Toggle("Left Light", 3*time.Second)
}
```
