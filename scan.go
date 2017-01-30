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
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

//Constants associated with Scanning
const (
	SSDPBROADCAST = "239.255.255.250:1900"
	MSEARCH       = "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 10\r\nST: %s\r\nUSER-AGENT: unix/5.1 UPnP/1.1 crash/1.0\r\n\r\n"
	LOCATION      = "LOCATION: "
)

// scan the multicast
func (w *Wemo) scan(urn string, timeout time.Duration) ([]*url.URL, error) {
	// open a udp port for us to receive multicast messages
	udpAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:0", w.ipAddr))
	if err != nil {
		return nil, err
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	defer udpConn.Close()

	//send the
	mAddr, err := net.ResolveUDPAddr("udp", SSDPBROADCAST)
	if err != nil {
		return nil, err
	}

	if w.Debug {
		log.Printf("Found multi-cast address %v", mAddr)
	}
	packet := fmt.Sprintf(MSEARCH, urn)

	if w.Debug {
		log.Printf("Writing discovery packet")
	}
	_, err = udpConn.WriteTo([]byte(packet), mAddr)
	if err != nil {
		return nil, err
	}

	if w.Debug {
		log.Printf("Setting read deadline")
	}
	err = udpConn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, err
	}

	locations := make(map[string]*url.URL)
	for {
		buffer := make([]byte, 2048)
		n, err := udpConn.Read(buffer)
		if err != nil {
			break
		}
		read := string(buffer[:n])
		lines := strings.Split(read, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, LOCATION) {
				temp := strings.TrimSpace(line[len(LOCATION):])
				u, err := url.Parse(temp)
				if err != nil {
					return nil, err
				}
				locations[temp] = u
			}
		}

		if w.Debug {
			log.Printf("Read : %v\n", string(buffer[:n]))
		}
	}

	var results []*url.URL
	for _, value := range locations {
		results = append(results, value)
	}

	return results, nil
}
