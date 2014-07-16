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

This is an example of discovering devices, subscribing to there events and being notified of changed to there state. Resubscriptions and managed automatically at the timeout specified. Subscriber details and state are maintained in a map.

```
package main

import (
	"github.com/savaki/go.wemo"
	"time"
  "log"
)

func main() {
  
  listenerAddress := "192.168.0.6:6767"
  timeout := 300
  
  api, _ := wemo.NewByInterface("en0")
  
  devices, _ := api.DiscoverAll(3*time.Second)
 
  subscriptions := make(map[string]*wemo.SubscriptionInfo)

  for _, device := range devices {
    _, err := device.ManageSubscription(listenerAddress, timeout, subscriptions)
    if err != 200 {
      log.Println("Initial Error Subscribing: ", err)   
    }
  }
  
  cs := make(chan wemo.SubscriptionEvent)

  go wemo.Listener(listenerAddress, cs)

  for m := range cs{
    if _, ok := subscriptions[m.Sid]; ok {
      subscriptions[m.Sid].State = m.State
      log.Println("---Subscriber Event: ", subscriptions[m.Sid])
    } else {
      log.Println("Does'nt exist, ", m.Sid)
    }
  }

}
```



