package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	message "github.com/toby-cloud/toby-go/message"
)

// MessageHandler
type OnMessageHandler func(string, message.Message)

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

// NewBot constructs a new Bot
func NewBot() *Bot {
	return &Bot{
		ClientID:      "",
		BotID:         "",
		Secret:        "",
		Subscriptions: nil,
		OnConnect:     nil,
		OnDisconnect:  nil,
		OnMessage:     nil,
		MqttClient:    nil,
	}
}

// Start starts the Toby MQTT connection
func (b *Bot) Start() error {
	// if connection successful, subscribe to bot data

	// set MQTT options
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://toby.cloud:444")
	opts.SetClientID(b.BotID)
	opts.SetUsername(b.BotID)
	opts.SetPassword(b.Secret)
	opts.SetCleanSession(false)
	opts.SetKeepAlive(30 * time.Second)

	// new MQTT client
	b.MqttClient = MQTT.NewClient(opts)

	// connect to mqtt broker
	if token := b.MqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	// execute onConnect callback when connected successfully
	b.OnConnect()

	// subscribe to bot data and call onMessage when messages are received
	if token := b.MqttClient.Subscribe("client/"+b.BotID+"/#", byte(0), b.onMessage); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// End ends the Toby MQTT connection
func (b *Bot) End() error {
	// TODO: throw error if OnDisconnect() throws an error
	b.OnDisconnect()

	b.MqttClient.Unsubscribe("#")
	b.MqttClient.Disconnect(250)

	return nil
}

// Send sends a message with one tag
func (b *Bot) Send(m message.Message) error {
	msg, err := m.String()
	if err != nil {
		return err
	}

	token := b.MqttClient.Publish("server/"+b.BotID+"/send", byte(0), false, msg)
	token.Wait()
	fmt.Println("Sent message: " + msg)

	return nil
}

// HooksOn turns on webhooks
func (b *Bot) HooksOn(hookSecret, ackTag string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// json the hookSecret and ackTag
	m, err := json.Marshal(map[string]string{"hookSecret": hookSecret, "ackTag": ackTag})
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/hooks-on", byte(0), false, string(m))
	token.Wait()

	return nil
}

// HooksOff turns off webhooks
func (b *Bot) HooksOff(ackTag string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// json the ackTag
	m, err := json.Marshal(map[string]string{"ackTag": ackTag})
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/hooks-off", byte(0), false, string(m))
	token.Wait()

	return nil
}

// Info gets the bot info from MQTT
func (b *Bot) Info(ackTag string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// json the ackTag
	m, err := json.Marshal(map[string]string{"ackTag": ackTag})
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/info", byte(0), false, string(m))
	token.Wait()

	return nil
}

// CreateBot creates the bot in MQTT
func (b *Bot) CreateBot(name, password, ackTag string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// construct payload
	payload := map[string]string{
		"id":     name,
		"secret": password,
		"ackTag": ackTag,
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/create-bot", byte(0), false, string(m))
	token.Wait()

	return nil
}

// CreateSocket creates the socket
func (b *Bot) CreateSocket(persist bool, ackTag string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// construct payload
	payload := map[string]interface{}{
		"ackTag":  ackTag,
		"persist": persist,
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/create-socket", byte(0), false, string(m))
	token.Wait()

	return nil
}

// RemoveBot removes the bot
func (b *Bot) RemoveBot(targetId, ackTag string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// construct payload
	payload := map[string]string{
		"ackTag": ackTag,
		"botId":  targetId,
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/remove-bot", byte(0), false, string(m))
	token.Wait()

	return nil
}

// RemoveSocket removes the socket
func (b *Bot) RemoveSocket(targetId, ackTag string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// construct payload
	payload := map[string]string{
		"ackTag": ackTag,
		"botId":  targetId,
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/remove-socket", byte(0), false, string(m))
	token.Wait()

	return nil
}

// Follow adds bot subscription to tag
func (b *Bot) Follow(tag, ackTag string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// construct payload
	payload := map[string]string{
		"tags":   tag,
		"ackTag": ackTag,
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/follow", byte(0), false, string(m))
	token.Wait()
	fmt.Println("Followed #" + tag)

	return nil
}

// Unfollow removes bot subscription to tag
func (b *Bot) Unfollow(tag, ackTag string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// construct payload
	payload := map[string]string{
		"tags":   tag,
		"ackTag": ackTag,
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/unfollow", byte(0), false, string(m))
	token.Wait()
	fmt.Println("Followed #" + tag)

	return nil
}

// SetBotID sets the BotID and ClientID
func (b *Bot) SetBotID(id string) {
	// TODO: validity check
	b.BotID = id
	b.ClientID = id
}

func (b *Bot) SetSecret(s string) {
	// TODO: validity check
	b.Secret = s
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

// onMessage processes the received MQTT message then calls OnMessage handler
func (b *Bot) onMessage(client MQTT.Client, msg MQTT.Message) {
	// preprocess message topics
	topics := strings.Split(msg.Topic(), "/")
	topics = topics[2:]
	t := strings.Join(topics, "/")

	// TODO: some error handling in here
	m := message.Message{}
	if err := json.Unmarshal(msg.Payload(), &m); err != nil {
		t = ""
	}

	b.OnMessage(t, m)
}
