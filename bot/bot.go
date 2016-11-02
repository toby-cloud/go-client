package bot

import (
	"encoding/json"
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	message "github.com/toby-cloud/toby-go/message"
)

// MessageHandler
type OnMessageHandler func(message.Message)

// ConnectionLostHandler
type OnDisconnectHandler func()

// OnConnectHandler
type OnConnectHandler func()

type Bot struct {
	ClientID      string
	BotID         string
	Secret        string
	Subscriptions []string
	OnConnect     OnConnectHandler
	OnDisconnect  OnDisconnectHandler
	OnMessage     OnMessageHandler
	MqttClient    MQTT.Client
}

func NewBot() *Bot {
	o := &Bot{
		ClientID:      "",
		BotID:         "",
		Secret:        "",
		Subscriptions: nil,
		OnConnect:     nil,
		OnDisconnect:  nil,
		OnMessage:     nil,
		MqttClient:    nil,
	}
	return o
}

func (b *Bot) SetBotID(id string) *Bot {
	b.BotID = id
	b.ClientID = id
	return b
}

func (b *Bot) SetSecret(s string) *Bot {
	b.Secret = s
	return b
}

func (b *Bot) SetOnMessageHandler(handler OnMessageHandler) *Bot {
	b.OnMessage = handler
	return b
}

func (b *Bot) SetOnConnectHandler(handler OnConnectHandler) *Bot {
	b.OnConnect = handler
	return b
}

func (b *Bot) SetOnDisconnectHandler(handler OnDisconnectHandler) *Bot {
	b.OnDisconnect = handler
	return b
}

// private functions
func (b *Bot) onMessage(client MQTT.Client, msg MQTT.Message) {
	// TODO preprocess message into Toby Message
	var m message.Message
	if err := json.Unmarshal(msg.Payload(), &m); err != nil {
		panic(err)
	}
	b.OnMessage(m) // call user onMessage
}
func (b *Bot) onConnect() {
	b.OnConnect()
}
func (b *Bot) onDisconnect() {
	b.OnDisconnect()
}

// Set bot subscriptions to tag
func (b *Bot) Follow(tag string) {
	token := b.MqttClient.Publish("server/"+b.BotID+"/follow", byte(0), false, "{\"tags\": [\""+tag+"\"]}")
	token.Wait()
	fmt.Println("Followed #" + tag)
}

// Send bot message with one tag
func (b *Bot) Send(m message.Message) {
	token := b.MqttClient.Publish("server/"+b.BotID+"/send", byte(0), false, m.String())
	token.Wait()
	fmt.Println("Sent message")
}

// Start Toby MQTT connection
// If connection succesful, subscribe to bot data
func (b *Bot) Start() {

	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://toby.cloud:444")
	opts.SetClientID(b.BotID)
	opts.SetUsername(b.BotID)
	opts.SetPassword(b.Secret)
	opts.SetCleanSession(false)
	opts.SetKeepAlive(30 * time.Second)

	// CONNECT TO BROKER
	b.MqttClient = MQTT.NewClient(opts)
	if token := b.MqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	b.onConnect()

	// SUBSCRIBE TO BOT DATA
	if token := b.MqttClient.Subscribe("client/"+b.BotID+"/#", byte(0), b.onMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

}
