package bot

import (
	"testing"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
)

// TestMQTTMock tests that the mock MQTT works
func TestMQTTMock(t *testing.T) {
	c := newMockClient()

	if tk := c.Connect(); tk == nil {
		t.Fatal("got nil token")
	}

	if tk := c.Subscribe("mock", 0, func(cm mqtt.Client, m mqtt.Message) {
		t.Logf("Received payload %+v", string(m.Payload()))
	}); tk == nil {
		t.Fatal("got nil token")
	}

	if tk := c.Publish("mock", 0, false, []byte(`hello world`)); tk == nil {
		t.Fatal("got nil token")
	}

	if tk := c.Unsubscribe("mock"); tk == nil {
		t.Fatal("got nil token")
	}

	c.Disconnect(0)
}

func TestNewBot(t *testing.T) {
	assert := assert.New(t)

	b := NewBot()

	assert.Equal("", b.ClientID)
	assert.Equal("", b.BotID)
	assert.Equal("", b.Secret)
	assert.Nil(b.Subscriptions)
	assert.Nil(b.OnConnect)
	assert.Nil(b.OnDisconnect)
	assert.Nil(b.OnMessage)
	assert.Nil(b.MqttClient)
}
