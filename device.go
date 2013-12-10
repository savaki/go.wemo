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

	fmt.Printf("%+v\n", resp.Device)

	return &resp.Device, nil
}

func (self *Device) FetchBinaryState() int {
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
	if binaryState := self.FetchBinaryState(); binaryState == 0 {
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
		log.Printf("unable to post message, %+v\n", err)
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
