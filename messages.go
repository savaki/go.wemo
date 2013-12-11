package wemo

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func post(hostAndPort, action, body string) (*http.Response, error) {
	tcpConn, err := timeoutDialer(2*time.Second, 2*time.Second)("tcp", hostAndPort)
	if err != nil {
		return nil, err
	}
	defer tcpConn.Close()

	preamble := fmt.Sprintf("POST http://%v/upnp/control/basicevent1 HTTP/1.1\r\nContent-type: text/xml; charset=\"utf-8\"\r\nSOAPACTION: \"urn:Belkin:service:basicevent:1#%s\"\r\nContent-Length: %v\r\n\r\n", hostAndPort, action, len(body))
	tcpConn.Write([]byte(preamble + body))

	data, err := ioutil.ReadAll(tcpConn)
	if err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(bytes.NewReader(data)), nil)
}

func newGetBinaryStateMessage() string {
	return `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>
    <u:GetBinaryState xmlns:u="urn:Belkin:service:basicevent:1"></u:GetBinaryState>
  </s:Body>
</s:Envelope>`
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
