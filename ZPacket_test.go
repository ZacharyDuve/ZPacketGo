package zpacketgo

import "testing"

func TestThatNewZPacketWithValidAddressesButNilDataCreatesWithoutError(t *testing.T) {
	zp, err := NewZPacket(ZPacketAddress(0x03), ZPacketAddress(0x01), nil)

	if err == nil {

	}
}
