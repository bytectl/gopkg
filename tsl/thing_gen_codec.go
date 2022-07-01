package tsl

import (
	"bytes"
	"go/format"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func (s *Thing) GenerateGoCodec(tp string) (string, error) {

	s.init() // initialize
	buf := new(bytes.Buffer)
	fm := template.FuncMap{
		"title": strings.Title,
	}
	tmpl, err := template.New("thing").Funcs(sprig.TxtFuncMap()).Funcs(fm).Parse(strings.TrimSpace(tp))
	if err != nil {
		return "", err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return "", err
	}

	str := buf.String()
	str = strings.ReplaceAll(str, "CODEBLOCK", "`")
	bs, err := format.Source([]byte(str))
	return string(bs), err
}

func (s *Event) ConstName() string {
	if s.Identifier == "post" {
		return "Property"
	}
	return s.Identifier
}

func (s *Event) ParamPrefixName() string {
	if s.Identifier == "post" {
		return ""
	}
	return s.Identifier
}
