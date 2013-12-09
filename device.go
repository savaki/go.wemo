package wemo

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Device struct {
	Host string
}

type DeviceInfo struct {
	DeviceType      string `xml:"deviceType" json:"device-type"`
	FriendlyName    string `xml:"friendlyName" json:"friendly-name"`
	MacAddress      string `xml:"macAddress" json:"mac-address"`
	FirmwareVersion string `xml:"firmwareVersion" json:"firmware-version"`
	SerialNumber    string `xml:"serialNumber" json:"serial-number"`
	BinaryState     int    `xml:"binaryState" json:"binary-state"`
}

type BelkinResponse struct {
	Device DeviceInfo `xml:"device"`
}

func (self *Device) FetchDeviceInfo() (*DeviceInfo, error) {
	uri := fmt.Sprintf("http://%s/setup.xml", self.Host)
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	resp := new(BelkinResponse)
	err = xml.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return &resp.Device, nil
}

func (self *Device) Off() {
	self.changeState(false)
}

func (self *Device) On() {
	self.changeState(true)
}

func (self *Device) Toggle() {
	deviceInfo, err := self.FetchDeviceInfo()
	if err == nil {
		newState := deviceInfo.BinaryState == 0
		self.changeState(newState)
	}
}

func (self *Device) changeState(newState bool) error {
	fmt.Printf("changeState(%v)\n", newState)
	message := newSetBinaryStateMessage(newState)
	response, err := post(self.Host, message)
	if err != nil {
		log.Printf("unable to post message, %s\n", err.Error())
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("couldn't read body from message => %s\n", err.Error())
			return err
		}
		content := string(data)

		gripe := fmt.Sprintf("changeState(%v) => %s", newState, content)
		log.Println(gripe)
		return errors.New(gripe)
	}

	return nil
}
