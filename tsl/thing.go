package tsl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type ThingEntity struct {
	ID        string          `json:"id"`      // 消息ID，String类型的数字，取值范围0~4294967295，且每个消息ID在当前设备中具有唯一性。
	Version   string          `json:"version"` // 协议版本号，目前协议版本号唯一取值为1.0。
	Params    json.RawMessage `json:"params,omitempty"`
	Method    string          `json:"method"`
	Timestamp string          `json:"timestamp"`
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

// 校验实体数据, 传入实体字节数据
func (s *Thing) ValidateEntityBytes(bs []byte) error {
	var thingEntity ThingEntity
	err := json.Unmarshal(bs, &thingEntity)
	if err != nil {
		return err
	}
	return s.ValidateEntity(&thingEntity)
}

// 校验实体数据 传入为结构体
func (s *Thing) ValidateEntity(thingEntity *ThingEntity) error {

	var err error
	if thingEntity == nil {
		return fmt.Errorf("thingEntity is nil")
	}

	method, err := NewThingMethod(thingEntity.Method)
	if err != nil {
		return err
	}
	if method.IsService() {
		err = s.ValidateService(method.Action, thingEntity.Params, thingEntity.Data)
	} else if method.IsEvent() {
		err = s.ValidateEvent(method.Action, thingEntity.Params)
	} else {
		err = fmt.Errorf("thingEntity.method(%s) no service or event", thingEntity.Method)
	}
	return err
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

func (s *Thing) ValidateEvent(identifier string, params []byte) error {
	s.init() // initialize
	event, ok := s.Value.Events[identifier]
	if !ok {
		return fmt.Errorf("event.identifier: (%s) no found", identifier)
	}
	err := event.ValidateEntity(params)
	if err != nil {
		return fmt.Errorf("events[%s].%v", identifier, err)
	}
	return nil
}

func (s *Thing) ValidateService(identifier string, params, data []byte) error {
	s.init() // initialize
	service, ok := s.Value.Services[identifier]
	if !ok {
		return fmt.Errorf("service.identifier: (%s) no found", identifier)
	}
	err := service.ValidateEntity(params, data)
	if err != nil {
		return fmt.Errorf("services[%s].%v", identifier, err)
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
	for k, v := range s.OutputData {
		err = v.ValidateSpec()
		if err != nil {
			return fmt.Errorf("outputData[%d].%v", k, err)
		}
	}
	return nil
}
func (s *Event) ValidateEntity(outputData []byte) error {
	var err error
	s.init() // initialize
	if outputData != nil {
		err = validateEntityParams(s.Value.OutputData, outputData)
		if err != nil {
			return fmt.Errorf("outputData.%v", err)
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
		Timestamp: "时间戳",
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

func (s *Service) ValidateEntity(inputData, outputData []byte) error {
	var err error
	s.init() // initialize
	if inputData != nil {
		err = validateEntityParams(s.Value.InputData, inputData)
		if err != nil {
			return fmt.Errorf("inputData.%v", err)
		}
	}
	if outputData != nil {
		err = validateEntityParams(s.Value.OutputData, outputData)
		if err != nil {
			return fmt.Errorf("outputData.%v", err)
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
		Timestamp: "时间戳",
		Params:    inputBytes,
		Data:      outputBytes,
		Method:    strings.Join(methodStrs, ","),
	}
}

func validateEntityParams(specData map[string]*Property, data []byte) error {
	var err error
	if specData == nil {
		return fmt.Errorf("validateEntityParams: specData==nil")
	}
	paramMap := make(map[string]interface{})
	decoder := json.NewDecoder(bytes.NewReader(data))
	// 使用json number
	decoder.UseNumber()
	if err := decoder.Decode(&paramMap); err != nil {
		return err
	}
	for k, v := range paramMap {
		param, ok := specData[k]
		if !ok {
			return fmt.Errorf("[%s] err:  not exist", k)
		}
		err = param.ValidateValue(v)
		if err != nil {
			return fmt.Errorf("[%s].%v", k, err)
		}
	}
	return nil
}
