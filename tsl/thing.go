package tsl

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type EntityRequest struct {
	ID        string          `json:"id"`      // 消息ID，String类型的数字，取值范围0~4294967295，且每个消息ID在当前设备中具有唯一性。
	Version   string          `json:"version"` // 协议版本号，目前协议版本号唯一取值为1.0。
	Params    json.RawMessage `json:"params,omitempty"`
	Method    string          `json:"method"`
	Timestamp int64           `json:"timestamp"`
	// Sys        interface{} `json:"sys,omitempty"` //	扩展功能的参数，其下包含各功能字段。
}

type EntityReply struct {
	ID        string          `json:"id"` // 消息ID，String类型的数字，取值范围0~4294967295，且每个消息ID在当前设备中具有唯一性。
	Code      int             `json:"code"`
	Data      json.RawMessage `json:"data,omitempty"`
	Method    string          `json:"method"`
	Timestamp int64           `json:"timestamp"`
	// Sys        interface{} `json:"sys,omitempty"` //	扩展功能的参数，其下包含各功能字段。
}

type ThingEntity struct {
	ID        string          `json:"id"`      // 消息ID，String类型的数字，取值范围0~4294967295，且每个消息ID在当前设备中具有唯一性。
	Version   string          `json:"version"` // 协议版本号，目前协议版本号唯一取值为1.0。
	Params    json.RawMessage `json:"params,omitempty"`
	Method    string          `json:"method"`
	Timestamp int64           `json:"timestamp"`
	Code      int             `json:"code"`
	Data      json.RawMessage `json:"data,omitempty"`
	// Sys        interface{} `json:"sys,omitempty"` //	扩展功能的参数，其下包含各功能字段。
}

// 物模型
type Thing struct {
	Profile *Profile
	// 事件
	Events []*Event
	// 服务
	Services []*Service
	// 属性
	Properties []*Property
	Value      struct {
		// 事件
		Events map[string]*Event
		// 服务
		Services map[string]*Service
	}
}

// NewThing
func NewThing(bs []byte) (*Thing, error) {
	var thing Thing
	err := json.Unmarshal(bs, &thing)
	if err != nil {
		return nil, err
	}
	err = thing.ValidateSpec()
	if err != nil {
		return nil, err
	}
	return &thing, nil
}

func (s *Thing) init() {
	if s.Value.Events == nil {
		s.Value.Events = make(map[string]*Event)
		for _, event := range s.Events {
			s.Value.Events[event.Identifier] = event
		}
		// 添加属性上报事件
		s.Value.Events["post"] = &Event{
			Identifier: "post",
			Name:       "属性上报",
			Desc:       "",
			Method:     "thing.event.property.post",
			Type:       "info",
			OutputData: s.Properties,
		}
	}
	if s.Value.Services == nil {
		s.Value.Services = make(map[string]*Service)
		for _, service := range s.Services {
			s.Value.Services[service.Identifier] = service
		}
		// 添加属性设置服务
		serviceSetProperties := []*Property{}
		for _, v := range s.Properties {
			if v.AccessMode == "rw" {
				serviceSetProperties = append(serviceSetProperties, v)
			}
		}
		s.Value.Services["set"] = &Service{
			Identifier: "set",
			Name:       "属性设置",
			Desc:       "",
			Method:     "thing.service.property.set",
			CallType:   "sync",
			Required:   true,
			InputData:  serviceSetProperties,
		}
		s.Value.Services["get"] = &Service{
			Identifier: "get",
			Name:       "属性获取",
			Desc:       "",
			Method:     "thing.service.property.get",
			CallType:   "sync",
			Required:   true,
			OutputData: s.Properties,
		}
	}
}

func (s *Thing) ValidateSpec() error {
	var err error
	if s.Profile == nil {
		return fmt.Errorf("Thing Profile is nil")
	}
	err = s.Profile.ValidateSpec()
	if err != nil {
		return fmt.Errorf("profile.%v", err)
	}
	for k, event := range s.Events {
		err = event.ValidateSpec()
		if err != nil {
			return fmt.Errorf("events[%d].%v", k, err)
		}
	}
	for k, service := range s.Services {
		err = service.ValidateSpec()
		if err != nil {
			return fmt.Errorf("services[%d].%v", k, err)
		}
	}
	for k, property := range s.Properties {
		err = property.ValidateSpec()
		if err != nil {
			return fmt.Errorf("properties[%d,(%s)].%v", k, property.Identifier, err)
		}
	}
	return nil
}

func (s *Thing) ToEntityString() string {
	var m struct {
		Events   []*ThingEntity `json:"events"`
		Services []*ThingEntity `json:"services"`
	}
	s.init() // initialize
	for _, v := range s.Value.Services {
		m.Services = append(m.Services, v.ToEntity())
	}
	for _, v := range s.Value.Events {
		m.Events = append(m.Events, v.ToEntity())
	}
	bs, _ := json.MarshalIndent(m, "", "  ")
	return string(bs)
}

type Profile struct {
	ProductKey string
	DeviceName string
}

func (s *Profile) ValidateSpec() error {
	if s.ProductKey == "" {
		return fmt.Errorf("productKey err: productKey is empty")
	}
	return nil
}

// 事件
type Event struct {
	Identifier string
	Name       string
	Desc       string
	Method     string
	Type       string
	OutputData []*Property
	Value      struct {
		OutputData map[string]*Property
	}
}

func (s *Event) init() {
	if len(s.OutputData) != 0 && s.Value.OutputData == nil {
		s.Value.OutputData = propertiesToMap(s.OutputData)
	}
}

func (s *Event) ValidateSpec() error {
	var err error
	if s.Identifier == "" {
		return fmt.Errorf("identifier err: identifier is empty")
	}
	if s.Name == "" {
		return fmt.Errorf("name err: name is empty")
	}
	if s.Method == "" {
		return fmt.Errorf("method err: method is empty")
	}
	if !strings.HasPrefix(s.Method, "thing.event.") {
		return fmt.Errorf("method err: method is thing.event.*")
	}
	_, err = NewThingMethod(s.Method)
	if err != nil {
		return fmt.Errorf("method err: %v", err)
	}
	for k, v := range s.OutputData {
		err = v.ValidateSpec()
		if err != nil {
			return fmt.Errorf("outputData[%d].%v", k, err)
		}
	}
	return nil
}

func (s *Event) ToEntity() *ThingEntity {
	outputData := propertyToEntityMap(s.OutputData)
	outputBytes, _ := json.Marshal(outputData)
	methodStrs := []string{
		s.Method,
	}
	if s.Name != "" {
		methodStrs = append(methodStrs, s.Name)
	}
	if s.Desc != "" {
		methodStrs = append(methodStrs, s.Desc)
	}
	return &ThingEntity{
		ID:        "int64,消息id",
		Version:   "1.0",
		Timestamp: time.Now().UnixMilli(),
		Params:    outputBytes, // event 为上报, 参数到平台放outputData中
		Method:    strings.Join(methodStrs, ", "),
	}
}

// 服务
type Service struct {
	Identifier string
	Name       string
	Desc       string
	Method     string
	CallType   string
	Required   bool
	InputData  []*Property
	OutputData []*Property
	Value      struct {
		InputData  map[string]*Property
		OutputData map[string]*Property
	}
}

func (s *Service) init() {
	if len(s.InputData) != 0 && s.Value.InputData == nil {
		s.Value.InputData = propertiesToMap(s.InputData)
	}
	if len(s.OutputData) != 0 && s.Value.OutputData == nil {
		s.Value.OutputData = propertiesToMap(s.OutputData)
	}
}

func (s *Service) ValidateSpec() error {
	s.init() // 初始化
	var err error
	if s.Identifier == "" {
		return fmt.Errorf("identifier err: identifier is empty")
	}
	if s.Name == "" {
		return fmt.Errorf("name err: name is empty")
	}
	if s.CallType == "" {
		return fmt.Errorf("callType err: callType is empty")
	}
	if s.Method == "" {
		return fmt.Errorf("method err: method is empty")
	}
	if !strings.HasPrefix(s.Method, "thing.service.") {
		return fmt.Errorf("method err: method is thing.service.*")
	}
	_, err = NewThingMethod(s.Method)
	if err != nil {
		return fmt.Errorf("method err: %v", err)
	}
	for k, v := range s.InputData {
		err = v.ValidateSpec()
		if err != nil {
			return fmt.Errorf("inputData[%d].%v", k, err)
		}
	}
	for k, v := range s.OutputData {
		err = v.ValidateSpec()
		if err != nil {
			return fmt.Errorf("outputData[%d].%v", k, err)
		}
	}
	return nil
}

func (s *Service) ToEntity() *ThingEntity {
	inputData := propertyToEntityMap(s.InputData)
	outputData := propertyToEntityMap(s.OutputData)
	inputBytes, _ := json.Marshal(inputData)
	outputBytes, _ := json.Marshal(outputData)
	methodStrs := []string{
		s.Method,
	}
	if s.Name != "" {
		methodStrs = append(methodStrs, s.Name)
	}
	if s.Desc != "" {
		methodStrs = append(methodStrs, s.Desc)
	}
	return &ThingEntity{
		ID:        "int64,消息id",
		Version:   "1.0",
		Timestamp: time.Now().UnixMilli(),
		Params:    inputBytes,
		Data:      outputBytes,
		Method:    strings.Join(methodStrs, ","),
	}
}
