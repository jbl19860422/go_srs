package protocol

type SrsChunkStream struct {
	fmt byte
	cid int32
	header SrsMessageHeader	
	extended_timestamp bool 
	msg_count int32
}
