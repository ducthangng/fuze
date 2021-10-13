package chat

import (
	"fuze/srv/chat/protocol"
	"log"
	"net"
)

type ChatClient interface {
	Dial(address string) error
	Start()
	Close()
	Send(command interface{}) error
	SetName(name string) error
	SendMessage(message string) error
	Error() chan error
	Incoming() chan protocol.MessageCommand
}

// Name of the TCPClient is managed by the server.
type TcpChatClient struct {
	conn      net.Conn
	cmdReader *protocol.CommandReader
	cmdWriter *protocol.CommandWriter
	error     chan error
	incoming  chan protocol.MessageCommand
}

func NewTCPClient() *TcpChatClient {
	return &TcpChatClient{
		incoming: make(chan protocol.MessageCommand),
		error:    make(chan error),
	}
}

func (c *TcpChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp4", address)

	if err == nil {
		c.conn = conn
		c.cmdReader = protocol.NewCommandReader(conn)
		c.cmdWriter = protocol.NewCommandWriter(conn)
	}

	return err
}

func (c *TcpChatClient) Start() {
	for {
		cmd, err := c.cmdReader.Read()

		if err != nil {
			c.error <- err
			break // TODO: find a way to recover from this
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

func (c *TcpChatClient) Close() {
	c.conn.Close()
}

func (c *TcpChatClient) Incoming() chan protocol.MessageCommand {
	return c.incoming
}

func (c *TcpChatClient) Error() chan error {
	return c.error
}

func (c *TcpChatClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

func (c *TcpChatClient) SetName(name string) error {
	return c.Send(protocol.NameCommand{
		Name: name,
	})
}

func (c *TcpChatClient) SendMessage(message string) error {
	return c.Send(protocol.SendCommand{
		Message: message,
	})
}
