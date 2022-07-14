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

func Subscribe{{.ServiceType}}(c paho_mqtt_golang.Client, m *mqtt.MQTTSubscribe) {
	{{- range .SubscribeMethods}}
	{{if .Path }}m.Subscribe(c,"{{.Path}}",0){{end}}
	{{- end}}
}

func Register{{.ServiceType}}MQTTServer(s *mqtt.Server, srv {{.ServiceType}}MQTTServer) {
	r := s.Route()
	{{- range .MethodSets}}
	r.Handle("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_MQTT_Handler(srv))
	{{- end}}
}

{{range .MethodSets}}
func _{{$svrType}}_{{.Name}}{{.Num}}_MQTT_Handler(srv {{$svrType}}MQTTServer) func(mqtt.Context)  {
	return func(ctx mqtt.Context)  {
		glog.Debugf("receive mqtt topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in :=&{{.Request}}{}
		err := ctx.Bind(in)
		if err != nil {
			glog.Error("message error:", err)
			return
		}
		err = ctx.BindVars(in)
		if err != nil {
			glog.Error("var Params error:", err)
		}
		glog.Debugf("receive mqtt topic:%v, in: %+v", ctx.Message().Topic(), in)
		err = in.Validate()
		if err != nil {
			glog.Error("validate error:", err)
			return
		}
		glog.Debugf("receive mqtt request:%+v",in)
		reply, err := srv.{{.Name}}(ctx, in)
		if reply == nil {
			glog.Debugf(" mqtt topic:%v, no need reply", ctx.Message().Topic())
			if err != nil {
				glog.Error("{{.Name}} error:", err)
			}
			return
		}
		if err != nil {
			glog.Error("{{.Name}} error:", err)
			ctx.ReplyErr(err)
			return
		}
		err = ctx.Reply(reply)
		if err != nil {
			glog.Error("{{.Name}} error:", err)
			ctx.ReplyErr(err)
			return
		}
	}
}
{{end}}

`

type serviceDesc struct {
	ServiceType      string // Greeter
	ServiceName      string // helloworld.Greeter
	Metadata         string // api/helloworld/helloworld.proto
	Methods          []*methodDesc
	MethodSets       map[string]*methodDesc
	SubscribeMethods []*methodDesc // 用于订阅topic
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
