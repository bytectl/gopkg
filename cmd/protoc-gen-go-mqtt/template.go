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
	{{.Name}}(mqtt.Context,*{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}
func SetLogger(logger log.Logger){
	glog = log.NewHelper(logger)
}

func Subscribe{{.ServiceType}}(c paho_mqtt_golang.Client, r *mqtt.Router) {
	{{- range .Methods}}
	r.Subscribe(c,"{{.Path}}",0)
	{{- end}}
}

func Register{{.ServiceType}}MQTTServer(r *mqtt.Router, srv {{.ServiceType}}MQTTServer) {
	{{- range .MethodSets}}
	r.Handle("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_MQTT_Handler(srv))
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_MQTT_Handler(srv {{$svrType}}MQTTServer) func(mqtt.Context)  {
	return func(ctx mqtt.Context)  {
		glog.Debugf("receive mqtt topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in :=&{{.Request}}{}
		err := jsonCodec.Unmarshal(ctx.Message().Payload(), in)
		if err != nil {
			glog.Error("message error:", err)
			return
		}
		vars := mqtt.ParamsFromContext(ctx)
		err = formCodec.Unmarshal([]byte(vars.Encode()), in)
		if err != nil {
			glog.Error("var Params error:", err)
		}
		err = in.Validate()
		if err != nil {
			glog.Error("validate error:", err)
			return
		}
		glog.Debugf("receive mqtt request:%+v",in)
		reply, err := srv.{{.Name}}(ctx, in)
		if err != nil {
			glog.Error("{{.Name}} error:", err)
		}
		if reply == nil {
			return
		}
		bs, err := jsonCodec.Marshal(reply)
		if err != nil {
			glog.Errorf("topic:%v, err: %v", ctx.Message().Topic(), err)
			return
		} 
		replyTopic := strings.TrimPrefix(ctx.Message().Topic(),ServerTopicPrefix) 
		replyTopic = fmt.Sprintf("%s%s_reply", DeviceTopicPrefix, replyTopic)
		ctx.Client().Publish(replyTopic, 0, false, bs)
		glog.Debugf("reply mqtt topic:%v,body: %v", replyTopic, string(bs))
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
