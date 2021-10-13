package app

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) AskForService() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Choose the service by enter the number accordingly: ")
	fmt.Println("◙ Join Room Chat:   1")
	fmt.Println("◙ Create Room Chat: 2")
	fmt.Println("◙ Share Documents:  3")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')

		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, " ", "", -1)

		code, err := strconv.Atoi(text)
		if err != nil {
			log.Println("We do not recognize the service")
		} else {
			return code
		}
	}
}

func (a *App) AskForAcceptance(request string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(request)

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, " ", "", -1)
		text = strings.ToLower(text)

		switch text {
		case "y":
			return true
		case "n":
			return false
		default:
			return false
		}
	}
}

func (a *App) AskForString(request string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(request)

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		return text
	}
}

func (a *App) Announce(messages []string) {
	for _, v := range messages {
		fmt.Println(v)
	}

}

func (a *App) WaitForString() string {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		return text
	}
}

// RetrieveIP retrieves the current IP Address of the user.
func (a *App) RetrieveIP() (IPs []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return IPs, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				IPs = append(IPs, ipnet.IP.String())
			}
		}
	}

	return IPs, err
}
