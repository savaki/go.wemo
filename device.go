// Copyright 2014 Matt Ho
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package wemo

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
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
}

type BelkinResponse struct {
	Device DeviceInfo `xml:"device"`
}

func (self *Device) FetchDeviceInfo() (*DeviceInfo, error) {
	uri := fmt.Sprintf("http://%s/setup.xml", self.Host)
	response, err := client.Get(uri)
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

	fmt.Printf("%+v\n", resp.Device)

	return &resp.Device, nil
}

func (self *Device) GetBinaryState() int {
	message := newGetBinaryStateMessage()
	response, err := post(self.Host, "GetBinaryState", message)
	if err != nil {
		log.Printf("unable to fetch BinaryState => %s\n", err)
		return -1
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("GetBinaryState returned status code => %d\n", response.StatusCode)
		return -1
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("unable to read data => %s\n", err)
		return -1
	}

	re := regexp.MustCompile(`.*<BinaryState>(\d+)</BinaryState>.*`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) != 2 {
		log.Printf("unable to find BinaryState response in message => %s\n", string(data))
		return -1
	}

	result, _ := strconv.Atoi(matches[1])
	return result
}

func (self *Device) Off() {
	self.changeState(false)
}

func (self *Device) On() {
	self.changeState(true)
}

func (self *Device) Toggle() {
	if binaryState := self.GetBinaryState(); binaryState == 0 {
		self.On()
	} else {
		self.Off()
	}
}

func (self *Device) changeState(newState bool) error {
	fmt.Printf("changeState(%v)\n", newState)
	message := newSetBinaryStateMessage(newState)
	response, err := post(self.Host, "SetBinaryState", message)
	if err != nil {
		log.Println("unable to SetBinaryState")
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("couldn't read body from message => " + err.Error())
			return err
		}
		content := string(data)

		gripe := fmt.Sprintf("changeState(%v) => %s", newState, content)
		log.Println(gripe)
		return errors.New(gripe)
	}

	return nil
}
