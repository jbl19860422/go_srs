package protocol

type SrsHandshake interface {
	HandShakeWithClient() int
}