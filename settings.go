package mosquittoscope

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Settings stores configuration values for mosquittoscope as an alternative
// to specifying all the configuration on the command line
type Settings struct {
	Mqtt struct {
		Hostname string
		Port     int
		Username string
		Password string
		ClientID string
	}
}

// NewSettings returns a Settings struct populated with the settings from filename
// or with default values if the load fails.
func NewSettings(filename string) Settings {
	s := Settings{}
	s.Load(filename)
	return s
}

// Load settings from filename into the settings struct
func (s Settings) Load(filename string) {
	s.Mqtt.Hostname = "localhost"
	s.Mqtt.Port = 1443
	s.Mqtt.Username = "username"
	s.Mqtt.Password = "password"
	s.Mqtt.ClientID = "mosquittoscope"

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Failed to open settings file %s. Using defaults.", filename)
		return
	}
	if err := yaml.Unmarshal([]byte(data), &s); err != nil {
		log.Printf("Failed to parse settings file %s. Using defaults.", filename)
		return
	}
}

func (s Settings) String() string {
	d, err := yaml.Marshal(&s)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(d)
}
