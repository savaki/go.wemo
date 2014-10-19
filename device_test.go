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
	"encoding/json"
	"testing"
)

func assert(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("expected %s; got %s", expected, actual)
		t.Fail()
	}
}

func TestParseResponseXML(t *testing.T) {
	data := []byte(`<?xml version="1.0"?>
<root xmlns="urn:Belkin:device-1-0">
  <specVersion>
    <major>1</major>
    <minor>0</minor>
  </specVersion>
  <device>
    <deviceType>urn:Belkin:device:controllee:1</deviceType>
    <friendlyName>Pirate Light Right</friendlyName>
    <manufacturer>Belkin International Inc.</manufacturer>
    <manufacturerURL>http://www.belkin.com</manufacturerURL>
    <modelDescription>Belkin Plugin Socket 1.0</modelDescription>
    <modelName>Socket</modelName>
    <modelNumber>1.0</modelNumber>
    <modelURL>http://www.belkin.com/plugin/</modelURL>
    <serialNumber>221248K0102C92</serialNumber>
    <UDN>uuid:Socket-1_0-221248K0102C92</UDN>
    <UPC>123456789</UPC>
    <macAddress>EC1A5974B1EC</macAddress>
    <firmwareVersion>WeMo_US_2.00.2769.PVT</firmwareVersion>
    <iconVersion>0|49153</iconVersion>
    <binaryState>1</binaryState>
    <iconList>
      <icon>
        <mimetype>jpg</mimetype>
        <width>100</width>
        <height>100</height>
        <depth>100</depth>
        <url>icon.jpg</url>
      </icon>
    </iconList>
    <serviceList>
      <service>
        <serviceType>urn:Belkin:service:WiFiSetup:1</serviceType>
        <serviceId>urn:Belkin:serviceId:WiFiSetup1</serviceId>
        <controlURL>/upnp/control/WiFiSetup1</controlURL>
        <eventSubURL>/upnp/event/WiFiSetup1</eventSubURL>
        <SCPDURL>/setupservice.xml</SCPDURL>
      </service>
      <service>
        <serviceType>urn:Belkin:service:timesync:1</serviceType>
        <serviceId>urn:Belkin:serviceId:timesync1</serviceId>
        <controlURL>/upnp/control/timesync1</controlURL>
        <eventSubURL>/upnp/event/timesync1</eventSubURL>
        <SCPDURL>/timesyncservice.xml</SCPDURL>
      </service>
      <service>
        <serviceType>urn:Belkin:service:basicevent:1</serviceType>
        <serviceId>urn:Belkin:serviceId:basicevent1</serviceId>
        <controlURL>/upnp/control/basicevent1</controlURL>
        <eventSubURL>/upnp/event/basicevent1</eventSubURL>
        <SCPDURL>/eventservice.xml</SCPDURL>
      </service>
      <service>
        <serviceType>urn:Belkin:service:firmwareupdate:1</serviceType>
        <serviceId>urn:Belkin:serviceId:firmwareupdate1</serviceId>
        <controlURL>/upnp/control/firmwareupdate1</controlURL>
        <eventSubURL>/upnp/event/firmwareupdate1</eventSubURL>
        <SCPDURL>/firmwareupdate.xml</SCPDURL>
      </service>
      <service>
        <serviceType>urn:Belkin:service:rules:1</serviceType>
        <serviceId>urn:Belkin:serviceId:rules1</serviceId>
        <controlURL>/upnp/control/rules1</controlURL>
        <eventSubURL>/upnp/event/rules1</eventSubURL>
        <SCPDURL>/rulesservice.xml</SCPDURL>
      </service>

      <service>
        <serviceType>urn:Belkin:service:metainfo:1</serviceType>
        <serviceId>urn:Belkin:serviceId:metainfo1</serviceId>
        <controlURL>/upnp/control/metainfo1</controlURL>
        <eventSubURL>/upnp/event/metainfo1</eventSubURL>
        <SCPDURL>/metainfoservice.xml</SCPDURL>
      </service>

      <service>
        <serviceType>urn:Belkin:service:remoteaccess:1</serviceType>
        <serviceId>urn:Belkin:serviceId:remoteaccess1</serviceId>
        <controlURL>/upnp/control/remoteaccess1</controlURL>
        <eventSubURL>/upnp/event/remoteaccess1</eventSubURL>
        <SCPDURL>/remoteaccess.xml</SCPDURL>
      </service>

      <service>
        <serviceType>urn:Belkin:service:deviceinfo:1</serviceType>
        <serviceId>urn:Belkin:serviceId:deviceinfo1</serviceId>
        <controlURL>/upnp/control/deviceinfo1</controlURL>
        <eventSubURL>/upnp/event/deviceinfo1</eventSubURL>
        <SCPDURL>/deviceinfoservice.xml</SCPDURL>
      </service>

    </serviceList>
    <presentationURL>/pluginpres.html</presentationURL>
  </device>
</root>
	`)

	// When
	deviceInfo, _ := unmarshalDeviceInfo(data)
	assert(t, deviceInfo.DeviceType, "urn:Belkin:device:controllee:1")
	assert(t, deviceInfo.FirmwareVersion, "WeMo_US_2.00.2769.PVT")
	assert(t, deviceInfo.FriendlyName, "Pirate Light Right")
	assert(t, deviceInfo.MacAddress, "EC1A5974B1EC")
	assert(t, deviceInfo.SerialNumber, "221248K0102C92")
}

func TestParseResponseJSON(t *testing.T) {
	data := []byte(`{
			"device-type":"das device",
			"friendly-name":"das name",
			"mac-address":"das address",
			"firmware-version":"das firmware",
			"serial-number":"das serial"
		}`)

	deviceInfo := new(DeviceInfo)
	json.Unmarshal(data, deviceInfo)

	// Then
	assert(t, deviceInfo.DeviceType, "das device")
	assert(t, deviceInfo.FirmwareVersion, "das firmware")
	assert(t, deviceInfo.FriendlyName, "das name")
	assert(t, deviceInfo.MacAddress, "das address")
	assert(t, deviceInfo.SerialNumber, "das serial")
}
