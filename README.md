wemo
====

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

```
package main

import (
	"fmt"
	"github.com/savaki/wemo"
)

func main() {
  device     := &Device{Host:"10.0.1.32:49153"}
  deviceInfo := device.FetchDeviceInfo()
  fmt.Printf("Found => %+v\n", deviceInfo)

  device.On()
  device.Off()
}
```