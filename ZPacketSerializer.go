package zpacketgo

import "errors"

type packetWriteState int

const (
	writeDestination packetWriteState = iota
	writeSender
	writeDataLength
	writeData
	writeCRC
)

type ZPacketSerializer struct {
	packetQueue   *ZPacketQueue
	packetToSend  *ZPacket
	curWriteState packetWriteState
	curDataI      int
	curCalcedCRC  byte
}

// TODO: Verify that we implemented io.Writer correctly and not as we thought was best
func (this *ZPacketSerializer) Write(data []byte) (int, error) {
	numBytesWritten := 0
	dataLen := len(data)

	if dataLen == 0 {
		return 0, errors.New("Unable to write packet due to empty buffer")
	}

	for i := 0; i < dataLen; i++ {
		//Need to get the next packet to send if we don't have one ready
		if this.packetToSend == nil {
			if this.packetQueue != nil {
				this.packetToSend = this.packetQueue.Pull()

				if this.packetToSend == nil {
					if numBytesWritten == 0 {
						//There is nothing to send and nothing was sent
						return 0, errors.New("No packet currently available to send")
					} else {
						//We actually sent some data but nothing else to send
						return numBytesWritten, nil
					}
				}
			}

		}

		if this.curWriteState == writeDestination {
			this.curCalcedCRC = byte(this.packetToSend.destinationAddress)
			data[i] = byte(this.packetToSend.destinationAddress)
			this.curWriteState = writeSender
		} else if this.curWriteState == writeSender {
			this.XORCRC(byte(this.packetToSend.senderAddress))
			data[i] = byte(this.packetToSend.senderAddress)
			this.curWriteState = writeDataLength
		} else if this.curWriteState == writeDataLength {
			this.XORCRC(byte(dataLen))
			data[i] = byte(dataLen)

			if dataLen == 0 {
				//If there is no data then jump straight to writing crc instead of sending data since there is no data
				this.curWriteState = writeCRC
			} else {
				this.curWriteState = writeData
			}
		} else if this.curWriteState == writeData {
			this.XORCRC(this.packetToSend.data[this.curDataI])
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
			panic(errors.New("Panic as ZPacketSerializer is an unknown state which should never have happened"))
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

func (this *ZPacketSerializer) AddPacketToSend(p *ZPacket) error {
	if p == nil {
		return errors.New("Unable to send nil ZPacket")
	}
	//We are already sending a packet lets not interrupt
	if this.packetQueue == nil {
		this.packetQueue = &ZPacketQueue{}
	}

	this.packetQueue.Push(p)

	return nil
}

func (this *ZPacketSerializer) AvailableToSendPacket() bool {
	return this.packetToSend == nil
}

func (this *ZPacketSerializer) XORCRC(d byte) {
	this.curCalcedCRC ^= d
}
