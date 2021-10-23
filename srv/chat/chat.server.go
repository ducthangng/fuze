package chat

import (
	"errors"
	"fuze/srv/chat/protocol"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	ErrUnknownClient = errors.New("unknown client")
)

type ChatSrv struct {
	Clients []*Client
	mu      *sync.Mutex
}

type Client struct {
	Name         string
	AssignedConn net.Conn
	writer       *protocol.CommandWriter
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		writer:       protocol.NewCommandWriter(conn),
		AssignedConn: conn,
	}
}

func NewChatSrv() *ChatSrv {
	return &ChatSrv{
		mu: &sync.Mutex{},
	}
}

// func (chat *ChatSrv) RunningConnection() {
// 	log.Printf("Connected conn: %v", len(chat.Clients))
// }

func (chat *ChatSrv) Start(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "closed network connection") {
				return
			}

			log.Fatalf("Connection failed, %s", err)
		} else {
			client := chat.AddClient(conn)
			go chat.Process(client)
		}
	}
}

func (chat *ChatSrv) AddClient(conn net.Conn) *Client {
	chat.mu.Lock()
	defer chat.mu.Unlock()

	client := NewClient(conn)
	chat.Clients = append(chat.Clients, client)

	return client
}

func (chat *ChatSrv) RemoveClient(client *Client) {
	chat.mu.Lock()
	defer chat.mu.Unlock()

	for i, v := range chat.Clients {
		if v == client {
			chat.Clients = append(chat.Clients[:i], chat.Clients[i+1:]...)
		}
	}

	client.AssignedConn.Close()
}

func (chat *ChatSrv) Process(client *Client) {
	cmdReader := protocol.NewCommandReader(client.AssignedConn)
	defer chat.RemoveClient(client)

	for {
		cmd, err := cmdReader.Read()

		if err != nil && err != io.EOF {
			log.Printf("Read error: %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.SendCommand:
				go chat.Broadcast(protocol.MessageCommand{
					Message: v.Message,
					Name:    client.Name,
				}, client)

			case protocol.NameCommand:
				client.Name = v.Name
				// log.Printf("Client: %v just joined chat\n", client.Name)
			}
		}

		if err == io.EOF {
			break
		}
	}
}

func (chat *ChatSrv) Broadcast(command interface{}, speaker *Client) error {
	for _, client := range chat.Clients {
		client.writer.Write(command)
	}

	return nil
}

func (chat *ChatSrv) Send(name string, command interface{}) error {
	for _, client := range chat.Clients {
		if client.Name == name {
			return client.writer.Write(command)
		}
	}

	return ErrUnknownClient
}
