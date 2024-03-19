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

//-------------------------- Read ---------------------------------------

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
