package zpacketgo

type ZPacketDeserializer struct {
}

func (this *ZPacketDeserializer) Read(data []byte) (int, error) {

	if len(data) == 0 {
		//Nothing to read since the buffer has zero length
		return 0, nil
	}

	numberBytesRead := 0

	return numberBytesRead, nil
}
