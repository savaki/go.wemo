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
	"golang.org/x/net/context"
	"time"
)

func (self *Wemo) foreach(friendlyName string, timeout time.Duration, callback func(*Device)) error {
	ctx := context.Background()

	devices, err := self.DiscoverAll(timeout)
	if err != nil {
		return err
	}

	for _, device := range devices {
		deviceInfo, err := device.FetchDeviceInfo(ctx)
		if err != nil {
			return err
		}

		if deviceInfo.FriendlyName == friendlyName {
			callback(device)
		}
	}

	return nil
}

func (self *Wemo) On(friendlyName string, timeout time.Duration) {
	self.foreach(friendlyName, timeout, func(device *Device) {
		device.On()
	})
}

func (self *Wemo) Off(friendlyName string, timeout time.Duration) {
	self.foreach(friendlyName, timeout, func(device *Device) {
		device.Off()
	})
}

func (self *Wemo) Toggle(friendlyName string, timeout time.Duration) {
	self.foreach(friendlyName, timeout, func(device *Device) {
		device.Toggle()
	})
}
