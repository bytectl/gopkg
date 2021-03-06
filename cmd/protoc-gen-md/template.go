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
{{if .Description}}
{{.Description}}
{{end}}
- 请求路径: INNERLINEBLOCK{{.Method}} {{.Path}}INNERLINEBLOCK
- 请求参数
NEWLINE
{{if .Params}}CODEBLOCKjson
{{.Params}}
CODEBLOCK{{else if .TableParams}}|参数名 | 类型 | 可选 | 说明|
|------|------|-----|-----|
{{.TableParams}}{{else}}无NEWLINE{{end}}

- 返回参数
NEWLINE
{{if .ResponseParams}}CODEBLOCKjson
{{.ResponseParams}}
CODEBLOCK{{else}}无{{end}}
NEWLINE
- 请求示例
NEWLINE
CODEBLOCKhttp
{{.Method}} {{$hostString}}{{.Path}}{{.PathParams}}
Content-Type: application/json
NEWLINE
{{.Params}}
CODEBLOCK

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
	Name           string
	Comment        string
	Description    string
	Params         string
	ResponseParams string
	PathParams     string
	TableParams    string
	Num            int
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
	tmpl, err := template.New("md").Parse(strings.TrimSpace(restTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	src := strings.Replace(buf.String(), "\t", "", -1)
	// src = strings.Replace(src, "\r\n", "", -1)
	src = strings.Replace(src, "NEWLINE", "", -1)
	src = strings.Replace(src, "CODEBLOCK", "```", -1)
	src = strings.Replace(src, "INNERLINEBLOCK", "`", -1)
	return src
}
