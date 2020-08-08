package mosquittoscope

import (
	"container/list"
	"strings"

	"github.com/gizak/termui/v3/widgets"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const messageBufferLength = 10

// NewTopic returns a Topic struct, ready to be populated via UpdateTopics.
func NewTopic(name string) *Topic {
	t := Topic{}
	t.Name = name
	t.Messages = list.New()
	t.Box = widgets.NewParagraph()
	return &t
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
