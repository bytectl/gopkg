package mqtt

import (
	"fmt"
	"strings"

	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// SupportPackageIsVersion1 These constants should not be referenced from any other code.
const SupportPackageIsVersion1 = true
const (
	ServerTopicPrefix = "/sys"
	DeviceTopicPrefix = "/device"
)

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func([]byte, interface{}) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(c pmqtt.Client, topic string, v interface{}) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(c pmqtt.Client, topic string, err error)

// DefaultRequestDecoder decodes the request body to object.
func DefaultRequestDecoder(data []byte, v interface{}) error {
	if len(data) == 0 {
		return nil
	}
	codec := encoding.GetCodec(json.Name)
	if err := codec.Unmarshal(data, v); err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	return nil
}

// DefaultResponseEncoder encodes the object to the mqtt reply.
func DefaultResponseEncoder(c pmqtt.Client, topic string, v interface{}) error {
	if v == nil {
		return nil
	}
	codec := encoding.GetCodec(json.Name)
	body, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	replyTopic := makeReplyTopic(topic)
	log.Debugf("reply mqtt topic:%v,body: %v", replyTopic, string(body))
	if c != nil {
		c.Publish(replyTopic, 0, false, body)
	}
	return nil
}

// DefaultErrorEncoder encodes the error to the mqtt response.
func DefaultErrorEncoder(c pmqtt.Client, topic string, err error) {
	var reply struct {
		Id      string `json:"id"`
		Code    int32  `json:"code"`
		Reason  string `json:"reason,omitempty"`
		Message string `json:"message"`
	}
	se := errors.FromError(err)
	reply.Id = se.Metadata["id"]
	reply.Code = se.Code
	reply.Message = se.Message
	reply.Reason = se.Reason
	codec := encoding.GetCodec(json.Name)
	body, _ := codec.Marshal(reply)
	replyTopic := makeReplyTopic(topic)
	log.Debugf("reply mqtt topic:%v,body: %v", replyTopic, string(body))
	if c != nil {
		c.Publish(replyTopic, 0, false, body)
	}
}

func makeReplyTopic(topic string) string {
	replyTopic := topic
	if strings.HasPrefix(topic, ServerTopicPrefix) {
		replyTopic = strings.TrimPrefix(topic, ServerTopicPrefix)
		replyTopic = fmt.Sprintf("%s%s_reply", DeviceTopicPrefix, replyTopic)
	} else if strings.HasPrefix(topic, DeviceTopicPrefix) {
		replyTopic = strings.TrimPrefix(topic, DeviceTopicPrefix)
		replyTopic = fmt.Sprintf("%s%s_reply", ServerTopicPrefix, replyTopic)
	}
	return replyTopic
}
