package wemo

import (
	"fmt"
	"testing"
)

var (
	testMessageHeader = `<?xml version="1.0" encoding="utf-8"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body>`
	testMessageFooter = `</s:Body></s:Envelope>`
)

func TestNewGetBinaryStateMessage(t *testing.T) {
	expected := `<?xml version="1.0" encoding="utf-8"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body><u:GetBinaryState xmlns:u="urn:Belkin:service:basicevent:1"></u:GetBinaryState></s:Body></s:Envelope>`
	actual := newGetBinaryStateMessage()
	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestNewSetBinaryStateMessage(t *testing.T) {
	t.Run("Binary State: On", testNewSetBinaryStateMessage(true))
	t.Run("Binary State: Off", testNewSetBinaryStateMessage(false))
}

func testNewSetBinaryStateMessage(state bool) func(*testing.T) {
	v := 0
	if state {
		v = 1
	}

	msg := `<u:SetBinaryState xmlns:u="urn:Belkin:service:basicevent:1"><BinaryState>%v</BinaryState></u:SetBinaryState>`

	expected := fmt.Sprintf(testMessageHeader+msg+testMessageFooter, v)

	return func(t *testing.T) {
		actual := newSetBinaryStateMessage(state)
		if actual != expected {
			t.Errorf("Expected: %s, got: %s", expected, actual)
		}
	}
}

func TestNewGetInsightParamsMessage(t *testing.T) {
	expected := `<?xml version="1.0" encoding="utf-8"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body><u:GetInsightParams xmlns:u="urn:Belkin:service:insight:1"></u:GetInsightParams></s:Body></s:Envelope>`
	actual := newGetInsightParamsMessage()
	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestNewGetBridgeEndDevices(t *testing.T) {
	udn := "uuid:Bridge-1_0-231503B01005A4"
	expected := fmt.Sprintf(testMessageHeader+`<u:GetEndDevices xmlns:u="urn:Belkin:service:bridge:1"><DevUDN>%s</DevUDN><ReqListType>PAIRED_LIST</ReqListType></u:GetEndDevices>`+testMessageFooter, udn)
	actual := newGetBridgeEndDevices(udn)
	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestNewSetBulbStatus(t *testing.T) {
	id := "94103EF6BF42867F"

	t.Run("Bulb State: On", testNewSetBulbStatus(id, "10006", "1", false))
	t.Run("Bulb State: Off", testNewSetBulbStatus(id, "10006", "0", false))

	t.Run("Bulb State: On", testNewSetBulbStatus(id, "10006", "1", true))
	t.Run("Bulb State: Off", testNewSetBulbStatus(id, "10006", "0", true))

	t.Run("Bulb State: Dim", testNewSetBulbStatus(id, "10008", "0", false))
	t.Run("Bulb State: Dim", testNewSetBulbStatus(id, "10008", "255", false))

	t.Run("Bulb State: Dim", testNewSetBulbStatus(id, "10008", "0", true))
	t.Run("Bulb State: Dim", testNewSetBulbStatus(id, "10008", "255", true))
}

func testNewSetBulbStatus(id, capability, value string, group bool) func(*testing.T) {
	g := "NO"
	if group {
		g = "YES"
	}

	msg := `<u:SetDeviceStatus xmlns:u="urn:Belkin:service:bridge:1">
			<DeviceStatusList>
		&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;&lt;DeviceStatus&gt;&lt;IsGroupAction&gt;%s&lt;/IsGroupAction&gt;&lt;DeviceID available=&quot;YES&quot;&gt;%s&lt;/DeviceID&gt;&lt;CapabilityID&gt;%s&lt;/CapabilityID&gt;&lt;CapabilityValue&gt;%s&lt;/CapabilityValue&gt;&lt;/DeviceStatus&gt;
		</DeviceStatusList>
	</u:SetDeviceStatus>`

	expected := fmt.Sprintf(testMessageHeader+msg+testMessageFooter, g, id, capability, value)

	return func(t *testing.T) {
		actual := newSetBulbStatus(id, capability, value, group)
		if actual != expected {
			t.Errorf("Expected: %s, got: %s", expected, actual)
		}
	}
}

func TestNewGetBulbStatus(t *testing.T) {

	id := "94103EF6BF42867F"

	msg := `<u:GetDeviceStatus xmlns:u="urn:Belkin:service:bridge:1">
		<DeviceIDs>%s</DeviceIDs>
		</u:GetDeviceStatus>`

	expected := fmt.Sprintf(testMessageHeader+msg+testMessageFooter, id)

	actual := newGetBulbStatus(id)
	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}
