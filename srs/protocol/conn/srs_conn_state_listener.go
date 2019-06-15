package conn

type SrsConnStateListener interface {
	OnDisconnect() error
}

