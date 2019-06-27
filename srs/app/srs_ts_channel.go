package app

type SrsTsChannel struct {
	pid	int
    apply SrsTsPidApply
    stream SrsTsStream
    msg *SrsTsMessage
    context *SrsTsContext
    // for encoder.
    continuityCounter uint8
}

func NewSrsTsChannel() *SrsTsChannel {
	return &SrsTsChannel{}
}