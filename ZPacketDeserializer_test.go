package zpacketgo

import "testing"

// --------------------------- Creation State --------------------------------------
func TestThatZPacketDeserializerStartsInReadDestionationState(t *testing.T) {
	zpd := ZPacketDeserializer{}

	if zpd.curState != readDestinationAddress {
		t.Errorf("Expected ZPacketDeserializer to start in %v state but instead it started in %v", readDestinationAddress, zpd.curState)
	}
}

//--------------------------- Calling functions on nil reference ------------------

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

//-------------------------- Read Edge Cases ---------------------------------------

func TestThatCallingReadWithNilDataInputReturns0ForBytesReadAndNilForError(t *testing.T) {
	zpd := ZPacketDeserializer{}

	numBytesRead, err := zpd.Read(nil)

	if numBytesRead != 0 {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestThatCallingReadWithZeroLengthDataInputReturns0ForBytesReadAndNilForError(t *testing.T) {
	zpd := ZPacketDeserializer{}

	numBytesRead, err := zpd.Read(make([]byte, 0))

	if numBytesRead != 0 {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

//----------------------- Read Actual Data -------------------------

func TestDeserializeOfWholePacketWithNoDataReturnsSinglePacket(t *testing.T) {
	destAddr := 0x01
	sendAddr := 0x02
	data := []byte{0x81, 0x02, 0x00, 0x83}

	zpd := ZPacketDeserializer{}

	numBytesRead, err := zpd.Read(data)

	if numBytesRead != len(data) {
		t.Fatalf("Expected %d bytes to be read but instead only %d were read. The Deserializer needs to be greedy for this test to work", len(data), numBytesRead)
	}

	if err != nil {
		t.Fatalf("Read of valid packet is not supposed to return error but instead one was returned %s", err)
	}

	packets := zpd.Packets()

	if len(packets) == 0 {
		t.Fatal("Expected deserializer to return a packet after reading an entire packet but none were returned")
	}

	if len(packets) > 1 {
		t.Fatalf("Only expected one packet to be read but instead %d packets were deserialized", len(packets))
	}

	if packets[0].DestinationAddress() != ZPacketAddress(destAddr) {
		t.Fatalf("Destination address of %d does not match expected %d", packets[0].DestinationAddress(), destAddr)
	}

	if packets[0].SenderAddress() != ZPacketAddress(sendAddr) {
		t.Fatalf("Sender address of %d does not match expected %d", packets[0].SenderAddress(), sendAddr)
	}

	if len(packets[0].Data()) != 0 {
		t.Fatal("Packet data should have been 0 bytes long but instead it was", len(packets[0].Data()))
	}
}
