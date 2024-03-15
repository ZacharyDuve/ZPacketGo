package zpacketgo

type zPacketQueue struct {
	head   *zPQNode
	tail   *zPQNode
	length int
}

type zPQNode struct {
	next *zPQNode
	zp   *ZPacket
}

func (this *zPacketQueue) Pull() *ZPacket {
	if this.head == nil {
		return nil
	}

	zpNode := this.head

	//Make the head the next packet in the queue
	this.head = zpNode.next
	//Decrease length since we pulled a packet off
	this.length--

	//If we pulled the last packet from the queue then we also need to clear tail
	if this.head == nil {
		this.tail = nil
	}

	return zpNode.zp
}

func (this *zPacketQueue) Push(zp *ZPacket) {
	zpNode := &zPQNode{zp: zp}
	//Need to increment the length as we are adding in a packet
	this.length++

	if this.tail == nil {
		//if we don't have a tail then first one becomes tail and head
		this.tail = zpNode
		this.head = zpNode
	} else {
		//otherwise add to tail as normal
		this.tail.next = zpNode
		this.tail = zpNode
	}
}

func (this *zPacketQueue) Len() int {
	return this.length
}
