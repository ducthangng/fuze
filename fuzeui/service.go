package fuzeui

import (
	"fmt"
	"fuze/srv/chat"
	"io"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	IPFocusStat    = false
	ModalFocusStat = false
	InputFocusStat = false
	ListFocusStat  = true

	app        *tview.Application
	feed       *tview.TextView
	grid       *tview.Grid
	inputField *tview.InputField
	IP         *tview.InputField
	action     *tview.TextView

	// IPAddress = []string{}
	Clients []string
)

func newPrimitive(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
}

func AddClient(name string) {
	Clients = append(Clients, name)
}

func CheckClient(name string) bool {
	r := false
	for _, v := range Clients {
		if v == name {
			return true
		}
	}

	if !r {
		AddClient(name)
	}

	return r
}

func SetUI(client *chat.TcpChatClient, isServer bool, IP string) {
	//defer client.Close()
	go client.Start()

	app = tview.NewApplication()
	if isServer {
		action = tview.NewTextView().SetText(fmt.Sprintf("Host: %v", IP)).SetTextAlign(tview.AlignCenter)
	} else {
		action = tview.NewTextView().SetText(fmt.Sprintf("CLI: %v", IP)).SetTextAlign(tview.AlignCenter)
	}

	menu := tview.NewTextView().SetText("People In Chat").SetTextAlign(tview.AlignCenter).SetRegions(true).SetScrollable(true).SetChangedFunc(func() {
		app.Draw()
	})

	feed = tview.NewTextView().SetRegions(true).SetScrollable(true).SetChangedFunc(func() {
		app.Draw()
	})

	inputField = tview.NewInputField().
		SetLabel("Type your name: ").
		SetFieldWidth(100).
		SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				msg := inputField.GetText()

				switch inputField.GetLabel() {
				case "Type your name: ":
					client.SetName(msg)
					inputField.SetLabel("Type chat: ")

				case "Type chat: ":
					client.SendMessage(msg)
				}

				inputField.SetText("")
			}

			if key == tcell.KeyESC {
				app.Stop()
			}
		})

	grid = tview.NewGrid().
		SetRows(3, 0, 5).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Your are Fuzing ..."), 0, 0, 1, 3, 0, 0, false).
		AddItem(feed, 1, 1, 1, 2, 0, 0, false).
		AddItem(newPrimitive("Message Box"), 2, 1, 1, 2, 0, 0, false).
		AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(inputField, 2, 1, 1, 2, 0, 100, true).
		AddItem(action, 2, 0, 1, 1, 0, 100, false)

	go func() {
		for {
			select {
			case msg := <-client.Incoming():
				fmt.Fprintf(feed, "%v: %v \n", msg.Name, msg.Message)
				if !CheckClient(msg.Name) {
					fmt.Fprintf(menu, "\n%v\n", msg.Name)
				}
			case err := <-client.Error():
				if err == io.EOF {
					app.Stop()
				} else {
					panic(err)
				}
			}
		}
	}()

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
