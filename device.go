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
	"golang.org/x/net/context"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/savaki/httpctx"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Device struct {
	Host   string
	Logger func(string, ...interface{}) (int, error)
}

type DeviceInfo struct {
	Device          *Device `json:"-"`
	DeviceType      string  `xml:"deviceType" json:"device-type"`
	FriendlyName    string  `xml:"friendlyName" json:"friendly-name"`
	MacAddress      string  `xml:"macAddress" json:"mac-address"`
	FirmwareVersion string  `xml:"firmwareVersion" json:"firmware-version"`
	SerialNumber    string  `xml:"serialNumber" json:"serial-number"`
}

type DeviceInfos []*DeviceInfo

func (d DeviceInfos) Len() int           { return len(d) }
func (d DeviceInfos) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DeviceInfos) Less(i, j int) bool { return d[i].FriendlyName < d[j].FriendlyName }

func (d *Device) printf(format string, args ...interface{}) {
	if d.Logger != nil {
		d.Logger(format, args...)
	}
}

func unmarshalDeviceInfo(data []byte) (*DeviceInfo, error) {
	resp := struct {
		DeviceInfo DeviceInfo `xml:"device"`
	}{}
	err := xml.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.DeviceInfo, nil
}

func (d *Device) FetchDeviceInfo(ctx context.Context) (*DeviceInfo, error) {
	var data []byte

	uri := fmt.Sprintf("http://%s/setup.xml", d.Host)
	err := httpctx.NewClient().Get(ctx, uri, nil, &data)
	if err != nil {
		return nil, err
	}

	deviceInfo, err := unmarshalDeviceInfo(data)
	if err != nil {
		return nil, err
	}

	deviceInfo.Device = d
	return deviceInfo, nil
}

func (d *Device) GetBinaryState() int {
	message := newGetBinaryStateMessage()
	response, err := post(d.Host, "GetBinaryState", message)
	if err != nil {
		d.printf("unable to fetch BinaryState => %s\n", err)
		return -1
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		d.printf("GetBinaryState returned status code => %d\n", response.StatusCode)
		return -1
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		d.printf("unable to read data => %s\n", err)
		return -1
	}

	re := regexp.MustCompile(`.*<BinaryState>(\d+)</BinaryState>.*`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) != 2 {
		d.printf("unable to find BinaryState response in message => %s\n", string(data))
		return -1
	}

	result, _ := strconv.Atoi(matches[1])
	return result
}

func (d *Device) Off() {
	d.changeState(false)
}

func (d *Device) On() {
	d.changeState(true)
}

func (d *Device) Toggle() {
	if binaryState := d.GetBinaryState(); binaryState == 0 {
		d.On()
	} else {
		d.Off()
	}
}

func (d *Device) changeState(newState bool) error {
	fmt.Printf("changeState(%v)\n", newState)
	message := newSetBinaryStateMessage(newState)
	response, err := post(d.Host, "SetBinaryState", message)
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
