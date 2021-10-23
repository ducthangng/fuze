package main

import (
	"fuze/fuzeui"
	"fuze/srv"
	"fuze/srv/app"
	"fuze/srv/chat"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.Println("Welcome to fuze, homie!")

	app := app.NewApp()
	service := app.AskForService()

	switch service {
	case 1:
		RunClient(app, false, "")
	case 2:
		RunServer(app)
	case 3:
		os.Exit(1)
	}
}

func RunClient(app *app.App, isServer bool, IP string) {
	address := IP
	client := chat.NewTCPClient()
	for {
		if len(address) == 0 {
			address = app.AskForString("Room IP: ")
			if net.ParseIP(address) == nil {
				app.Announce([]string{"Invalid, type again."})
				continue
			}
		}

		addr := strings.Join([]string{address, strconv.Itoa(3333)}, ":")
		err := client.Dial(addr)
		if err != nil {
			app.Announce([]string{"Invalid, try again."})
		} else {
			break
		}

	}

	fuzeui.SetUI(client, isServer, address)
}

func RunServer(app *app.App) {
	server := srv.NewFuze()
	server.Run()

	IP, err := app.RetrieveIP()
	if err != nil {
		log.Fatalf("exit with error: %v", err)
	}

	app.Announce(IP)
	go server.StartChatServer()

	RunClient(app, true, IP[0])
}
