package topic

import (
	"errors"
	"strings"
)

const (
	EventTopic = iota
	ServiceTopic
	PropertyTopic
)

const (
	LogClassify1 = iota
	otaClassify1
	ThingClassify1
	subClassify1
)

var (
	TopicPrefix = []string{
		"/sys/",
		"/ext/session/",
	}
	// Classify1 = map[string][int]{
	// 	"log": ,
	// 	"ota",
	// 	"thing",
	// 	"sub",
	// }
	TopicReplySuffix      = "reply"
	ErrInvalidTopicPrefix = errors.New("invalid topic prefix")
	ErrInvalidTopic       = errors.New("invalid topic")
)

type Topic struct {
	OrigTopic  string
	Prefix     string
	ProductKey string
	DeviceName string
	Classify1  string
	Classify2  string
	SubDirs    []string
	IsReply    bool
}

func ParseTopic(topic string) (Topic, error) {
	var (
		t        Topic
		subTopic string
	)
	t.OrigTopic = topic
	for _, prefix := range TopicPrefix {
		if strings.HasPrefix(topic, prefix) {
			t.Prefix = prefix
			subTopic = topic[len(prefix):]
			break
		}
	}
	if t.Prefix == "" {
		return t, ErrInvalidTopicPrefix
	}

	if strings.HasSuffix(subTopic, TopicReplySuffix) {
		subTopic = subTopic[:len(subTopic)-len(TopicReplySuffix)-1]
		t.IsReply = true
	}

	subSlice := strings.Split(subTopic, "/")
	if len(subSlice) < 4 {
		return t, ErrInvalidTopic
	}
	t.ProductKey = subSlice[0]
	t.DeviceName = subSlice[1]
	t.Classify1 = subSlice[2]
	t.Classify2 = subSlice[3]
	t.SubDirs = subSlice[4:]
	return t, nil
}
