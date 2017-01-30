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
	"regexp"
	"time"
)

//Constants for URNS
const (
	Basic      = "urn:Belkin:service:basicevent:1"
	Bridge     = "urn:Belkin:device:bridge:1"
	Controllee = "urn:Belkin:device:controllee:1"
	Light      = "urn:Belkin:device:light:1"
	Sensor     = "urn:Belkin:device:sensor:1"
	NetCam     = "urn:Belkin:device:netcam:1"
	Insight    = "urn:Belkin:device:insight:1"
)

var (
	//var belkinRE *regexp.Regexp = regexp.MustCompile(`http://([^/]+)/setup.xml`)
	belkinRE = regexp.MustCompile(`http://([^/]+)/setup.xml`)
)

// Wemo ...
type Wemo struct {
	ipAddr string
	Debug  bool
}

// DiscoverAll ...
func (w *Wemo) DiscoverAll(timeout time.Duration) ([]*Device, error) {
	urns := []string{Basic, Bridge, Controllee, Light, Sensor, NetCam, Insight}
	var all []*Device
	for _, urn := range urns {
		devices, _ := w.Discover(urn, timeout)
		for _, device := range devices {
			all = append(all, device)
		}
	}

	return all, nil
}

// Discover ...
func (w *Wemo) Discover(urn string, timeout time.Duration) ([]*Device, error) {
	locations, err := w.scan(urn, timeout)
	if err != nil {
		return nil, err
	}

	var devices []*Device
	for _, uri := range locations {
		if matches := belkinRE.FindStringSubmatch(uri.String()); len(matches) == 2 {
			host := matches[1]
			devices = append(devices, &Device{Host: host})
		}
	}
	return devices, nil
}
