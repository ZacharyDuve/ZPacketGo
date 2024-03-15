package zpacketgo

import (
	"fmt"
)

type ZPacketAddress byte

const (
	DATA_MAX_LENGTH = 255
	MAX_ADDRESS     = 63
)

type ZPacket struct {
	destinationAddress ZPacketAddress
	senderAddress      ZPacketAddress
	data               []byte
}

func NewZPacket(destination, sender ZPacketAddress, d []byte) (*ZPacket, error) {
	if destination > MAX_ADDRESS {
		return nil, fmt.Errorf("Unable to create ZPacket due to Destination of %d being greater than max address of %d", destination, MAX_ADDRESS)
	}

	if sender > MAX_ADDRESS {
		return nil, fmt.Errorf("Unable to create ZPacket due to Sender of %d being greater than max address of %d", sender, MAX_ADDRESS)
	}

	if len(d) > DATA_MAX_LENGTH {
		return nil, fmt.Errorf("Unable to create ZPacket due to data length of %d being greater than max data length %d", len(d), DATA_MAX_LENGTH)
	}

	p := &ZPacket{destinationAddress: destination, senderAddress: sender}
	//Make a deep copy of the data to protect from modification
	p.data = make([]byte, len(d))

	if d != nil {
		copy(p.data, d)
	}

	return p, nil
}

func (this *ZPacket) DestinationAddress() ZPacketAddress {
	return this.destinationAddress
}

func (this *ZPacket) SenderAddress() ZPacketAddress {
	return this.senderAddress
}

func (this *ZPacket) Data() []byte {
	return this.data
}
