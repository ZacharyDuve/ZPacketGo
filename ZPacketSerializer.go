package zpacketgo

import "errors"

type packetWriteState int

const (
	writeDestination packetWriteState = iota
	writeSender
	writeDataLength
	writeData
	writeCRC

	nilSerializerReferenceErrorMessage string = "error Read called on a nil reference to ZPacketDeserializer"
)

type ZPacketSerializer struct {
	packetQueue   zPacketQueue
	packetToSend  *ZPacket
	curWriteState packetWriteState
	curDataI      int
	curCalcedCRC  byte
}

// TODO: Verify that we implemented io.Writer correctly and not as we thought was best
func (this *ZPacketSerializer) Write(data []byte) (int, error) {

	if this == nil {
		//Our ref
		return 0, nil
	}

	numBytesWritten := 0
	dataLen := len(data)

	//No Need to error check for len(data) == 0 as that is ok as defined by io.Reader. In that case return 0 and nil
	for i := 0; i < dataLen; i++ {
		//Need to get the next packet to send if we don't have one ready
		if this.packetToSend == nil {
			this.packetToSend = this.packetQueue.Pull()

			if this.packetToSend == nil {
				if numBytesWritten == 0 {
					//There is nothing to send and nothing was sent
					return 0, errors.New("no packet currently available to send")
				} else {
					//We actually sent some data but nothing else to send
					return numBytesWritten, nil
				}
			}

		}

		if this.curWriteState == writeDestination {
			this.curCalcedCRC = byte(this.packetToSend.destinationAddress)
			data[i] = byte(this.packetToSend.destinationAddress)
			this.curWriteState = writeSender
		} else if this.curWriteState == writeSender {
			this.curWriteStxorateCRC(byte(this.packetToSend.senderAddress))
			data[i] = byte(this.packetToSend.senderAddress)
			this.curWriteState = writeDataLength
		} else if this.curWriteState == writeDataLength {
			this.curWriteStxorateCRC(byte(dataLen))
			data[i] = byte(dataLen)

			if dataLen == 0 {
				//If there is no data then jump straight to writing crc instead of sending data since there is no data
				this.curWriteState = writeCRC
			} else {
				this.curWriteState = writeData
			}
		} else if this.curWriteState == writeData {
			this.curWriteStxorateCRC(this.packetToSend.data[this.curDataI])
			data[i] = this.packetToSend.data[this.curDataI]
			this.curDataI++

			if this.curDataI == len(this.packetToSend.data) {
				//If we have written all of the packets data then we need to write the crc
				this.curWriteState = writeCRC
			}
		} else if this.curWriteState == writeCRC {
			data[i] = this.curCalcedCRC
			this.reset()

		} else {
			panic(errors.New("panic as ZPacketSerializer is an unknown state which should never have happened"))
		}

		numBytesWritten++
	}

	//If we got here then we should have been able to fill the buffer up
	return numBytesWritten, nil
}

func (this *ZPacketSerializer) reset() {
	this.curDataI = 0
	this.curWriteState = writeDestination
	this.packetToSend = nil
}

func (this *ZPacketSerializer) SerializePacket(p *ZPacket) error {
	if p == nil {
		return errors.New("unable to send nil ZPacket")
	}
	//We are already sending a packet lets not interrupt
	if this.packetQueue == nil {
		this.packetQueue = &zPacketQueue{}
	}

	this.packetQueue.Push(p)

	return nil
}

func (this *ZPacketSerializer) curWriteStxorateCRC(d byte) {
	this.curCalcedCRC ^= d
}
