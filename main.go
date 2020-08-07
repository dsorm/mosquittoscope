package main

import (
	"flag"
	"fmt"

	"./mosquittoscope"
)

var mqttUsername = flag.String("u", "", "MQTT broker username")
var mqttPassword = flag.String("P", "", "MQTT broker password")
var mqttHostname = flag.String("h", "", "MQTT broker hostname")
var mqttPort = flag.String("p", "", "MQTT broker port")

func main() {
	flag.Parse()
	fmt.Println(flag.Args())
	s := mosquittoscope.NewSettings("boop.yaml")
	if len(*mqttUsername) > 0 {
		s.MQTT.Username = *mqttUsername
	}
	m := mosquittoscope.NewMQTTMonitor(s)

	fmt.Printf("%q\n", *mqttUsername)
	fmt.Printf("%q\n", m)
}
