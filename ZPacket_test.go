package zpacketgo

import "testing"

// ------------------------------ Happy Path Tests ---------------------------------
func TestThatNewZPacketWithValidAddressesButNilDataCreatesWithoutError(t *testing.T) {
	zp, err := NewZPacket(ZPacketAddress(0x03), ZPacketAddress(0x01), nil)

	if err != nil {
		t.Errorf("Creating a packet with valid address and no data is not supposed to return an error %v", err)
	}

	if zp == nil {
		t.Error("Creating a packet with valid address and no data is supposed to return a valid")
	}
}

func TestThatNewZPacketWithValidAddressesAnd1ByteDataCreatesWithoutError(t *testing.T) {
	zp, err := NewZPacket(ZPacketAddress(0x03), ZPacketAddress(0x01), make([]byte, 1))

	if err != nil {
		t.Errorf("Creating a packet with valid address and one byte of data is not supposed to return an error %v", err)
	}

	if zp == nil {
		t.Error("Creating a packet with valid address and one byte of is supposed to return a valid")
	}
}

func TestThatNewZPacketWithValidAddressesAndMaxSizeDataCreatesWithoutError(t *testing.T) {
	zp, err := NewZPacket(ZPacketAddress(0x03), ZPacketAddress(0x01), make([]byte, DATA_MAX_LENGTH))

	if err != nil {
		t.Errorf("Creating a packet with valid address and %d bytes of data is not supposed to return an error %v", DATA_MAX_LENGTH, err)
	}

	if zp == nil {
		t.Errorf("Creating a packet with valid address and %d bytes of data is supposed to return a valid", DATA_MAX_LENGTH)
	}
}

// ------------------------------- NewPacket Error Tests ----------------------------
func TestThatNewZPacketWithValidAddressesAndMaxSizePlus1DataDoesNotCreatesWithError(t *testing.T) {
	zp, err := NewZPacket(ZPacketAddress(0x03), ZPacketAddress(0x01), make([]byte, DATA_MAX_LENGTH+1))

	if err == nil {
		t.Errorf("Creating a packet with valid address and %d bytes of data is supposed to return an error of too long", DATA_MAX_LENGTH+1)
	}

	if zp != nil {
		t.Errorf("Creating a packet with valid address and %d bytes of data is supposed to return nil for ZPackewt", DATA_MAX_LENGTH+1)
	}
}

func TestThatNewZPacketWithOutOfRangeDestinationAddressReturnsErrorAndNoZPacket(t *testing.T) {
	zp, err := NewZPacket(ZPacketAddress(MAX_ADDRESS+1), ZPacketAddress(0), nil)

	if err == nil {
		t.Fail()
	}

	if zp != nil {
		t.Fail()
	}
}

func TestThatNewZPacketWithOutOfRangeSenderAddressReturnsErrorAndNoZPacket(t *testing.T) {
	zp, err := NewZPacket(ZPacketAddress(0), ZPacketAddress(MAX_ADDRESS+1), nil)

	if err == nil {
		t.Fail()
	}

	if zp != nil {
		t.Fail()
	}
}

// -------------------------------- Address Getters --------------------------------
func TestThatZPacketReturnsDestinationAddressThatItWasCreatedWith(t *testing.T) {
	destinationAddress := ZPacketAddress(4)
	zp, _ := NewZPacket(destinationAddress, ZPacketAddress(0), nil)

	if zp.DestinationAddress() != destinationAddress {
		t.Fail()
	}
}

func TestThatZPacketReturnsSenderAddressThatItWasCreatedWith(t *testing.T) {
	senderAddress := ZPacketAddress(4)
	zp, _ := NewZPacket(ZPacketAddress(0), senderAddress, nil)

	if zp.SenderAddress() != senderAddress {
		t.Fail()
	}
}

// ------------------------------- Data Getters -------------------------------------
func TestThatZPacketReturnsEmptySliceForDataWhenNilWasPassedInForData(t *testing.T) {
	zp, _ := NewZPacket(ZPacketAddress(0), ZPacketAddress(1), nil)

	if zp.Data() == nil || len(zp.Data()) != 0 {
		t.Fail()
	}
}

func TestThatZPacketDataReturnsMatchingValuesForWhatWasPassedIn(t *testing.T) {
	testData := []byte{0x00, 0xAA, 0x55, 0xFF}
	zp, _ := NewZPacket(ZPacketAddress(0), ZPacketAddress(1), testData)

	for i, curZPDataByte := range zp.Data() {
		if testData[i] != curZPDataByte {
			t.Fail()
		}
	}
}

func TestThatZPacketCopiesDataOnCreationAndNotCreatesReferenceToPassedInSlice(t *testing.T) {
	testData := make([]byte, 0, DATA_MAX_LENGTH)

	//First set the input data to something
	for i, _ := range testData {
		testData[i] = 0xAA
	}

	zp, _ := NewZPacket(ZPacketAddress(0), ZPacketAddress(1), testData)

	//Then set the input data to something else after the zpacket is created
	for i, _ := range testData {
		testData[i] = 0xF4
	}

	for i, curZPDataByte := range zp.Data() {
		//Due to the change we expect all data to be different if copy occured
		if testData[i] == curZPDataByte {
			t.Fail()
		}
	}
}
