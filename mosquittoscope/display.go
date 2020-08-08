package mosquittoscope

import (
	"container/list"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Display is used to manage the user interface
type Display struct {
	c chan mqtt.Message
	p *widgets.Paragraph
	t *Topic
}

// Topic defines a tree used to hierarchically store messages received over MQTT
type Topic struct {
	Name      string
	Messages  *list.List
	Subtopics []*Topic
	Box       *widgets.Paragraph
}

// NewDisplay builds a new display interface thingo
func NewDisplay(s *Settings) *Display {
	d := Display{}
	d.t = NewTopic("")
	return &d
}

// SetTopicChannel is used to register the channel delivering the topics to the display code
func (d *Display) SetTopicChannel(c chan mqtt.Message) {
	d.c = c
}

// DisplayLoop blocks, updating the display and handling user input
func (d *Display) DisplayLoop(done chan bool) {
	defer func(d chan bool) { d <- true }(done)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	d.p = widgets.NewParagraph()
	d.p.Text = "Hello World!"
	d.p.SetRect(0, 0, 25, 5)

MainLoop:
	for {
		time.Sleep(1 * time.Millisecond)
		select {
		case msg := <-d.c:
			d.p.Text = fmt.Sprintf("TOPIC: %s\n", msg.Topic())
			d.t.UpdateTopics(msg)
		case e := <-ui.PollEvents():
			if e.Type == ui.KeyboardEvent {
				break MainLoop
			}
		}
		ui.Render(d.p)
	}
}
