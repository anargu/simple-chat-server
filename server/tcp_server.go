package server

import (
	"anargu/simple-chat-server/protocol"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type client struct {
	conn   net.Conn
	name   string
	writer *protocol.CommandWriter
}

// TCPChatServer is the chat server
type TCPChatServer struct {
	listener net.Listener
	clients  []*client
	mutex    *sync.Mutex
}

var (
	// ErrUnknownClient means no client identified
	ErrUnknownClient = errors.New("Unknown client")
)

// NewServer creates a new TCPChatServer
func NewServer() *TCPChatServer {
	return &TCPChatServer{
		mutex: &sync.Mutex{},
	}
}

// Listen in an address
func (s *TCPChatServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)
	if err == nil {
		s.listener = l
	}

	log.Printf("Listening on %v", address)

	return err
}

// Close close connections
func (s *TCPChatServer) Close() {
	s.listener.Close()
}

// Start start a new connection
func (s *TCPChatServer) Start() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
		} else {
			client := s.accept(conn)
			go s.serve(client)
		}
	}
}

// Broadcast send a message to all connected clients
func (s *TCPChatServer) Broadcast(command interface{}) error {
	for _, client := range s.clients {
		client.writer.Write(command)
	}
	return nil
}

// Send sends a message to a client
func (s *TCPChatServer) Send(name string, command interface{}) error {
	for _, client := range s.clients {
		if client.name == name {
			client.writer.Write(command)
		}
	}
	return ErrUnknownClient
}

func (s *TCPChatServer) accept(conn net.Conn) *client {
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(s.clients)+1)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &client{
		conn:   conn,
		writer: protocol.NewCommandWriter(conn),
	}
	s.clients = append(s.clients, client)

	return client
}

func (s *TCPChatServer) remove(client *client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}

	log.Printf("Closing connection from %v", client.conn.RemoteAddr().String())
	client.conn.Close()
}

func (s *TCPChatServer) serve(client *client) {
	cmdReader := protocol.NewCommandReader(client.conn)

	defer s.remove(client)

	for {
		cmd, err := cmdReader.Read()
		if err != nil && err != io.EOF {
			log.Printf("read error: %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.SendCommand:
				go s.Broadcast(protocol.MessageCommand{
					Message: v.Message,
					Name:    client.name,
				})
			case protocol.NameCommand:
				client.name = v.Name
			}
		}

		if err == io.EOF {
			break
		}
	}
}
