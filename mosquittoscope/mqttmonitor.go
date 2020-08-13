package mosquittoscope

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

// MQTTMonitor monitors mqtt... ?
type MQTTMonitor struct {
	c            mqtt.Client
	s            *Settings
	topicChannel chan mqtt.Message
}

// NewMQTTMonitor returns a pointer to an instance of MQTTMonitor
func NewMQTTMonitor(s *Settings) *MQTTMonitor {
	m := MQTTMonitor{}
	m.s = s
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stderr, "", 0)
	fullPath := fmt.Sprintf("tcp://%s:%d", m.s.MQTT.Hostname, m.s.MQTT.Port)
	opts := mqtt.NewClientOptions().AddBroker(fullPath).SetClientID(m.s.MQTT.ClientID)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetUsername(m.s.MQTT.Username)
	opts.SetPassword(m.s.MQTT.Password)

	m.c = mqtt.NewClient(opts)
	if token := m.c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connected to MQTT broker")
	return &m
}

func (m *MQTTMonitor) publishCallback(client mqtt.Client, msg mqtt.Message) {
	if m.topicChannel == nil {
		fmt.Println("Got message, but no channel, so dropping on floor.")
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		return
	}
	// This feels like witchcraft
	m.topicChannel <- msg
}

// Subscribe subscribes to the provided topic.
func (m *MQTTMonitor) Subscribe(topic string) (err error) {
	if token := m.c.Subscribe("t/#", 0, m.publishCallback); token.Wait() && token.Error() != nil {
		return fmt.Errorf("Failed to subscribe to %q", topic)
	}
	return nil
}

// GetTopicChannel returns a channel from which received MQTT messages can be... got?
// I don't know the proper terminology
func (m *MQTTMonitor) GetTopicChannel() chan mqtt.Message {
	m.topicChannel = make(chan mqtt.Message)
	return m.topicChannel
}

// SubscribeAndGetChannel will subscribe to the given topic and return a channel through which
// publications to that topic will be fed.
func (m *MQTTMonitor) SubscribeAndGetChannel(topic string) (chan mqtt.Message, error) {
	channel := make(chan mqtt.Message)
	callback := func(client mqtt.Client, msg mqtt.Message) {
		channel <- msg
	}
	if token := m.c.Subscribe("#", 0, callback); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("Failed to subscribe to %q", topic)
	}
	return channel, nil
}
