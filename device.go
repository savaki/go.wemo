// Package wemo ...
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
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/savaki/httpctx"
)

// Device struct
type Device struct {
	Host   string
	Logger func(string, ...interface{}) (int, error)
}

// DeviceInfo struct
type DeviceInfo struct {
	Device          *Device `json:"-"`
	DeviceType      string  `xml:"deviceType" json:"device-type"`
	FriendlyName    string  `xml:"friendlyName" json:"friendly-name"`
	MacAddress      string  `xml:"macAddress" json:"mac-address"`
	FirmwareVersion string  `xml:"firmwareVersion" json:"firmware-version"`
	SerialNumber    string  `xml:"serialNumber" json:"serial-number"`
	UDN             string  `xml:"UDN" json:"UDN"`
}

// DeviceInfos slice
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

// FetchDeviceInfo from device
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

// GetBinaryState ...
func (d *Device) GetBinaryState() int {
	message := newGetBinaryStateMessage()
	response, err := post(d.Host, "basicevent", "GetBinaryState", message)
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

// Off toggle state to Off
func (d *Device) Off() {
	d.changeState(false)
}

// On toggle state to On
func (d *Device) On() {
	d.changeState(true)
}

// Toggle state
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
	response, err := post(d.Host, "basicevent", "SetBinaryState", message)
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

// InsightParams ...
type InsightParams struct {
	Power int // mW
}

// GetInsightParams ...
func (d *Device) GetInsightParams() *InsightParams {
	message := newGetInsightParamsMessage()
	response, err := post(d.Host, "insight", "GetInsightParams", message)
	if err != nil {
		d.printf("unable to fetch Power => %s\n", err)
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		d.printf("GetInsightParams returned status code => %d\n", response.StatusCode)
		return nil
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		d.printf("unable to read data => %s\n", err)
		return nil
	}

	// <s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body>
	// <u:GetInsightParamsResponse xmlns:u="urn:Belkin:service:metainfo:1">
	// <InsightParams>8|1471416661|8|3244|3182|15377|19|7300|1011115|1011115.000000|8000</InsightParams>
	// </u:GetInsightParamsResponse>

	re := regexp.MustCompile(`.*<InsightParams>(.+)</InsightParams>.*`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) != 2 {
		d.printf("unable to find GetInsightParams response in message => %s\n", string(data))
		return nil
	}
	split := strings.Split(matches[1], "|")
	if len(split) < 7 {
		d.printf("unable to parse InsightParams response => %s\n", string(data))
		return nil
	}
	power, err := strconv.Atoi(split[7])
	if err != nil {
		d.printf("failed to parse power: %v", err)
	}
	return &InsightParams{
		Power: power,
	}
}

// EndDevice ...
type EndDevice struct {
	DeviceID        string `xml:"DeviceID"`
	FriendlyName    string `xml:"FriendlyName"`
	FirmwareVersion string `xml:"FirmwareVersion"`
	CapabilityIDs   string `xml:"CapabilityIDs"`
	CurrentState    string `xml:"CurrentState"`
	Manufacturer    string `xml:"Manufacturer"`
	ModelCode       string `xml:"ModelCode"`
	ProductName     string `xml:"productName"`
	WeMoCertified   string `xml:"WeMoCertified"`
}

// GetBridgeEndDevices ...
func (d *Device) GetBridgeEndDevices(uuid string) {
	a := "GetEndDevices"
	b := newGetBridgeEndDevices(uuid)

	response, err := post(d.Host, "bridge", a, b)
	if err != nil {
		d.printf("unable to fetch bridge end devices => %s\n", err)
		//return nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		d.printf("GetBridgeEndDevices returned status code => %d\n", response.StatusCode)
		//return nil
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		d.printf("unable to read data => %s\n", err)
		//return nil
	}

	fmt.Println("Response Body:", string(data))
}
