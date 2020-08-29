package twitchchatbot

import (
	"net"

	"gopkg.in/irc.v3"
)

//Im just generating oauth from https://twitchapps.com/tmi/
//Need to add the auto-functionality

//Result is a struct for a channel to return, incase the Bot has an error
type Result struct {
	Source  string //The bot that called the error
	Message string //Message to explain where error occurred
	Error   error  //The error itself
}

//ErrorChannel is a channel you should listen to for any errors that come back across the bots
var ErrorChannel chan Result

//On init, allocate the channel
func init() {
	ErrorChannel = make(chan Result)
}

// BasicBot is a basic bot that connects to twitch
// Set ID, Nickname, Password, Channel, and ReceiveMessage before calling ConnectToChat. 
//Listin on HasJoined to know when the bot has successfully connected to the chat room
type BasicBot struct {
	ID             string //An ID given to the bot to track its progression
	Nickname       string
	Password       string                          //oath:... //This will be replaced at some point with an automatic login function
	Channel        string                          //The username of the chat channel you wish to connect to
	ReceiveMessage func(username, message string)  //A simpler handler, only gives relevant string information. Call bot.SendMessage to reply
	client         *irc.Client
	HasJoined      chan bool //Becomes true when a bot has successfully joined the channel
}

// ConnectToChat connects the bot to the chatroom
func (b *BasicBot) ConnectToChat() {
	b.HasJoined = make(chan bool)
	//This is the unsecured line, follow https://dave.cheney.net/2010/10/05/how-to-dial-remote-ssltls-services-in-go
	//to see if we can get that secure connection working
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	//conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6697")
	pError(err)

	config := irc.ClientConfig{
		Nick:    b.Nickname,
		Pass:    b.Password,
		Handler: irc.HandlerFunc(outsideHandler(b)),
	}

	b.client = irc.NewClient(conn, config)

	go func() {
		err := b.client.Run()
		rs := Result{Error: err, Source: b.ID, Message: "Error running the client"}
		ErrorChannel <- rs
	}()
}

//SendMessage sends a string, and attaches all the headers that you need to send a message to the current channel
//If the bot losses connection to the internet, it never knows, need to fix
func (b *BasicBot) SendMessage(msg string) {
	fmsg := irc.Message{}
	fmsg.Command = "PRIVMSG"
	fmsg.Params = []string{"#" + b.Channel, msg}
	err := b.client.WriteMessage(&fmsg)
	rs := Result{Error: err, Source: b.ID, Message: "Error sending message"}
	ErrorChannel <- rs
}

func outsideHandler(b *BasicBot) func(c *irc.Client, m *irc.Message) {
	return func(c *irc.Client, m *irc.Message) {
		//I think this is good since we won't receiver this until we actually join
		if m.Command == "PRIVMSG" {
			b.ReceiveMessage(m.Name, m.Params[1])
			return
		} else if m.Command == "PING" {
			msg := irc.Message{}
			msg.Command = "PONG"
			msg.Params = append(msg.Params, "tmi.twitch.tv")
			c.WriteMessage(&msg)
			return
		} else if m.Command == irc.RPL_ENDOFMOTD {
			c.Write("JOIN #" + b.Channel)
			//fmt.Println("Joining: ", b.Channel)
			return
		} else if m.Command == irc.RPL_ENDOFNAMES {
			b.HasJoined <- true
		}
	}
}

func pError(err error) {
	if err != nil {
		panic(err)
	}
}
