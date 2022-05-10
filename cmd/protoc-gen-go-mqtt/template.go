package main

import (
	"bytes"
	"strings"
	"text/template"
)

var mqttTemplate = `
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
type {{.ServiceType}}MQTTServer interface {
{{- range .MethodSets}}
	{{.Name}}(context.Context, paho_mqtt_golang.Client, paho_mqtt_golang.Message, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

func Register{{.ServiceType}}MQTTServer(r *mqttrouter.Router, srv {{.ServiceType}}MQTTServer) {
	{{- range .Methods}}
	r.Handle("{{.Path}}", 0, _{{$svrType}}_{{.Name}}{{.Num}}_MQTT_Handler(srv))
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_MQTT_Handler(srv {{$svrType}}MQTTServer) func(context.Context, paho_mqtt_golang.Client, paho_mqtt_golang.Message)  {
	return func(ctx context.Context, c paho_mqtt_golang.Client, msg paho_mqtt_golang.Message)  {
		var in {{.Request}}
		vars := mqttrouter.ParamsFromContext(ctx)
		bs, _ := json.Marshal(vars)
		err := json.Unmarshal(bs, &in)
		if err != nil {
			log.Error("var Params error:", err)
		}
		err = json.Unmarshal(msg.Payload(), &in)
		if err != nil {
			log.Error("message error:", err)
			return
		}
		err = in.Validate()
		if err != nil {
			log.Error("validate error:", err)
			return
		}
		reply, err := srv.{{.Name}}(ctx, c, msg, &in)
		if err != nil {
			log.Error("{{.Name}} error:", err)
			return
		}
		if reply == nil {
			return
		}
		bs, err = json.Marshal(reply)
		if err != nil {
			log.Errorf("topic:%v, err: %v", msg.Topic(), err)
			return
		} else {
			log.Debugf("reply mqtt topic:%v,body: %v", msg.Topic(), string(bs))
		}
		c.Publish(msg.Topic() + "_reply", 0, false, bs)
	}
}
{{end}}

`

type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/helloworld/helloworld.proto
	Methods     []*methodDesc
	MethodSets  map[string]*methodDesc
}

type methodDesc struct {
	// method
	Name    string
	Num     int
	Request string
	Reply   string
	// http_rule
	Path         string
	Method       string
	HasVars      bool
	HasBody      bool
	Body         string
	ResponseBody string
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("mqtt").Parse(strings.TrimSpace(mqttTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}
