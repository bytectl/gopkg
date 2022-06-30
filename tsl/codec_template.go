package tsl

var CodecPrefix = `// xxx 产品编解码插件
package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type Event struct {
	Params map[string]interface{} CODEBLOCKjson:"params"CODEBLOCK
	Method string                 CODEBLOCKjson:"method"CODEBLOCK
}
type EventData struct {
	Events []Event CODEBLOCKjson:"events"CODEBLOCK
}

type Params map[string]interface{}

`
var CodecTail = `
// 解码	
func Decode(payload, metadata []byte) ([]byte, error) {

	buffer := bytes.NewBuffer(payload)
	var decodeData struct {
		//you decode fields   uint8
	
	}
	if err := binary.Read(buffer, binary.BigEndian, &decodeData); err != nil {
		return nil, err
	}

	params := make(Params)
	// 构建参数
	//params.Setxxx
	
	// ... 

	eventData := &EventData{
		Events: []Event{
			{
				Params: params,
				// Method: EventPropertyMethod,
			},
		},
	}
	return json.Marshal(eventData)
}

// Encode 编码
func Encode(data, metadata []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implement")
}

`
