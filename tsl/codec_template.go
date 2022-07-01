package tsl

var DefaultCodecTmpl = `// {{ .Profile.ProductKey }} 产品编解码插件
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

const ({{range $index, $element := .Value.Events }}
	Event{{ title $element.ConstName }} = "{{$element.Method}}" // {{$element.Name}}事件{{end}} 
)

{{- range $index, $element := .Value.Events }}
{{- if $element.OutputData }}
// Event{{ title $element.ConstName }} {{$element.Name}}事件参数{{ end }}
{{$prefixName := title $element.ParamPrefixName}}
{{- range $element.OutputData}}func (p Params) Set{{ $prefixName }}{{- title .Identifier }}(v {{.DataType.GenerateGoType -}}) { p["{{- .Identifier -}}"] = v } // {{.Name}}
{{end -}}{{end -}}

// 解码	
func Decode(payload, metadata []byte) ([]byte, error) {

	buffer := bytes.NewBuffer(payload)
	var decodeData struct {
		//TODO: please change to you decode fields
		{{- range .Value.Events}}


		{{$prefixName := title .ParamPrefixName}}
		{{- if not $prefixName }}
		{{- range .OutputData}}{{- title .Identifier }}  {{.DataType.GenerateGoType}} 
		{{end -}}{{end -}}{{end -}}
	}
	if err := binary.Read(buffer, binary.BigEndian, &decodeData); err != nil {
		return nil, err
	}

	params := make(Params)
	//TODO: please set your params
	{{- range .Value.Events}}
	{{$prefixName := title .ParamPrefixName}}
	{{- if not $prefixName }}
	{{- range .OutputData}}params.Set{{- title .Identifier }}({{- title .Identifier }} ) // {{.Name}}
	{{end -}}{{end -}}{{end -}}
	//TODO: please make up your other event params

	eventData := &EventData{
		Events: []Event{
			{
				Params: params,
				//TODO: please change to you need event method
				Method: EventPropertyMethod,
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
