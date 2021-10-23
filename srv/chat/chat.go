package chat

import (
	"fuze/srv/chat/protocol"
	"net"
)

type ChatServer interface {
	Start(l net.Listener)
	AddClient(conn net.Conn) *Client
	RemoveClient(client *Client)
	Process(client *Client)
	Broadcast(command interface{}, speaker *Client) error
	Send(name string, command interface{}) error
}

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
