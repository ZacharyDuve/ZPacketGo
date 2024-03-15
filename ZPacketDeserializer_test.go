package zpacketgo

import "testing"

func TestThatZPacketDeserializerStartsInReadDestionationState(t *testing.T) {
	zpd := ZPacketDeserializer{}

	if zpd.curState != readDestinationAddress {
		t.Errorf("Expected ZPacketDeserializer to start in %v state but instead it started in %v", readDestinationAddress, zpd.curState)
	}
}

func TestThatNilReferenceToZPacketDeserializerWhenReadCalledReturnsAppropriateError(t *testing.T) {
	var zpd *ZPacketDeserializer = nil
	_, err := zpd.Read(make([]byte, 0))

	if !IsZPacketDeserializerNilReferenceError(err) {
		t.Fail()
	}
}

func TestThatNilReferenceToZPacketDeserializerWhenPacketsIsCalledReturnsNil(t *testing.T) {
	var zpd *ZPacketDeserializer = nil

	if zpd.Packets() != nil {
		t.Fail()
	}
}
