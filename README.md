go.wemo
=======

[![GoDoc](http://godoc.org/github.com/savaki/go.wemo?status.png)](http://godoc.org/github.com/savaki/go.wemo)

Simple package to interface with Belkin wemo devices.

### Example - Device discovery

```
package main

import (
	"fmt"
	"github.com/savaki/go.wemo"
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
  "github.com/savaki/go.wemo"
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
  device.GetBinaryState() // returns 0 or 1
}
```

### Example - Control a named light

As a convenience method, you can control lights through a more generic interface.

```
package main

import (
  "github.com/savaki/go.wemo"
  "time"
)

func main() {
  api, _ := wemo.NewByInterface("en0")
  api.On("Left Light", 3*time.Second)
  api.Off("Left Light", 3*time.Second)
  api.Toggle("Left Light", 3*time.Second)
}
```

###Example - Managing Subscriptions

This is a work in progress (I am learning Go ;-\)


```
package main

import (
	"fmt"
	"github.com/savaki/go.wemo"
	"time"
)

func main() {
  
  listenerAddress := "192.168.0.6:6767"
  timeout := 300
  
  api, _ := wemo.NewByInterface("en0")
  
  devices, _ := api.DiscoverAll(3*time.Second)
  
  subscriptions := make(map[string]wemo.SubscriptionInfo)
  
  for _, device := range devices {
    
    info, _ := device.FetchDeviceInfo()
    id, err := device.Subscribe(listenerAddress, timeout)
    if err == 200 {
      subscriptions[device.Host] = wemo.SubscriptionInfo{*info, "0", timeout, id }      
    }
  }
  
  wemo.Listener(listenerAddress)
   
  fmt.Println("Subscription List: ", subscriptions)
  
}

```



