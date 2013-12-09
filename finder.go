package wemo

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"time"
)

type Wemo struct {
	ipAddr string
	Debug  bool
}

var re *regexp.Regexp = regexp.MustCompile(`http://([^/]+)/setup.xml`)

func (self *Wemo) FindAll(timeout time.Duration) ([]*Device, error) {
	var all []*Device
	urns := []string{"urn:Belkin:device:controllee:1", "urn:Belkin:device:light:1", "urn:Belkin:device:sensor:1"}
	for _, urn := range urns {
		devices, _ := self.Find(urn, timeout)
		for _, device := range devices {
			all = append(all, device)
		}
	}

	return all, nil
}

func (self *Wemo) Find(urn string, timeout time.Duration) ([]*Device, error) {
	locations, err := self.scan(urn, timeout)
	if err != nil {
		return nil, err
	}

	var devices []*Device
	for _, uri := range locations {
		if matches := re.FindStringSubmatch(uri.String()); len(matches) == 2 {
			device := &Device{matches[1]}
			devices = append(devices, device)
		}
	}

	return devices, nil
}
