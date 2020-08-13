package mosquittoscope

import (
	"container/list"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gizak/termui/v3/widgets"
)

const messageBufferLength = 10

// Topic defines a tree used to hierarchically store messages received over MQTT
type Topic struct {
	Name      string
	Messages  *list.List
	Subtopics []*Topic
	Box       *widgets.Paragraph
	Height    int
	Width     int
}

// NewTopic returns a Topic struct, ready to be populated via UpdateTopics.
func NewTopic(name string) *Topic {
	t := Topic{}
	t.InitTopic(name)
	t.Box = widgets.NewParagraph()
	t.Box.BorderLeft = false
	// t.Box.BorderTop = false
	return &t
}

// InitTopic initialises a topic
func (t *Topic) InitTopic(name string) {
	t.Name = name
	t.Messages = list.New()
}

// UpdateTopics receives an MQTT message and populates the hierachy accordingly
func (t *Topic) UpdateTopics(m mqtt.Message) {
	if t := getTopic(t, m.Topic()); t != nil {
		t.addMessage(string(m.Payload()))
	}
}

func (t *Topic) addMessage(m string) {
	t.Messages.PushBack(m)
	if t.Messages.Len() > messageBufferLength {
		t.Messages.Remove(t.Messages.Front())
	}
	// for e := t.Messages.Front(); e != nil; e = e.Next() {
	// 	fmt.Printf("%q\n", e.Value)
	// }
}

// getTopic returns the topic object from t matching the topic path defined by s
func getTopic(t *Topic, s string) *Topic {
	split := strings.SplitN(s, "/", 2)
	topTopic := split[0]
	atBottom := len(split) == 1

	if t.Name == topTopic && atBottom {
		return t
	}

	for _, i := range t.Subtopics {
		if i.Name == topTopic {
			if atBottom {
				return i
			}
			return getTopic(i, split[1])
		}
	}

	new := NewTopic(topTopic)
	t.Subtopics = append(t.Subtopics, new)
	if atBottom {
		return new
	}
	return getTopic(new, split[1])
}

// LeafCount returns the number of leaves this topic has
func (t *Topic) LeafCount(n int) int {
	subLeafCount := 0
	// If this topic has a value, it is considered a leaf
	// even if it also has subtopics.
	// if t.Messages.Len() > 0 && len(t.Subtopics) == 0{
	if t.Messages.Len() > 0 {
		subLeafCount++
	}
	for _, i := range t.Subtopics {
		subLeafCount += i.LeafCount(0)
	}
	return n + subLeafCount
}
