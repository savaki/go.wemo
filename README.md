wemo
====

Simple package to interface with Belkin wemo devices.

### Example - Device discovery

```
package main

import (
	"fmt"
	"github.com/savaki/wemo"
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
  "github.com/savaki/wemo"
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
  "github.com/savaki/wemo"
  "time"
)

func main() {
  api, _ := wemo.NewByInterface("en0")
  api.On("Left Light", 3*time.Second)
  api.Off("Left Light", 3*time.Second)
  api.Toggle("Left Light", 3*time.Second)
}
```