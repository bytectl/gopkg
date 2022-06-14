package tsl

import (
	"fmt"
	"strings"
)

const (
	MethodCellLength         = 4
	MethodClassifyIndex      = 1
	MethodActionIndex        = 2
	MethodActionPropertyName = "property"
	MethodServiceName        = "service"
	MethodEventName          = "event"
)

type ThingMethod struct {
	Original string
	Classify string // 分类: 属性, 服务, 事件
	Action   string // 动作: set, get, post, {tsl.identifier}
}

func NewThingMethod(method string) (*ThingMethod, error) {

	strs := strings.Split(method, ".")
	if len(strs) < MethodCellLength {
		return nil, fmt.Errorf("method(%s) is invalid", method)
	}
	action := strs[MethodActionIndex]
	if strings.Compare(action, MethodActionPropertyName) == 0 {
		action = strs[MethodActionIndex+1]
	}
	return &ThingMethod{
		Original: method,
		Classify: strs[MethodClassifyIndex],
		Action:   action,
	}, nil
}

func (t *ThingMethod) IsService() bool {
	return strings.Compare(t.Classify, MethodServiceName) == 0
}
func (t *ThingMethod) IsEvent() bool {
	return strings.Compare(t.Classify, MethodEventName) == 0
}
