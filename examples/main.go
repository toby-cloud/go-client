package main

import (
	"fmt"
	"time"

	bot "toby-cloud/go-client/bot"
	message "toby-cloud/go-client/message"
)

func main() {

	b := bot.NewBot()

	b.SetBotID("")
	b.SetSecret("")

	b.SetOnConnectHandler(func() {
		fmt.Println("Connected to toby!")
		b.Follow("go")
	})

	b.SetOnDisconnectHandler(func() {
		fmt.Println("Disconnected")
	})

	b.SetOnMessageHandler(func(msg message.Message) {
		fmt.Println("Message Received:", msg.Message)
		if msg.AckTag == "" {
			fmt.Println("No ackTag")
		} else {
			fmt.Println("Acknowledgment requested: " + msg.AckTag)
			tags := []string{msg.AckTag}
			m := message.NewMessage("I received your message "+msg.Message+" in golang", "TEXT", tags, "test")
			b.Send(*m)
		}
	})

	b.Start()

}
