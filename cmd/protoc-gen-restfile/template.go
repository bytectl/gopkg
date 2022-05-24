package main

import (
	"bytes"
	"strings"
	"text/template"
)

var restTemplate = `
{{$hostString := .HostString}}
{{- range .Methods}}

### {{.Comment}}
{{.Method}} {{$hostString}}{{.Path}}{{.PathParams}}
Content-Type: application/json
{{.Params}}
{{- end}}
`

type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/helloworld/helloworld.proto
	Methods     []*methodDesc
	MethodSets  map[string]*methodDesc
	HostString  string
}

type methodDesc struct {
	// method
	Name       string
	Comment    string
	Params     string
	PathParams string
	Num        int
	// http_rule
	Path         string
	Method       string
	Body         string
	ResponseBody string
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("restfile").Parse(strings.TrimSpace(restTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return buf.String()
}
