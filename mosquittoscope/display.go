package mosquittoscope

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Display is used to manage the user interface
type Display struct {
	c chan mqtt.Message
}

// NewDisplay builds a new display interface thingo
func NewDisplay(s *Settings) *Display {
	d := Display{}

	return &d
}

// SetTopicChannel is used to register the channel delivering the topics to the display code
func (d *Display) SetTopicChannel(c chan mqtt.Message) {
	d.c = c
}

// DisplayLoop blocks forever, updating the display and handling user input
func (d *Display) DisplayLoop(done chan bool) {
	for {
		time.Sleep(1 * time.Millisecond)
		if d.c == nil {
			continue
		}
		msg := <-d.c
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}

	done <- true
}

// func main() {
// 	if err := ui.Init(); err != nil {
// 		log.Fatalf("failed to initialize termui: %v", err)
// 	}
// 	defer ui.Close()

// 	p := widgets.NewParagraph()
// 	p.Text = "Hello World!"
// 	p.SetRect(0, 0, 25, 5)

// 	ui.Render(p)

// 	for e := range ui.PollEvents() {
// 		if e.Type == ui.KeyboardEvent {
// 			break
// 		}
// 	}
// }
