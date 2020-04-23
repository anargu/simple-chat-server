package client

import (
	"anargu/simple-chat-server/protocol"
	"io"
	"log"
	"net"
)

// TCPChatClient tcp client
type TCPChatClient struct {
	conn      net.Conn
	cmdReader *protocol.CommandReader
	cmdWriter *protocol.CommandWriter
	name      string
	incoming  chan protocol.MessageCommand
}

// NewClient creates a new
func NewClient() *TCPChatClient {
	return &TCPChatClient{
		incoming: make(chan protocol.MessageCommand),
	}
}

// Dial call method
func (c *TCPChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)

	if err == nil {
		c.conn = conn
	}

	c.cmdReader = protocol.NewCommandReader(conn)
	c.cmdWriter = protocol.NewCommandWriter(conn)

	return err
}

// Start start method
func (c *TCPChatClient) Start() {
	for {
		cmd, err := c.cmdReader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Read error %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.MessageCommand:
				c.incoming <- v
			default:
				log.Printf("Unknown command: %v", v)
			}
		}
	}
}

// Close closes connection
func (c *TCPChatClient) Close() {
	c.conn.Close()
}

// Incoming method
func (c *TCPChatClient) Incoming() chan protocol.MessageCommand {
	return c.incoming
}

// Send write  by command
func (c *TCPChatClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

// SetName set a name for the chat
func (c *TCPChatClient) SetName(name string) error {
	return c.Send(protocol.NameCommand{Name: name})
}

// SendMessage send message for the chat
func (c *TCPChatClient) SendMessage(message string) error {
	return c.Send(protocol.SendCommand{
		Message: message,
	})
}
