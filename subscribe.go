package wemo

/*
  Need to set up a subscription service to Wemo events.

  1. Take a discovered device and
  2. Send a subscribe message to
    deviceIP:devicePort/upnp/event/basicevent1
  3. If the responce is 200, the subscription is successful and ...
  4. ... thus it should be added to the subscribed device list

Need to decide how to handle things next...
  5. When state is emitted record state changes

*/

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Subscriptions Structure
//var Subscriptions map[string]SubscriptionInfo

// Device and Subscription Info
type SubscriptionInfo struct {
  DeviceInfo
  State string
  Timeout int
  Sid string
}


// Structure for XML to Parse to
type Deviceevent struct {
	XMLName     xml.Name `xml:"propertyset"`
	BinaryState string   `xml:"property>BinaryState"`
}

// Listen for incomming subscribed state changes.
func Listener(listenerAddress string) {

	fmt.Println("Listening... ", listenerAddress)

	http.HandleFunc("/listener", func(w http.ResponseWriter, r *http.Request) {

		eventxml := Deviceevent{}

		if r.Method == "NOTIFY" {
			//This will be delt with when I work out how to handle events
			fmt.Println("SID: ", r.Header.Get("Sid"))
			fmt.Println("Host: ", r.Host)
			fmt.Println("Content-Type: ", r.Header.Get("Content-Type"))
			fmt.Println("Seq: ", r.Header.Get("Seq"))

			body, err := ioutil.ReadAll(r.Body)
			if err == nil {
				err := xml.Unmarshal([]byte(body), &eventxml)
				if err != nil {
					fmt.Println("Unmarshal error: ", err)
					return
				}

				//Need to work out how to handle events....
				fmt.Println("BinaryState: ", eventxml.BinaryState)
			}
		}
	})

	http.ListenAndServe(listenerAddress, nil)
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

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Client Request Error: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("Subscription Successful: ", resp.StatusCode)
		resp.Body.Close()
		return resp.Header.Get("Sid"), resp.StatusCode
	} else if resp.StatusCode == 400 {
		fmt.Println("Subscription Unsuccessful, Incompatible header fields: ", resp.StatusCode)
	} else if resp.StatusCode == 412 {
		fmt.Println("Subscription Unsuccessful, Precondition Failed: ", resp.StatusCode)
	} else {
		fmt.Println("Subscription Unsuccessful, Unable to accept renewal: ", resp.StatusCode)
	}

	resp.Body.Close()
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

	resp.Body.Close()
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

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Client Request Error: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("Unsubscription Successful: ", resp.StatusCode)
		resp.Body.Close()
		return resp.Header.Get("Sid"), resp.StatusCode
	} else if resp.StatusCode == 400 {
		fmt.Println("Unsubscription Unsuccessful, Incompatible header fields: ", resp.StatusCode)
	} else if resp.StatusCode == 412 {
		fmt.Println("Unsubscription Unsuccessful, Precondition Failed: ", resp.StatusCode)
	} else {
		fmt.Println("Unsubscription Unsuccessful, Unable to accept renewal: ", resp.StatusCode)
	}

	resp.Body.Close()
	return "", resp.StatusCode

}
