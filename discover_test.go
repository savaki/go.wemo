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
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"time"
)

func TestDisoverAll(t *testing.T) {
	if os.Getenv("integration") != "" {
		// device := &Device{Host: "10.0.1.32:49153"}
		// device := &Device{Host: "10.0.1.19:49153"}

		// device.Toggle()
		// <-time.After(1 * time.Second)

		// device.Toggle()
		// <-time.After(1 * time.Second)

		// device.Toggle()
		// <-time.After(1 * time.Second)

		// device.Toggle()
		// <-time.After(1 * time.Second)
		api, _ := NewByInterface("en0")
		devices, _ := api.DiscoverAll(3 * time.Second)
		for _, device := range devices {
			fmt.Printf(">> %+v\n", device)
		}
	}
}

func TestRENoMatch(t *testing.T) {
	Convey("Given a url that does not contain setup.xml", t, func() {
		var url = "http://10.0.11:80/description.xml"

		Convey("When I call belkinRE.FindStringSubmatch", func() {
			matches := belkinRE.FindStringSubmatch(url)

			Convey("Then I expect no errors", func() {
				So(len(matches), ShouldEqual, 0)
			})
		})
	})
}

func TestREMatch(t *testing.T) {
	Convey("Given a URL", t, func() {
		var url = "http://10.0.1.17:49153/setup.xml"

		Convey("When I call belkinRE.FindStringSubmatch", func() {
			// When
			matches := belkinRE.FindStringSubmatch(url)

			Convey("Then I expect 2 matches", func() {
				So(len(matches), ShouldEqual, 2)
			})

			Convey("And I expect a host match", func() {
				So(matches[1], ShouldEqual, "10.0.1.17:49153")
			})
		})
	})
}
