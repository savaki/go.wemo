package wemo

import (
	"fmt"
	"testing"
	"time"
)

func TestSample(t *testing.T) {
	// device := &Device{Host: "10.0.1.32:49153"}
	// device := &Device{Host: "10.0.1.19:49153"}
	// device.On()
	api, _ := NewByInterface("en0")
	devices, _ := api.FindAll(3 * time.Second)
	for _, device := range devices {
		fmt.Printf(">> %+v\n", device)
	}
}

func TestRENoMatch(t *testing.T) {
	var url = "http://10.0.11:80/description.xml"

	// When
	matches := re.FindStringSubmatch(url)

	// Then
	if len(matches) != 0 {
		t.Fail()
	}
}

func TestREMatch(t *testing.T) {
	var url = "http://10.0.1.17:49153/setup.xml"

	// When
	matches := re.FindStringSubmatch(url)

	// Then
	if len(matches) != 2 {
		t.Fail()
	}
	if matches[1] != "10.0.1.17:49153" {
		t.Fail()
	}
}
