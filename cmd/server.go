package cmd

import (
	"fmt"
	"fuze/srv"
	"fuze/srv/app"
	"fuze/srv/chat"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	GetAddress     = "What is the service address: "
	InvalidAddress = "Invalid address"
	GetUsername    = "What is your username: "
	ConnectionLost = "Connection closed connection from server."
)

// main executed command.
var rootCmd = &cobra.Command{
	Use:   "fuze",
	Short: "minimal LAN tool, support file sharing and chatting",
	Run: func(cmd *cobra.Command, args []string) {
		app := app.NewApp()

		// Declare service
		fuze := srv.NewFuze()

		// Prompt user for service
		service := app.AskForService()
		switch service {
		case 1:
			StartChatClient()
		case 2:
			StartChatServer(fuze)
		case 3:
			StartFileRetrieve(fuze)
		default:
			fmt.Println("not available")
		}

	},
}

func StartFileRetrieve(fuze *srv.FuzeServer) {
	fuze.RetrieveFile()
}

func StartChatClient() {
	app := app.NewApp()
	var address string

	for {
		address = app.AskForString(GetAddress)
		if net.ParseIP(address) == nil {
			app.Announce([]string{InvalidAddress})
		} else {
			break
		}
	}

	client := chat.NewTCPClient()

	addr := strings.Join([]string{address, strconv.Itoa(3333)}, ":")
	err := client.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	go client.Start()

	name := app.AskForString(GetUsername)
	client.SetName(name)

	// Listen to keyboard
	go func() {
		for {
			message := app.WaitForString()
			client.SendMessage(message)
		}
	}()

	// Listen server broadcast
	func() {
		for {
			select {
			case err := <-client.Error():
				if err == io.EOF {
					app.Announce([]string{ConnectionLost})
				} else {
					panic(err)
				}

			case msg := <-client.Incoming():
				app.Announce([]string{fmt.Sprintf("%v: %v", msg.Name, msg.Message)})
			}
		}
	}()
}

func StartChatServer(fuze *srv.FuzeServer) {
	app := app.NewApp()
	fuze.Run()

	// Listen to keyboard
	go func() {
		for {
			message := app.WaitForString()
			if message == "exit" {
				fuze.Stop()
			}
		}
	}()

	fuze.StartChatServer()
}
