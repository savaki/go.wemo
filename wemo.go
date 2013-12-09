package wemo

import (
	"time"
)

func (self *Wemo) foreach(friendlyName string, timeout time.Duration, callback func(*Device)) error {
	devices, err := self.DiscoverAll(timeout)
	if err != nil {
		return err
	}

	for _, device := range devices {
		deviceInfo, err := device.FetchDeviceInfo()
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
