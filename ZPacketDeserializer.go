package zpacketgo

import "errors"

type zPacketDeserializerState int

const (
	readDestinationAddress zPacketDeserializerState = iota
	readSenderAddress
	readDataLength
	readData
	readCRC

	nilReferenceErrorMessage string = "ZPacketDeserializer Error: Read called on a nil reference to ZPacketDeserializer"
)

type ZPacketDeserializer struct {
	zpQueue              zPacketQueue
	curState             zPacketDeserializerState
	curPacketDestination byte
	curPacketSender      byte
	curPacketDataLength  byte
	curPacketData        []byte
	curPacketDataIndex   int
	curPacketCalcedCRC   byte
}

// Implementation of io.Reader
func (this *ZPacketDeserializer) Read(data []byte) (int, error) {

	if this == nil {
		return 0, errors.New(nilReferenceErrorMessage)
	}

	// If this curPacketData is zero due to first run through then we need to initialize the curPacketData
	if len(this.curPacketData) == 0 {
		//Initialize it to the full length as we are going to ensure that the backing array is only ever allocated once
		this.curPacketData = make([]byte, DATA_MAX_LENGTH)
	}

	if len(data) == 0 {
		//Nothing to read since the buffer has zero length
		return 0, nil
	}

	numberBytesRead := 0

	for _, curByte := range data {
		switch this.curState {
		case readDestinationAddress:
			{
				this.curPacketDestination = curByte
				this.curPacketCalcedCRC = curByte
				this.curState = readSenderAddress
			}
		case readSenderAddress:
			{
				this.curPacketSender = curByte
				this.curPacketCalcedCRC ^= curByte
				this.curState = readDataLength
			}
		case readDataLength:
			{
				this.curPacketDataLength = curByte
				this.curPacketCalcedCRC ^= curByte

				if this.curPacketDataLength == 0 {
					//If the data length is 0 then there is no data so skip to read crc
					this.curState = readCRC
				} else {
					//Otherwise there is data that we need to read
					this.curState = readData
				}
			}
		case readData:
			{
				this.curPacketData[this.curPacketDataIndex] = curByte
				this.curPacketDataIndex++
				this.curPacketCalcedCRC ^= curByte

				if this.curPacketDataIndex == int(this.curPacketDataLength) {
					//If the dataIndex has been incremented to the length then we have read all data and can move onto read crc
					this.curState = readCRC
				}
			}
		case readCRC:
			{
				if this.curPacketCalcedCRC == curByte {
					zp, err := NewZPacket(ZPacketAddress(this.curPacketDestination), ZPacketAddress(this.curPacketSender), this.curPacketData[:this.curPacketDataLength])
					if err != nil {
						//Error encountered so lets stop and send out the error and how many bytes we read
						return numberBytesRead, err
					}
					//Otherwise lest push the packet onto the queue and continue
					this.zpQueue.Push(zp)
				}
			}
		}
	}

	return numberBytesRead, nil
}

// This function will return a slice of *ZPackets unless there are none then nil is returned
func (this *ZPacketDeserializer) Packets() []*ZPacket {

	if this == nil || this.zpQueue.Len() == 0 {
		//Nothing here return nil
		return nil
	}

	packets := make([]*ZPacket, 0, this.zpQueue.Len())

	for curPacket := this.zpQueue.Pull(); curPacket != nil; curPacket = this.zpQueue.Pull() {
		packets = append(packets, curPacket)
	}

	return packets
}

func IsZPacketDeserializerNilReferenceError(err error) bool {
	return err != nil && err.Error() == nilReferenceErrorMessage
}
