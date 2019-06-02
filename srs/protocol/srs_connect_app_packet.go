package protocol

type SrsConnectAppPacket struct {
	/**
	 * Name of the command. Set to "connect".
	 */
	command_name string
	/**
	 * Always set to 1.
	 */
	transaction_id float64
}

func (s *SrsConnectAppPacket) Encode(payload []byte, size int32) int32 {
	return 0
}
