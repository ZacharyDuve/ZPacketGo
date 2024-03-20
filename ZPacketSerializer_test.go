package zpacketgo

import "testing"

func TestThatCallingSerializePacketOnNilZPacketSerializerReturnsNilZPacketSerializerReferenceError(t *testing.T) {
	var zps *ZPacketSerializer = nil

	zp, _ := NewZPacket(ZPacketAddress(1), ZPacketAddress(2), nil)

	if !IsZPacketSerializerNilReferenceError(zps.SerializePacket(zp)) {
		t.Fail()
	}
}

func TestThatCallingSerializePacketWithNilPacketReturnsError(t *testing.T) {
	zps := ZPacketSerializer{}

	if err := zps.SerializePacket(nil); err == nil {
		t.Fail()
	}
}
