package wemo

import (
	"regexp"
	"time"
)

var belkinRE *regexp.Regexp = regexp.MustCompile(`http://([^/]+)/setup.xml`)

type Wemo struct {
	ipAddr string
	Debug  bool
}

func (self *Wemo) DiscoverAll(timeout time.Duration) ([]*Device, error) {
	var all []*Device
	urns := []string{"urn:Belkin:device:controllee:1", "urn:Belkin:device:light:1", "urn:Belkin:device:sensor:1"}
	for _, urn := range urns {
		devices, _ := self.Discover(urn, timeout)
		for _, device := range devices {
			all = append(all, device)
		}
	}

	return all, nil
}

func (self *Wemo) Discover(urn string, timeout time.Duration) ([]*Device, error) {
	locations, err := self.scan(urn, timeout)
	if err != nil {
		return nil, err
	}

	var devices []*Device
	for _, uri := range locations {
		if matches := belkinRE.FindStringSubmatch(uri.String()); len(matches) == 2 {
			host := matches[1]
			devices = append(devices, &Device{host})
		}
	}

	return devices, nil
}
