package gui

import (
	"fyne.io/fyne/v2/canvas"
	"github.com/Git5737/lexchanger/client/internal/chat"
	pb "github.com/Git5737/lexchanger/proto/chat/proto"
	"image/color"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func chatBubble(text, sender string, isMe bool) *fyne.Container {
	const bubbleWidth = 300

	senderLabel := widget.NewLabelWithStyle(sender, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	msg := widget.NewLabel(text)
	msg.Wrapping = fyne.TextWrapWord
	msg.Alignment = fyne.TextAlignLeading

	content := container.NewVBox(senderLabel, msg)

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(bubbleWidth, 0))

	limited := container.NewMax(spacer, container.NewPadded(content))

	bgColor := color.NRGBA{R: 50, G: 50, B: 50, A: 255}
	if isMe {
		bgColor = color.NRGBA{R: 0, G: 122, B: 255, A: 255}
	}
	bg := canvas.NewRectangle(bgColor)

	bubble := container.NewMax(bg, limited)

	// Вирівнювання
	if isMe {
		return container.NewVBox(container.NewHBox(layout.NewSpacer(), bubble))
	}
	return container.NewVBox(container.NewHBox(bubble, layout.NewSpacer()))
}

func Start() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("lexchanger")

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter your name")

	chatList := container.NewVBox()
	chatScroll := container.NewVScroll(chatList)
	chatScroll.SetMinSize(fyne.NewSize(700, 500))

	msgEntry := widget.NewEntry()
	msgEntry.SetPlaceHolder("Type message...")

	sendBtn := widget.NewButton("Send", nil)
	bottom := container.NewBorder(nil, nil, nil, sendBtn, msgEntry)
	top := container.NewVBox(nameEntry)

	w.SetContent(container.NewBorder(top, bottom, nil, nil, chatScroll))
	w.Resize(fyne.NewSize(700, 700))
	w.Show()

	go func() {
		client, err := chat.NewClient("localhost:50051")
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		myName := ""

		sendBtn.OnTapped = func() {
			if myName == "" {
				myName = strings.TrimSpace(nameEntry.Text)
				client.Send(&pb.Events{
					Event: &pb.Events_ClientLogin{
						ClientLogin: &pb.Events_Login{
							Name: myName,
						},
					},
				})
			}

			msg := strings.TrimSpace(msgEntry.Text)
			if msg == "" || myName == "" {
				return
			}
			client.Send(&pb.Events{
				Event: &pb.Events_ClientMessage{
					ClientMessage: &pb.Events_Message{
						Name:    myName,
						Message: msg,
					},
				},
			})
			msgEntry.SetText("")
		}

		// Login after delay if name filled
		time.Sleep(time.Second)
		if name := strings.TrimSpace(nameEntry.Text); name != "" {
			myName = name
			client.Send(&pb.Events{
				Event: &pb.Events_ClientLogin{
					ClientLogin: &pb.Events_Login{
						Name: name,
					},
				},
			})
		}

		for {
			event, err := client.Recv()
			if err != nil {
				log.Println("recv error:", err)
				return
			}

			switch ev := event.Event.(type) {
			case *pb.Events_ClientMessage:
				isMe := ev.ClientMessage.Name == myName
				bubble := chatBubble(ev.ClientMessage.Message, ev.ClientMessage.Name, isMe)
				chatList.Add(bubble)
				chatScroll.ScrollToBottom()
			}
		}
	}()
	a.Run()
}
