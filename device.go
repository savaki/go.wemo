// Package wemo ...
/* Copyright 2014 Matt Ho
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
*/
package wemo

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html"
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
	EndDevices      EndDevices
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

	if deviceInfo.DeviceType == "urn:Belkin:device:bridge:1" {
		deviceInfo.EndDevices = *deviceInfo.Device.GetBridgeEndDevices(deviceInfo.UDN)
	}

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

// EndDevices ...
type EndDevices struct {
	DeviceListType string          `xml:"Body>GetEndDevicesResponse>DeviceLists>DeviceLists>DeviceList>DeviceListType"`
	EndDeviceInfo  []EndDeviceInfo `xml:"Body>GetEndDevicesResponse>DeviceLists>DeviceLists>DeviceList>DeviceInfos>DeviceInfo"`
}

// EndDeviceInfo ...
type EndDeviceInfo struct {
	DeviceIndex     string `xml:"DeviceIndex"`
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
func (d *Device) GetBridgeEndDevices(uuid string) *EndDevices {
	b := newGetBridgeEndDevices(uuid)

	response, err := post(d.Host, "bridge", "GetEndDevices", b)
	if err != nil {
		d.printf("unable to fetch bridge end devices => %s\n", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		d.printf("GetBridgeEndDevices returned status code => %d\n", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		d.printf("unable to read data => %s\n", err)
	}

	resp := EndDevices{}

	data = []byte(html.UnescapeString(string(data)))

	err = xml.Unmarshal(data, &resp)
	if err != nil {
		d.printf("Unmarshal Error: %s\n", err)
	}

	return &resp
}

//Bulb ...
func (d *Device) Bulb(id, cmd, value string, group bool) error {

	if id == "" {
		return errors.New("No ID provided")
	}

	capability := "10006"
	if cmd == "dim" {
		capability = "10008"

		s, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}

		if s > 255 || s < 0 {
			return errors.New("Dim value is out of bounds 0-255")
		}
	}

	if cmd == "on" {
		value = "1"
	} else if cmd == "off" {
		value = "0"
	}

	message := newSetBulbStatus(id, capability, value, group)

	response, err := post(d.Host, "bridge", "SetDeviceStatus", message)
	if err != nil {
		return errors.New("unable to SetDeviceStatus")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		_, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return errors.New("couldn't read body from message => " + err.Error())
		}

		return errors.New(string(response.StatusCode))
	}
	return nil
}

//BulbStatusList ...
type BulbStatusList struct {
	DeviceStatus []DeviceStatus `xml:"Body>GetDeviceStatusResponse>DeviceStatusList>DeviceStatusList>DeviceStatus"`
}

//DeviceStatus ...
type DeviceStatus struct {
	DeviceID        string `xml:"DeviceID"`
	CapabilityValue string `xml:"CapabilityValue"`
}

//GetBulbStatus return map of [DeviceID]status values
func (d *Device) GetBulbStatus(ids string) (map[string]string, error) {
	result := make(map[string]string)
	message := newGetBulbStatus(ids)

	response, err := post(d.Host, "bridge", "GetDeviceStatus", message)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch Bulb status => %s\n", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetBulbStatus returned status code => %d\n", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read data => %s\n", err)
	}

	data = []byte(html.UnescapeString(string(data)))

	statusInfo := BulbStatusList{}
	err = xml.Unmarshal(data, &statusInfo)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal Error: %s\n", err)
	}

	for k := range statusInfo.DeviceStatus {
		result[statusInfo.DeviceStatus[k].DeviceID] = statusInfo.DeviceStatus[k].CapabilityValue
	}

	return result, nil
}
