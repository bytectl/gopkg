package topic

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseTopic(t *testing.T) {
	tests := []struct {
		topic   string
		want    Topic
		wantErr bool
	}{
		{
			"/sys/{productKey}/{deviceName}/thing/sub/register",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/sub/register",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "sub",
				SubDirs:    []string{"register"},
				IsReply:    false,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/sub/register_reply",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/sub/register_reply",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "sub",
				SubDirs:    []string{"register"},
				IsReply:    true,
			},
			false,
		},
		{
			"/ext/session/{productKey}/{deviceName}/combine/login",
			Topic{
				OrigTopic:  "/ext/session/{productKey}/{deviceName}/combine/login",
				Prefix:     "/ext/session/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "combine",
				Classify2:  "login",
				SubDirs:    []string{},
				IsReply:    false,
			},
			false,
		},
		{
			"/ext/session/{productKey}/{deviceName}/combine/login_reply",
			Topic{
				OrigTopic:  "/ext/session/{productKey}/{deviceName}/combine/login_reply",
				Prefix:     "/ext/session/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "combine",
				Classify2:  "login",
				SubDirs:    []string{},
				IsReply:    true,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/discovery/open",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/discovery/open",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "discovery",
				Classify2:  "open",
				SubDirs:    []string{},
				IsReply:    false,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/topo/add",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/topo/add",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "topo",
				SubDirs:    []string{"add"},
				IsReply:    false,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/topo/add_reply",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/topo/add_reply",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "topo",
				SubDirs:    []string{"add"},
				IsReply:    true,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/event/property/post",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/event/property/post",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "event",
				SubDirs:    []string{"property", "post"},
				IsReply:    false,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/event/property/post_reply",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/event/property/post_reply",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "event",
				SubDirs:    []string{"property", "post"},
				IsReply:    true,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/service/property/get",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/service/property/get",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "service",
				SubDirs:    []string{"property", "get"},
				IsReply:    false,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/service/property/get_reply",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/service/property/get_reply",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "service",
				SubDirs:    []string{"property", "get"},
				IsReply:    true,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/service/{tsl.service.identifier}",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/service/{tsl.service.identifier}",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "service",
				SubDirs:    []string{"{tsl.service.identifier}"},
				IsReply:    false,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/service/{tsl.service.identifier}_reply",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/service/{tsl.service.identifier}_reply",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "service",
				SubDirs:    []string{"{tsl.service.identifier}"},
				IsReply:    true,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/property/desired/get",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/property/desired/get",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "property",
				SubDirs:    []string{"desired", "get"},
				IsReply:    false,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/property/desired/get_reply",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/property/desired/get_reply",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "property",
				SubDirs:    []string{"desired", "get"},
				IsReply:    true,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/delete",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/delete",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "delete",
				SubDirs:    []string{},
				IsReply:    false,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/thing/deviceinfo/delete",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/thing/deviceinfo/delete",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "thing",
				Classify2:  "deviceinfo",
				SubDirs:    []string{"delete"},
				IsReply:    false,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/ota/device/check",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/ota/device/check",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "ota",
				Classify2:  "device",
				SubDirs:    []string{"check"},
				IsReply:    false,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/log/upload",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/log/upload",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "log",
				Classify2:  "upload",
				SubDirs:    []string{},
				IsReply:    false,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/log/upload/reply",
			Topic{
				OrigTopic:  "/sys/{productKey}/{deviceName}/log/upload/reply",
				Prefix:     "/sys/",
				ProductKey: "{productKey}",
				DeviceName: "{deviceName}",
				Classify1:  "log",
				Classify2:  "upload",
				SubDirs:    []string{},
				IsReply:    true,
			},
			false,
		},
		{
			"/sys/{productKey}/{deviceName}/log",
			Topic{
				OrigTopic: "/sys/{productKey}/{deviceName}/log",
				Prefix:    "/sys/",
			},
			true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			got, err := ParseTopic(tt.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTopic error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTopic got = %v, want %v", got, tt.want)
			}
		})
	}
}
