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

// MqttMonitor monitors mqtt... ?
type MqttMonitor struct {
	c mqtt.Client
	s Settings
}

func (m MqttMonitor) mqttmonitor() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	fullPath := fmt.Sprintf("tcp://%s:%d", m.s.Mqtt.Hostname, m.s.Mqtt.Port)
	opts := mqtt.NewClientOptions().AddBroker(fullPath).SetClientID(m.s.Mqtt.ClientID)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetUsername(m.s.Mqtt.Username)
	opts.SetPassword(m.s.Mqtt.Password)

	m.c = mqtt.NewClient(opts)
	if token := m.c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := m.c.Subscribe("#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	m.c.Disconnect(250)

	time.Sleep(1 * time.Second)
}
