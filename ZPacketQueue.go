package zpacketgo

type ZPacketQueue struct {
	head *ZPQNode
	tail *ZPQNode
}

type ZPQNode struct {
	next *ZPQNode
	zp   *ZPacket
}

func (this *ZPacketQueue) Pull() *ZPacket {
	if this.head == nil {
		return nil
	}

	zpNode := this.head

	this.head = zpNode.next

	//If we pulled the last packet from the queue then we also need to clear tail
	if this.head == nil {
		this.tail = nil
	}

	return zpNode.zp
}

func (this *ZPacketQueue) Push(zp *ZPacket) {
	zpNode := &ZPQNode{zp: zp}
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
