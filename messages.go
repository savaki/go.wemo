package wemo

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
)

func post(hostAndPort, body string) (*http.Response, error) {
	tcpConn, err := net.Dial("tcp", hostAndPort)
	if err != nil {
		return nil, err
	}
	defer tcpConn.Close()

	preamble := fmt.Sprintf("POST http://%v/upnp/control/basicevent1 HTTP/1.1\r\nContent-type: text/xml; charset=\"utf-8\"\r\nSOAPACTION: \"urn:Belkin:service:basicevent:1#SetBinaryState\"\r\nContent-Length: %v\r\n\r\n", hostAndPort, len(body))
	tcpConn.Write([]byte(preamble + body))

	buffer := make([]byte, 2048)
	n, err := tcpConn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(bytes.NewReader(buffer[:n])), nil)
}

func newSetBinaryStateMessage(on bool) string {
	value := 0
	if on {
		value = 1
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>
    <u:SetBinaryState xmlns:u="urn:Belkin:service:basicevent:1">
      <BinaryState>%v</BinaryState>
    </u:SetBinaryState>
  </s:Body>
</s:Envelope>`, value)
}
