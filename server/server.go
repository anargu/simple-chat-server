package server

// ChatServer set the behaviour of a ChatServer
type ChatServer interface {
	Listen(address string) error
	Broadcast(command interface{}) error
	Start()
	Close()
}
