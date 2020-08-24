package twitchchatbot

import (
	"fmt"
	"net"

	"gopkg.in/irc.v3"
)
//Im just generating oauth from https://twitchapps.com/tmi/
//Need to add the auto-functionality

//Result is a struct for a channel to return, incase the Bot has an error
type Result struct{
	Source string //The bot that called the error
	Message string //Message to explain where error occurred
	Error error //The error itself
}
//ErrorChannel is a channel you should listen to for any errors that come back across the bots
var ErrorChannel chan Result

//On init, allocate the channel
func init(){
	ErrorChannel = make(chan Result)
}

// BasicBot is a basic bot taht connects to twitch
type BasicBot struct {
	ID string //An ID given to the bot to track its progression
	Nickname string
	Password string //oath:... //This will be replaced at some point with an automatic login function
	Handler  func(*irc.Client, *irc.Message) //Your functionality that you want to handle
	Client   *irc.Client
	Channel  string //The username of the chat channel you wish to connect to
}

// ConnectToChat connects the bot to the chatroom
func (b *BasicBot) ConnectToChat() {
	//This is the unsecured line, follow https://dave.cheney.net/2010/10/05/how-to-dial-remote-ssltls-services-in-go
	//to see if we can get that secure connection working
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	//conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6697")
	pError(err)

	config := irc.ClientConfig{
		Nick:    b.Nickname,
		Pass:    b.Password,
		Handler: irc.HandlerFunc(outsideHandler(b.Handler, b.Channel)),
	}

	b.Client = irc.NewClient(conn, config)

	go func(){
		err := b.Client.Run()
		rs := Result{Error: err, Source: b.ID, Message: "Error running the client"}
		ErrorChannel <- rs
	}()
}

//SendMessage sends a string, and attaches all the headers that you need to send a message to the current channel
//If the bot losses connection to the internet, it never knows, need to fix
func (b *BasicBot) SendMessage(msg string) {
	fmsg := irc.Message{}
	fmsg.Command = "PRIVMSG"
	fmsg.Params = []string{"#"+b.Channel, msg}
	err := b.Client.WriteMessage(&fmsg)
	rs := Result{Error: err, Source: b.ID, Message: "Error sending message"}
	ErrorChannel <- rs
}

func outsideHandler(h func(*irc.Client, *irc.Message), channel string) func(c *irc.Client, m *irc.Message) {
	return func(c *irc.Client, m *irc.Message) {
		if m.Command == "PING" {
			msg := irc.Message{}
			msg.Command = "PONG"
			msg.Params = append(msg.Params, "tmi.twitch.tv")
			c.WriteMessage(&msg)
			return
		}else if m.Command == irc.RPL_ENDOFMOTD {
			c.Write("JOIN #"+ channel)
			fmt.Println("Joining: ", channel)
			return
		}else{
			h(c, m)
		}
	}
}

func pError(err error) {
	if err != nil {
		panic(err)
	}
}
