package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"./mosquittoscope"
)

var mqttUsername = flag.String("u", "", "MQTT broker username")
var mqttPassword = flag.String("P", "", "MQTT broker password")
var mqttHostname = flag.String("h", "", "MQTT broker hostname")
var mqttPort = flag.String("p", "", "MQTT broker port")
var settingsFile = flag.String("s", "default.yaml", "Settings file path")

func main() {
	flag.Parse()
	fmt.Println(flag.Args())
	s := mosquittoscope.NewSettings(*settingsFile)
	if len(*mqttUsername) > 0 {
		s.MQTT.Username = *mqttUsername
	}
	if len(*mqttPassword) > 0 {
		s.MQTT.Password = *mqttPassword
	}
	if len(*mqttHostname) > 0 {
		s.MQTT.Hostname = *mqttHostname
	}
	if len(*mqttPort) > 0 {
		port, err := strconv.Atoi(*mqttPort)
		if err != nil {
			log.Printf("Failed to read port from command line arguments. Defaulting to port %d.", s.MQTT.Port)
		}
		s.MQTT.Port = port
	}
	m := mosquittoscope.NewMQTTMonitor(s)
	d := mosquittoscope.NewDisplay(s)
	d.SetTopicChannel(m.GetTopicChannel())

	if err := m.Subscribe("#"); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	fmt.Println("Cool beans")
	var done chan bool
	// Arguably this doesn't have to be a goroutine, but this is a learning exercise.
	go d.DisplayLoop(done)
	<-done
}
