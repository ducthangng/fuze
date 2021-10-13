package chat

import "net"

type ChatServer interface {
	Start(l net.Listener)
	AddClient(conn net.Conn) *Client
	RemoveClient(client *Client)
	Process(client *Client)
	Broadcast(command interface{}, speaker *Client) error
	Send(name string, command interface{}) error
}
