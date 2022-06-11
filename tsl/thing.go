package tsl

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// 校验接口
type Validator interface {
	Validate() error
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
}

func (s *Thing) Validate() error {
	var err error
	if s.Profile == nil {
		return fmt.Errorf("Thing Profile is nil")
	}
	err = s.Profile.Validate()
	if err != nil {
		return fmt.Errorf("profile.%v", err)
	}
	for k, event := range s.Events {
		err = event.Validate()
		if err != nil {
			return fmt.Errorf("events[%d].%v", k, err)
		}
	}
	for k, service := range s.Services {
		err = service.Validate()
		if err != nil {
			return fmt.Errorf("services[%d].%v", k, err)
		}
	}
	for k, property := range s.Properties {
		err = property.Validate()
		if err != nil {
			return fmt.Errorf("properties[%d].%v", k, err)
		}
	}
	return nil
}

type Profile struct {
	ProductKey string
	DeviceName string
}

func (s *Profile) Validate() error {
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
}

func (s *Event) Validate() error {
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
		err = v.Validate()
		if err != nil {
			return fmt.Errorf("outputData[%d].%v", k, err)
		}
	}

	return nil
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
}

func (s *Service) Validate() error {
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
		err = v.Validate()
		if err != nil {
			return fmt.Errorf("inputData[%d].%v", k, err)
		}
	}
	for k, v := range s.OutputData {
		err = v.Validate()
		if err != nil {
			return fmt.Errorf("outputData[%d].%v", k, err)
		}
	}
	return nil
}

// 属性
type Property struct {
	AccessMode string
	Identifier string
	Name       string
	Desc       string
	Required   bool
	DataType   *DataType
}

func (s *Property) Validate() error {
	if s.Identifier == "" {
		return fmt.Errorf("identifier err: identifier is empty")
	}
	if s.Name == "" {
		return fmt.Errorf("name  err: name is empty")
	}
	if s.DataType == nil {
		return fmt.Errorf("dataType err: dataType is empty")
	}
	err := s.DataType.Validate()
	if err != nil {
		return fmt.Errorf("dataType.%v", err)
	}
	return nil
}

// 数据类型
type DataType struct {
	// 数据类型
	Type string
	// 数据类型定义
	Specs json.RawMessage
}

func validateSpec(bs []byte, spec Validator) error {
	if len(bs) == 0 {
		return fmt.Errorf("spec is empty")
	}
	err := json.Unmarshal(bs, spec)
	if err != nil {
		return err
	}
	return spec.Validate()
}
func (s *DataType) Validate() error {
	var err error
	bs := s.Specs
	switch s.Type {
	case "int", "float", "double":
		err = validateSpec(bs, &DigitalSpec{})
	case "text":
		err = validateSpec(bs, &TextSpec{})
	case "array":
		err = validateSpec(bs, &ArraySpec{})
	case "struct":
		var structDataSpec []*StructDataSpec
		err = json.Unmarshal(bs, &structDataSpec)
		if err != nil {
			return fmt.Errorf("specs.%v", err)
		}
		for _, v := range structDataSpec {
			err = v.Validate()
			if err != nil {
				return fmt.Errorf("specs.%v", err)
			}
		}
	case "enum":

		enumSpec := &EnumSpec{
			Specs: bs,
		}
		err = enumSpec.Validate()
	case "date":
		return nil
	default:
		return fmt.Errorf("type err: type is invalid or unsupported for now")
	}
	if err != nil {
		return fmt.Errorf("specs.%v", err)
	}
	return nil
}

// 数值类型
type DigitalSpec struct {
	Max      string
	Min      string
	Step     string
	Unit     string
	UnitName string
}

func (s *DigitalSpec) Validate() error {
	_, err := strconv.ParseUint(s.Max, 10, 32)
	if err != nil {
		return fmt.Errorf("(digital).max err: %v", err)
	}
	_, err = strconv.ParseUint(s.Min, 10, 32)
	if err != nil {
		return fmt.Errorf("(digital).min err: %v", err)
	}
	_, err = strconv.ParseUint(s.Step, 10, 32)
	if err != nil {
		return fmt.Errorf("(digital).step err: %v", err)
	}
	return nil
}

// 数组类型
type ArraySpec struct {
	Size string
	Item *DataType
}

func (s *ArraySpec) Validate() error {
	const (
		maxSize = 512
		MinSize = 1
	)
	size, err := strconv.ParseUint(s.Size, 10, 32)
	if err != nil {
		return fmt.Errorf("(array).size err: %v", err)
	}
	if size > maxSize || size < MinSize {
		err = fmt.Errorf("size out of range [%v, %v]", MinSize, maxSize)
		return fmt.Errorf("(array).size err: %v", err)
	}
	err = s.Item.Validate()
	if err != nil {
		return fmt.Errorf("(array).item.%v", err)
	}
	return nil
}

// 结构体类型
type StructDataSpec struct {
	Identifier string
	Name       string
	DataType   *DataType
}

func (s *StructDataSpec) Validate() error {
	if s.Identifier == "" {
		return fmt.Errorf("(struct).identifier err: identifier is empty")
	}
	if s.Name == "" {
		return fmt.Errorf("(struct).name err: name is empty")
	}
	if s.DataType == nil {
		return fmt.Errorf("(struct).dataType err: dataType is empty")
	}
	if s.DataType.Type == "struct" {
		return fmt.Errorf("(struct).dataType.type  err: struct wrap struct, not support")
	}
	err := s.DataType.Validate()
	if err != nil {
		return fmt.Errorf("(struct).dataType.%v", err)
	}
	return nil
}

// 字符串类型
type TextSpec struct {
	Length string
}

func (s *TextSpec) Validate() error {
	const (
		maxLength = 10240
		MinLength = 1
	)
	length, err := strconv.ParseUint(s.Length, 10, 32)
	if err != nil {
		return fmt.Errorf("%v: %v", s, err)
	}
	if length > maxLength || length < MinLength {
		err = fmt.Errorf("length out of range [%v, %v]", MinLength, maxLength)
		return fmt.Errorf("(text).length err: %v", err)
	}
	return nil
}

// 布尔类型
type BooleanSpec struct {
	FalseValue string `json:"0"`
	TrueValue  string `json:"1"`
}

func (s *BooleanSpec) Validate() error {
	if s.FalseValue == "" {
		return fmt.Errorf("(bool).0  err: value is empty")
	}
	if s.TrueValue == "" {
		return fmt.Errorf("(bool).1, err: ßvalue is empty")
	}
	return nil
}

// 枚举类型
type EnumSpec struct {
	Specs json.RawMessage
}

func (s *EnumSpec) Validate() error {
	var specs map[string]string
	err := json.Unmarshal(s.Specs, &specs)
	if err != nil {
		return fmt.Errorf("(enum) err: %v", err)
	}
	for k, v := range specs {
		if v == "" {
			return fmt.Errorf("(enum).%v err: %v is empty", k, k)
		}
		_, err := strconv.ParseUint(k, 10, 32)
		if err != nil {
			return fmt.Errorf("(enum).%v err: %v is no enum", k, k)
		}
	}
	return nil
}
