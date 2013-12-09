package wemo

import (
	"errors"
	"log"
	"net"
	"regexp"
)

var ipAddrRE = regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/\d{1,3}$`)

func NewByIp(ipAddr string) *Wemo {
	return &Wemo{ipAddr: ipAddr, Debug: false}
}

// find the ip address associated with the specified interface
func NewByInterface(name string) (*Wemo, error) {
	// find the interface with the selected name
	iface, err := net.InterfaceByName(name)
	if err != nil {
		log.Printf("Unable to find interface, %s\n", name)
		return nil, err
	}

	// find all the addresses associated with this address
	addrs, err := iface.Addrs()
	if err != nil {
		log.Printf("No addresses associated with interface, %s\n", iface.Name)
		return nil, err
	}

	// and find the one that looks like an IPv4 address
	for _, addr := range addrs {
		if matches := ipAddrRE.FindStringSubmatch(addr.String()); len(matches) == 2 {
			return NewByIp(matches[1]), nil
		}
	}

	// nope, couldn't find one
	return nil, errors.New("unable to find ip address associated with interface, " + name)
}
