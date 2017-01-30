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
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	messageHeader = `<?xml version="1.0" encoding="utf-8"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body>`
	messageFooter = `</s:Body></s:Envelope>`
)

func post(hostAndPort, service, action, body string) (*http.Response, error) {
	tcpConn, err := timeoutDialer(2*time.Second, 2*time.Second)("tcp", hostAndPort)
	if err != nil {
		return nil, err
	}
	defer tcpConn.Close()

	preamble := fmt.Sprintf("POST http://%v/upnp/control/%s1 HTTP/1.1\r\nContent-type: text/xml; charset=\"utf-8\"\r\nSOAPACTION: \"urn:Belkin:service:%s:1#%s\"\r\nContent-Length: %v\r\n\r\n", hostAndPort, service, service, action, len(body))
	tcpConn.Write([]byte(preamble + body))

	data, err := ioutil.ReadAll(tcpConn)
	if err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(bytes.NewReader(data)), nil)
}

func newGetBinaryStateMessage() string {
	return messageHeader + `<u:GetBinaryState xmlns:u="urn:Belkin:service:basicevent:1"></u:GetBinaryState>` + messageFooter
}

func newSetBinaryStateMessage(on bool) string {
	value := 0
	if on {
		value = 1
	}

	return fmt.Sprintf(messageHeader+`<u:SetBinaryState xmlns:u="urn:Belkin:service:basicevent:1"><BinaryState>%v</BinaryState></u:SetBinaryState>`+messageFooter, value)
}

func newGetInsightParamsMessage() string {
	return messageHeader + `<u:GetInsightParams xmlns:u="urn:Belkin:service:insight:1"></u:GetInsightParams>` + messageFooter
}

func newGetBridgeEndDevices(u string) string {
	return fmt.Sprintf(messageHeader+`<u:GetEndDevices xmlns:u="urn:Belkin:service:bridge:1"><DevUDN>%s</DevUDN><ReqListType>PAIRED_LIST</ReqListType></u:GetEndDevices>`+messageFooter, u)
}

func newSetBulbStatus(id, capability, value string, group bool) string {
	g := "NO"
	if group {
		g = "YES"
	}

	return fmt.Sprintf(messageHeader+
		`<u:SetDeviceStatus xmlns:u="urn:Belkin:service:bridge:1">
			<DeviceStatusList>
		&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;&lt;DeviceStatus&gt;&lt;IsGroupAction&gt;%s&lt;/IsGroupAction&gt;&lt;DeviceID available=&quot;YES&quot;&gt;%s&lt;/DeviceID&gt;&lt;CapabilityID&gt;%s&lt;/CapabilityID&gt;&lt;CapabilityValue&gt;%s&lt;/CapabilityValue&gt;&lt;/DeviceStatus&gt;
		</DeviceStatusList>
	</u:SetDeviceStatus>`+
		messageFooter, g, id, capability, value)
}

func newGetBulbStatus(id string) string {
	return fmt.Sprintf(messageHeader+`<u:GetDeviceStatus xmlns:u="urn:Belkin:service:bridge:1">
		<DeviceIDs>%s</DeviceIDs>
		</u:GetDeviceStatus>`+messageFooter, id)
}
