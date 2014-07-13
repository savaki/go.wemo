package wemo

/*
  Need to set up a subscription service to Wemo events.

  1. Take a discovered device and
  2. Send a subscribe message to
    deviceIP:devicePort/upnp/event/basicevent1
  3. If the responce is 200, the subscription is successful and ...
  4. ... thus it should be added to the subscribed device list
  5. Subscriptions should be renewed around the timeout period
  6. When state is emitted record state changes against the subscription id (SID)

*/

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
  "time"
)


// Device and Subscription Info
type SubscriptionInfo struct {
  DeviceInfo
  State string
  Timeout int
  Sid string
  Host string
}

// Structure for XML to Parse to
type Deviceevent struct {
	XMLName     xml.Name `xml:"propertyset"`
	BinaryState string   `xml:"property>BinaryState"`
}

// Structure for XML to Parse to
type SubscriptionEvent struct {
	Sid string
	State string
}

// Listen for incomming subscribed state changes.
func Listener(listenerAddress string, cs chan SubscriptionEvent) {

	fmt.Println("Listening... ", listenerAddress)

	http.HandleFunc("/listener", func(w http.ResponseWriter, r *http.Request) {

		eventxml := Deviceevent{}

		if r.Method == "NOTIFY" {

			body, err := ioutil.ReadAll(r.Body)
			if err == nil {
				err := xml.Unmarshal([]byte(body), &eventxml)
				if err != nil {
					fmt.Println("Unmarshal error: ", err)
					return
				}
        
        cs <- SubscriptionEvent {r.Header.Get("Sid"), eventxml.BinaryState}
			}
		}
	})

	err := http.ListenAndServe(listenerAddress, nil)
  if err != nil {
    fmt.Println("From Listen and Serve an Err! ", err)
  }
}

// Manage firstly the subscription and then the resubscription of this device.
func (self *Device) ManageSubscription(listenerAddress string, timeout int) (string, int){
  /*  Subscribe to the device.
      Once the device has a SID, it should have resubscriptions requested before the timeout.
      Should a resubscription fail, an attempt should be made to unsubscribe and 
      then subscribe to the device in question. Returning the new SID or an error
  */
  
  id, err := self.Subscribe(listenerAddress, timeout)
  if err != 200 {
    fmt.Println("Error with initial subscription: ", err)
    return "", err
  }
  
  timer := time.NewTimer(time.Second * time.Duration(timeout))
  go func() (string, int){
    for _ = range timer.C {
      timer.Reset(time.Second * time.Duration(timeout))
  
      // Should really do something with the SID returned to check it has not changed
      _, err := self.ReSubscribe(id, timeout)
      if err != 200 {
        fmt.Println("Error with resubscription: ", err)
        return "", err
      }
    }
    return "", err
  }()
  
  return id, err
  
}

// Subscribe to the device event emitter, return the Subscription ID (sid) and StatusCode
func (self *Device) Subscribe(listenerAddress string, timeout int) (string, int) {

	host := self.Host

	address := fmt.Sprintf("http://%s/upnp/event/basicevent1", host)

	if timeout == 0 {
		timeout = 300
	}

	client := &http.Client{}

	req, err := http.NewRequest("SUBSCRIBE", address, nil)
	if err != nil {
		log.Fatal("http NewRequest Err: ", err)
	}

	req.Header.Add("host", fmt.Sprintf("http://%s", host))
	req.Header.Add("path", "/upnp/event/basicevent1")
	req.Header.Add("callback", fmt.Sprintf("<http://%s/listener>", listenerAddress))
	req.Header.Add("nt", "upnp:event")
	req.Header.Add("timeout", fmt.Sprintf("Second-%d", timeout))
  
  req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Client Request Error: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("Subscription Successful: ", resp.StatusCode)
		//resp.Body.Close()
		return resp.Header.Get("Sid"), resp.StatusCode
	} else if resp.StatusCode == 400 {
		fmt.Println("Subscription Unsuccessful, Incompatible header fields: ", resp.StatusCode)
	} else if resp.StatusCode == 412 {
		fmt.Println("Subscription Unsuccessful, Precondition Failed: ", resp.StatusCode)
	} else {
		fmt.Println("Subscription Unsuccessful, Unable to accept renewal: ", resp.StatusCode)
	}

	//resp.Body.Close()
	return "", resp.StatusCode

}

// According to the spec all subscribers must unsubscribe when the publisher is no longer required to provide state updates. Return the StatusCode
func (self *Device) UnSubscribe(sid string) int {

	host := self.Host

	address := fmt.Sprintf("http://%s/upnp/event/basicevent1", host)

	client := &http.Client{}

	req, err := http.NewRequest("UNSUBSCRIBE", address, nil)
	if err != nil {
		log.Fatal("http NewRequest Err: ", err)
	}

	req.Header.Add("host", fmt.Sprintf("http://%s", host))
	req.Header.Add("SID", sid)
  
  req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Client Request Error: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("Unsubscription Successful: ", resp.StatusCode)
	} else if resp.StatusCode == 400 {
		fmt.Println("Unsubscription Unsuccessful, Incompatible header fields: ", resp.StatusCode)
	} else if resp.StatusCode == 412 {
		fmt.Println("Unsubscription Unsuccessful, Precondition Failed: ", resp.StatusCode)
	} else {
		fmt.Println("Unsubscription Unsuccessful, Unable to accept renewal: ", resp.StatusCode)
	}

	//resp.Body.Close()
	return resp.StatusCode

}

// The subscription to the device must be renewed before the timeout. Return the Subscription ID (sid) and StatusCode
func (self *Device) ReSubscribe(sid string, timeout int) (string, int) {
  
	host := self.Host

	address := fmt.Sprintf("http://%s/upnp/event/basicevent1", host)

	if timeout == 0 {
		timeout = 300
	}

	client := &http.Client{}

	req, err := http.NewRequest("SUBSCRIBE", address, nil)
	if err != nil {
		log.Fatal("http NewRequest Err: ", err)
	}

	req.Header.Add("host", fmt.Sprintf("http://%s", host))
	req.Header.Add("SID", sid)
	req.Header.Add("timeout", fmt.Sprintf("Second-%d", timeout))

  req.Close = true
  
	resp, err := client.Do(req)
	if err != nil {
    fmt.Println("Response: ", self.Host, resp)
		log.Fatal("Client Request Error: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("Resubscription Successful: ", resp.StatusCode)
		//resp.Body.Close()
		return resp.Header.Get("Sid"), resp.StatusCode
	} else if resp.StatusCode == 400 {
		fmt.Println("Resubscription Unsuccessful, Incompatible header fields: ", resp.StatusCode)
	} else if resp.StatusCode == 412 {
		fmt.Println("Resubscription Unsuccessful, Precondition Failed: ", resp.StatusCode)
	} else {
		fmt.Println("Resubscription Unsuccessful, Unable to accept renewal: ", resp.StatusCode)
	}

	//resp.Body.Close()
	return "", resp.StatusCode

}
