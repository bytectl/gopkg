package tsl

import (
	"fmt"
	"strings"
)

const (
	MethodCellMinLength      = 3
	MethodCellMaxLength      = 4
	MethodClassifyIndex      = 1
	MethodActionIndex        = 2
	MethodActionPropertyName = "property"
	MethodServiceName        = "service"
	MethodEventName          = "event"
)

type ThingMethod struct {
	Original   string
	Classify   string // 分类: 属性, 服务, 事件
	Action     string // 动作: set, get, post, {tsl.identifier}
	IsProperty bool   // 是否是属性
	IsSet      bool   // 是否是set
	IsGet      bool   // 是否是get
}

func NewThingMethod(method string) (*ThingMethod, error) {
	isProperty := false
	isSet := false
	isGet := false
	strs := strings.Split(method, ".")
	if len(strs) < MethodCellMinLength {
		return nil, fmt.Errorf("method(%s) is invalid", method)
	}
	action := strs[MethodActionIndex]
	if strings.Compare(action, MethodActionPropertyName) == 0 {
		if len(strs) < MethodCellMaxLength {
			return nil, fmt.Errorf("method(%s) is invalid", method)
		}
		action = strs[MethodActionIndex+1]
		isProperty = true
		if 0 == strings.Compare(action, "set") {
			isSet = true
		}
		if 0 == strings.Compare(action, "get") {
			isGet = true
		}
	}
	return &ThingMethod{
		Original:   method,
		Classify:   strs[MethodClassifyIndex],
		Action:     action,
		IsProperty: isProperty,
		IsSet:      isSet,
		IsGet:      isGet,
	}, nil
}

func (t *ThingMethod) IsService() bool {
	return strings.Compare(t.Classify, MethodServiceName) == 0
}
func (t *ThingMethod) IsEvent() bool {
	return strings.Compare(t.Classify, MethodEventName) == 0
}
