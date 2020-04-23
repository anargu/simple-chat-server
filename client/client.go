package client

import "anargu/simple-chat-server/protocol"

type messageHandler func(string)

// ChatClient client
type ChatClient interface {
	Dial(address string) error
	Start()
	Close()
	Send(command interface{}) error
	SetName(name string) error
	SendMessage(message string) error
	Incoming() chan protocol.MessageCommand
}
