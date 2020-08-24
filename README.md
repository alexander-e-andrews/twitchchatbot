# twitchchatbot

## Overview [![GoDoc](https://godoc.org/github.com/alexander-e-andrews/twitchchatbot?status.svg)](https://godoc.org/github.com/alexander-e-andrews/twitchchatbot) [![Go Report Card](https://goreportcard.com/badge/github.com/alexander-e-andrews/twitchchatbot)](https://goreportcard.com/report/github.com/alexander-e-andrews/twitchchatbot)

A basic bot to connect and chat with a twitch channel. Still a work in progress.

## Install

```
go get github.com/alexander-e-andrews/twitchchatbot
```

## Example

```
import(
  basicbot "github.com/alexander-e-andrews/twitchchatbot"
  "gopkg.in/irc.v3"
)

func main(){
   b := basicbot.BasicBot{}
   b.ID = "DummyBot"
   b.Channel = "twitchUsername"
   b.Nickname = "AccountUsername"
   b.Password = "oauth:Code"
   b.Handler = dummyHandler
   b.ConnectToChat()

   msg := "I am a dummy bot"
   b.SendMessage(msg)
   fmt.Println(<-basicbot.ErrorChannel)
}
//Just print out all messages in the chat
func dummyHandler(c *irc.Client, m *irc.Message){
  fmt.Println(m)
}
```

## License

MIT.