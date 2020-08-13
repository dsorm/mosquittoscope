package mosquittoscope

import (
	"fmt"
	"log"
	"math"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Display is used to manage the user interface
type Display struct {
	c      chan mqtt.Message
	p      *widgets.Paragraph
	debug  *widgets.Paragraph
	t      *Topic
	height int
	width  int
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
	d.width, d.height = ui.TerminalDimensions()
	background := widgets.NewParagraph()
	background.SetRect(0, 0, d.width, d.height)
	background.Border = false
	// d.background.

	d.debug = widgets.NewParagraph()
	d.debug.Text = "Hello there"
	d.debug.SetRect(d.width-30, 0, d.width, d.height)
	// debug.SetRect(20, 10, 30, 30)
	uiEvents := ui.PollEvents()

	a := widgets.NewParagraph()
	a.SetRect(1, 30, 11, 35)
	a.Text = "a"
	b := widgets.NewParagraph()
	b.SetRect(11, 30, 15, 35)
	b.Text = "b"
	b.BorderLeft = false
	// b.BorderTop = false
	c := widgets.NewParagraph()
	c.SetRect(1, 34, 15, 40)
	c.Text = "c"
	c.BorderLeft = false
	c.BorderTop = false

MainLoop:
	for {
		time.Sleep(1 * time.Millisecond)
		select {
		case msg := <-d.c:
			d.t.UpdateTopics(msg)
		case e := <-uiEvents:
			if e.Type == ui.KeyboardEvent {
				// if e.Payload.(
				break MainLoop
			}
			if e.Type == ui.ResizeEvent {
				payload := e.Payload.(ui.Resize)
				d.width, d.height = payload.Width, payload.Height
				// d.debug.Text = fmt.Sprintf("Height: %d", d.height)
			}
		case <-time.After(20 * time.Millisecond):
			d.updateDisplay()
			ui.Render(c)
			ui.Render(b)
			ui.Render(a)
			// ui.Render(background)
		}
	}
}

func (d *Display) updateDisplay() {
	d.debug.Text = ""
	d.drawTopic(d.t, 1, 1, 0)
	d.debug.SetRect(d.width/2, 0, d.width, d.height)
	ui.Render(d.debug)
}

func (d *Display) drawTopic(t *Topic, x, y, w int) (int, int) {
	height := int(math.Max(float64(t.LeafCount(0)*2), 3))
	d.debug.Text += fmt.Sprintf("%v\n", height)
	t.Box.SetRect(x, y, x+w, y+height)
	t.Box.Text = t.Name
	ui.Render(t.Box)

	totalY := 0
	columnWidth := 0
	for _, st := range t.Subtopics {
		if len(st.Name) > columnWidth {
			columnWidth = len(st.Name)
		}
	}
	for _, st := range t.Subtopics {
		_, subY := d.drawTopic(st, x+w, y+totalY, columnWidth+2)
		totalY += subY
	}

	return 0, totalY + 3

}
