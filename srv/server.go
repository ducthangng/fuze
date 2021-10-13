package srv

import (
	"fmt"
	"fuze/srv/app"
	"fuze/srv/chat"
	"fuze/srv/network"
	"log"
	"net"
	"strconv"
)

var (
	Port                = 3333
	AcceptanceRequestIP = "You got package send from IP: %v"
)

type FuzeServer struct {
	listener net.Listener
	receiver network.FileRetriever
	chatter  chat.ChatServer
}

func NewFuze() *FuzeServer {
	r := network.NewReceiver()
	c := chat.NewChatSrv()

	return &FuzeServer{
		receiver: r,
		chatter:  c,
	}
}

// Run force the server to start listening to incoming packages.
func (s *FuzeServer) Run() error {
	listener, err := net.Listen("tcp4", ":"+strconv.Itoa(Port))
	if err != nil {
		log.Fatalf("Socket listen port %d failed, %s", Port, err)
	}

	s.listener = listener

	// spin := spinner.New(spinner.CharSets[13], 100*time.Millisecond)
	// spin.Start()
	log.Printf("listen ok!")

	return nil
}

func (s *FuzeServer) RetrieveFile() error {
	app := app.NewApp()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatalf("Connection port %d failed, %s", Port, err)
		}

		invi, _ := s.receiver.HasInvitation(conn)
		status := app.AskForAcceptance(fmt.Sprintf(AcceptanceRequestIP, invi.SenderIP))
		if status {
			err := s.receiver.Retrieve(invi.Filename, conn)
			if err != nil {
				log.Fatalf("Error retrieving package: %s", err)
			}
		} else {
			conn.Close()
			return nil
		}
	}
}

func (s *FuzeServer) StartChatServer() {
	app := app.NewApp()

	YourIPs, err := app.RetrieveIP()
	if err != nil {
		log.Fatalf("Exit retrieving IP with error: %v", err)
	}

	// Annouce your RoomIP to welcome connection
	app.Announce(YourIPs)
	s.chatter.Start(s.listener)
}

func (s *FuzeServer) Stop() error {
	return s.listener.Close()
}
