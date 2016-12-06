package bot

import (
	"encoding/json"
	"errors"
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

type CredsResponse struct {
	Id   string   `json:"id"`
	Sk   string   `json:"sk"`
	Tags []string `json:"tags"`
	From string   `json:"from"`
	Ack  string   `json:"ack"`
}

// NewBot constructs a new Bot.
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

// Start starts the Toby MQTT connection.
func (b *Bot) Start() {
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
		panic(token.Error())
	}

	// subscribe to bot data and call onMessage when messages are received
	if token := b.MqttClient.Subscribe("client/"+b.BotID, byte(0), b.onMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// execute onConnect callback when connected successfully
	b.OnConnect()

}

// Stop ends the Toby MQTT connection.
func (b *Bot) Stop() {
	b.OnDisconnect()
	b.MqttClient.Disconnect(250)
}

// Send sends a message with one tag.
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

// HooksOn turns on webhooks.
func (b *Bot) HooksOn(sk, ack string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// json the hookSecret and ackTag
	m, err := json.Marshal(map[string]string{"sk": sk, "ack": ack})
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/hooks-on", byte(0), false, string(m))
	token.Wait()

	return nil
}

// HooksOff turns off webhooks.
func (b *Bot) HooksOff(ack string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// json the ackTag
	m, err := json.Marshal(map[string]string{"ack": ack})
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/hooks-off", byte(0), false, string(m))
	token.Wait()

	return nil
}

// FIXME: handle info requests not using map[string]string
// Info gets the bot info from MQTT.
func (b *Bot) Info(ack string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// json the ackTag
	m, err := json.Marshal(map[string]string{"ack": ack})
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/info", byte(0), false, string(m))
	token.Wait()

	return nil
}

// CreateBot creates the bot in MQTT.
func (b *Bot) CreateBot(name, password, ack string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// construct payload
	payload := map[string]string{
		"id":  name,
		"sk":  password,
		"ack": ack,
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

// CreateSocket creates the socket.
func (b *Bot) CreateSocket(persist bool, ack string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// construct payload
	payload := map[string]interface{}{
		"ack":     ack,
		"persist": persist,
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	fmt.Println("creating socket", string(m))
	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/create-socket", byte(0), false, string(m))
	token.Wait()

	return nil
}

// RemoveBot removes the bot.
func (b *Bot) RemoveBot(targetId, ack string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// construct payload
	payload := map[string]string{
		"ack": ack,
		"id":  targetId,
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

// RemoveSocket removes the socket.
func (b *Bot) RemoveSocket(targetId, ack string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// construct payload
	payload := map[string]string{
		"ack": ack,
		"id":  targetId,
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

// Follow adds bot subscription to tag.
func (b *Bot) Follow(tags, ack string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// construct payload
	payload := map[string]string{
		"tags": tags,
		"ack":  ack,
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/follow", byte(0), false, string(m))
	token.Wait()

	return nil
}

// Unfollow removes bot subscription to tag.
func (b *Bot) Unfollow(tags, ack string) error {
	// return error if MqttClient is not connected
	if !b.MqttClient.IsConnected() {
		return errors.New("MQTTClient not connected")
	}

	// construct payload
	payload := map[string]string{
		"tags": tags,
		"ack":  ack,
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// publish to MqttClient
	token := b.MqttClient.Publish("server/"+b.BotID+"/unfollow", byte(0), false, string(m))
	token.Wait()

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

func (b *Bot) SetOnMessageHandler(handler OnMessageHandler) {
	b.OnMessage = handler
}

func (b *Bot) SetOnConnectHandler(handler OnConnectHandler) {
	b.OnConnect = handler
}

func (b *Bot) SetOnDisconnectHandler(handler OnDisconnectHandler) {
	b.OnDisconnect = handler
}

// onMessage processes the received MQTT message then calls OnMessage handler
func (b *Bot) onMessage(client MQTT.Client, msg MQTT.Message) {

	m := message.Message{}
	if err := json.Unmarshal(msg.Payload(), &m); err != nil {
		// try unmarshaling as creds response
		r := CredsResponse{}
		if err = json.Unmarshal(msg.Payload(), &r); err != nil {
			b.Stop()
			fmt.Println("Error unmarshaling message - disconnecting...")
			return
		}
		m.Id = r.Id
		m.Sk = r.Sk
		m.Tags = r.Tags
		m.Ack = r.Ack
		m.From = r.From
		// TODO: handle error unmarshaling
	}

	b.OnMessage(m)
}
