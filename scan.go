package wemo

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

const (
	SSDP_BROADCAST = "239.255.255.250:1900"
	M_SEARCH       = "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 10\r\nST: %s\r\nUSER-AGENT: unix/5.1 UPnP/1.1 crash/1.0\r\n\r\n"
	LOCATION       = "LOCATION: "
)

// scan the multicast
func (self *Wemo) scan(urn string, timeout time.Duration) ([]*url.URL, error) {
	// open a udp port for us to receive multicast messages
	udpAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:0", self.ipAddr))
	if err != nil {
		return nil, err
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	defer udpConn.Close()

	//send the
	mAddr, err := net.ResolveUDPAddr("udp", SSDP_BROADCAST)
	if err != nil {
		return nil, err
	}

	if self.Debug {
		log.Printf("Found multi-cast address %v", mAddr)
	}
	packet := fmt.Sprintf(M_SEARCH, urn)

	if self.Debug {
		log.Printf("Writing discovery packet")
	}
	_, err = udpConn.WriteTo([]byte(packet), mAddr)
	if err != nil {
		return nil, err
	}

	if self.Debug {
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

		if self.Debug {
			log.Printf("Read : %v\n", string(buffer[:n]))
		}
	}

	var results []*url.URL
	for _, value := range locations {
		results = append(results, value)
	}

	return results, nil
}
