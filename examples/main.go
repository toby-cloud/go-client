package main

import (
	"fmt"

	bot "github.com/toby-cloud/toby-go/bot"
	message "github.com/toby-cloud/toby-go/message"
)

func main() {

	b := bot.NewBot()

	b.SetBotID("tiff")
	b.SetSecret("tiff")

	b.SetOnConnectHandler(func() {
		fmt.Println("Connected to toby!")
		m := message.Message{
			From: "tiff",
		}
		b.Send(m)
		// b.Follow("go")
	})

	b.SetOnDisconnectHandler(func() {
		fmt.Println("Disconnected")
	})

	b.SetOnMessageHandler(func(msg message.Message) {
		fmt.Println("Message Received:", msg.Payload)
		/*if msg.AckTag == "" {
			fmt.Println("No ackTag")
		} else {
			fmt.Println("Acknowledgment requested: " + msg.AckTag)
			tags := []string{msg.AckTag}
			m := message.NewMessage("I received your message "+msg.Message+" in golang", "TEXT", tags, "test")
			b.Send(*m)
		}*/
	})

	b.Start()

	b.CreateSocket(false, "ack")

	for {
		/*		if !b.MqttClient.IsConnected() {
				b.Stop()
				return
			}*/
	}

}
