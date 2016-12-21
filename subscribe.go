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
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

//SubscriptionInfo struct
type SubscriptionInfo struct {
	DeviceInfo
	Timeout     int
	Sid         string
	Host        string
	Deviceevent Deviceevent
}

//Deviceevent Structure for XML to Parse to
type Deviceevent struct {
	XMLName     xml.Name   `xml:"propertyset"`
	BinaryState string     `xml:"property>BinaryState"`
	StateEvent  StateEvent `xml:"property>StatusChange>StateEvent"`
}

//StateEvent ...
type StateEvent struct {
	DeviceID     string `xml:"DeviceID"`
	CapabilityID string `xml:"CapabilityId"`
	Value        string `xml:"Value"`
}

//SubscriptionEvent Structure for sending subscribed event data with
type SubscriptionEvent struct {
	Sid         string
	Deviceevent Deviceevent
}

//Listener Listen for incomming subscribed state changes.
func Listener(listenerAddress string, cs chan SubscriptionEvent) {

	log.Println("Listening... ", listenerAddress)

	http.HandleFunc("/listener", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "NOTIFY" {
			err := emitEvent(r, cs)
			if err != nil {
				log.Println("Event emit error: ", err)
			}
		}
	})

	err := http.ListenAndServe(listenerAddress, nil)
	if err != nil {
		log.Println("From Listen and Serve Err! ", err)
	}
}

//ManageSubscription Manage firstly the subscription and then the resubscription of this device.
func (d *Device) ManageSubscription(listenerAddress string, timeout int, subscriptions map[string]*SubscriptionInfo) (string, int) {
	/*  Subscribe to the device. Add device to subscriptions list

	    Once the device has a SID, it should have resubscriptions requested before the timeout.

	    Should a resubscription fail, an attempt should be made to unsubscribe and
	    then subscribe to the device in question. Returning the new SID or an error

	    The new SID should be updated in the subscription list and the old item removed.
	*/

	// Initial Subscribe
	info, _ := d.FetchDeviceInfo(context.Background())
	address := fmt.Sprintf("http://%s/upnp/event/basicevent1", d.Host)
	path := "/upnp/event/basicevent1"
	if info.DeviceType == "urn:Belkin:device:bridge:1" {
		address = fmt.Sprintf("http://%s/upnp/event/bridge1", d.Host)
		path = "/upnp/event/bridge1"
	}

	id, err := d.Subscribe(listenerAddress, address, path, timeout)
	if err != 200 {
		log.Println("Error with initial subscription: ", err)
		return "", err
	}
	fmt.Println("Returned ID", id)
	subscriptions[id] = &SubscriptionInfo{*info, timeout, id, d.Host, Deviceevent{}}

	// Setup resubscription timer
	offset := 30 //Renew early by offset seconds
	timer := time.NewTimer(time.Second * time.Duration(timeout-offset))
	go func() (string, int) {
		for _ = range timer.C {
			timer.Reset(time.Second * time.Duration(timeout-offset))

			// Resubscribe
			_, err = d.ReSubscribe(id, address, timeout)
			if err != 200 {

				// Failed to resubscribe so try unsubscribe, it is likely to fail but don't care.
				d.UnSubscribe(id, address)

				// Setup a new subscription, if this fails, next attempt will be when timer triggers again
				var newID string
				newID, err = d.Subscribe(listenerAddress, address, path, timeout)
				if err != 200 {
					log.Println("Error with subscription attempt: ", err)
				} else {
					// If the subscription is successful. Check if the new SID exists and if not remove it. Then add the new SID
					_, ok := subscriptions[newID]
					if ok == false {
						delete(subscriptions, id)
					}

					subscriptions[newID] = &SubscriptionInfo{*info, timeout, newID, d.Host, Deviceevent{}}
					id = newID
				}

			}
		}
		return "", err
	}()

	return id, err
}

//Subscribe to the device event emitter, return the Subscription ID (sid) and StatusCode
func (d *Device) Subscribe(listenerAddress, address, path string, timeout int) (string, int) {

	if timeout == 0 {
		timeout = 300
	}

	client := &http.Client{}

	req, err := http.NewRequest("SUBSCRIBE", address, nil)
	if err != nil {
		log.Println("http NewRequest Err: ", err)
	}

	req.Header.Add("host", fmt.Sprintf("http://%s", d.Host))
	req.Header.Add("path", path)
	req.Header.Add("callback", fmt.Sprintf("<http://%s/listener>", listenerAddress))
	req.Header.Add("nt", "upnp:event")
	req.Header.Add("timeout", fmt.Sprintf("Second-%d", timeout))

	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Client Request Error: ", err)
		return "", 0 //TODO:Check that this return is correct.
	}
	defer resp.Body.Close()

	log.Println(statusMessage("Subscription", d.Host, resp.StatusCode))

	if resp.StatusCode == http.StatusOK {
		return resp.Header.Get("Sid"), resp.StatusCode
	}

	return "", resp.StatusCode
}

//UnSubscribe According to the spec all subscribers must unsubscribe when the publisher is no longer required to provide state updates. Return the StatusCode
func (d *Device) UnSubscribe(sid, address string) int {

	client := &http.Client{}

	req, err := http.NewRequest("UNSUBSCRIBE", address, nil)
	if err != nil {
		log.Println("http NewRequest Err: ", err)
		return 0
	}

	req.Header.Add("host", fmt.Sprintf("http://%s", d.Host))
	req.Header.Add("SID", sid)

	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Client Request Error: ", err)
		return 0 //TODO:Check that this return is correct
	}
	defer resp.Body.Close()

	log.Println(statusMessage("Unsubscription", d.Host, resp.StatusCode))

	return resp.StatusCode
}

//ReSubscribe The subscription to the device must be renewed before the timeout. Return the Subscription ID (sid) and StatusCode
func (d *Device) ReSubscribe(sid, address string, timeout int) (string, int) {

	if timeout == 0 {
		timeout = 300
	}

	client := &http.Client{}

	req, err := http.NewRequest("SUBSCRIBE", address, nil)
	if err != nil {
		log.Println("http NewRequest Err: ", err)
		return "", 0
	}

	req.Header.Add("host", fmt.Sprintf("http://%s", d.Host))
	req.Header.Add("SID", sid)
	req.Header.Add("timeout", fmt.Sprintf("Second-%d", timeout))

	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Client Request Error: ", err)
		return "", 0 //TODO:Check that this return is correct
	}
	defer resp.Body.Close()

	log.Println(statusMessage("Resubscription", d.Host, resp.StatusCode))

	if resp.StatusCode == http.StatusOK {
		return resp.Header.Get("Sid"), resp.StatusCode
	}

	return "", resp.StatusCode
}

func statusMessage(action, host string, statusCode int) string {
	switch statusCode {
	case http.StatusOK:
		return fmt.Sprintf("%s Successful: %s, %d", action, host, statusCode)
	case http.StatusBadRequest:
		return fmt.Sprintf("%s Unsuccessful, Incompatible header fields: %s, %d", action, host, statusCode)
	case http.StatusPreconditionFailed:
		return fmt.Sprintf("%s Unsuccessful, Precondition Failed: %s, %d", action, host, statusCode)
	default:
		return fmt.Sprintf("%s Unsuccessful, Unable to accept renewal: %s, %d", action, host, statusCode)
	}
}

func emitEvent(r *http.Request, cs chan SubscriptionEvent) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err == nil {
		eventxml := Deviceevent{}
		err = xml.Unmarshal([]byte(html.UnescapeString(string(body))), &eventxml)
		if err != nil {
			return err
		}

		cs <- SubscriptionEvent{r.Header.Get("Sid"), eventxml}
	}

	return err
}
