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

func Foo(name string) {
	body := newSetBinaryStateMessage(false)
	response, err := post("10.0.1.32:49153", body)
	if err != nil {
		log.Fatalf("unable to execute request => %s\n", err.Error())
	}
	defer response.Body.Close()
	fmt.Printf("response code => %d\n", response.StatusCode)
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(contents))
}

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
