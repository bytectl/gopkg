package tsl

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// 校验接口
type Validator interface {
	ValidateSpec() error
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
			return fmt.Errorf("properties[%d].%v", k, err)
		}
	}
	return nil
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

func (s *Service) ValidateSpec() error {
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

// 属性
type Property struct {
	AccessMode string
	Identifier string
	Name       string
	Desc       string
	Required   bool
	DataType   *DataType
}

func (s *Property) ValidateSpec() error {
	if s.Identifier == "" {
		return fmt.Errorf("identifier err: identifier is empty")
	}
	if s.Name == "" {
		return fmt.Errorf("name  err: name is empty")
	}
	if s.DataType == nil {
		return fmt.Errorf("dataType err: dataType is empty")
	}
	if strings.Compare(s.AccessMode, "r") != 0 && strings.Compare(s.AccessMode, "rw") != 0 {
		return fmt.Errorf("accessMode err: accessMode(%s) is invalid", s.AccessMode)
	}
	err := s.DataType.ValidateSpec()
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

var TypeSpecRegister = map[string]func([]byte) (Validator, error){
	"int":    NewDigitalSpec,
	"float":  NewDigitalSpec,
	"double": NewDigitalSpec,
	"text":   NewTextSpec,
	"enum":   NewEnumSpec,
	"bool":   NewBooleanSpec,
	"array":  NewArraySpec,
	"struct": NewStructSpec,
	"date":   NewEmptySpec,
}

func (s *DataType) ValidateSpec() error {
	var err error
	bs := s.Specs
	if len(bs) == 0 {
		return fmt.Errorf("spec is empty")
	}
	// 查找注册的类型函数
	newValidator, ok := TypeSpecRegister[s.Type]
	if !ok {
		return fmt.Errorf("type err: type is invalid or unsupported for now")
	}
	// 创建相应校验类型
	spec, er := newValidator(bs)
	if er != nil {
		return fmt.Errorf("specs.%v", err)
	}
	err = spec.ValidateSpec()
	if err != nil {
		return fmt.Errorf("specs.%v", err)
	}
	return nil
}

type EmptySpec struct{}

func (s *EmptySpec) ValidateSpec() error {
	fmt.Println("note: empty validateSpec.....")
	return nil
}
func NewEmptySpec(bs []byte) (Validator, error) {
	return &EmptySpec{}, nil
}

// 数值类型
type DigitalSpec struct {
	Max      string
	Min      string
	Step     string
	Unit     string
	UnitName string
}

func NewDigitalSpec(bs []byte) (Validator, error) {
	spec := &DigitalSpec{}
	err := json.Unmarshal(bs, spec)
	if err != nil {
		return nil, fmt.Errorf("(digital) err: %v", err)
	}
	return spec, nil
}
func (s *DigitalSpec) ValidateSpec() error {
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

func NewArraySpec(bs []byte) (Validator, error) {
	spec := &ArraySpec{}
	err := json.Unmarshal(bs, spec)
	if err != nil {
		return nil, fmt.Errorf("(array) err: %v", err)
	}
	return spec, nil
}
func (s *ArraySpec) ValidateSpec() error {
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
	err = s.Item.ValidateSpec()
	if err != nil {
		return fmt.Errorf("(array).item.%v", err)
	}
	return nil
}

// 结构体类型
type StructSpec struct {
	// Identifier Name dataType
	Properties []*Property
}

func NewStructSpec(bs []byte) (Validator, error) {
	var properties []*Property
	err := json.Unmarshal(bs, &properties)
	if err != nil {
		return nil, fmt.Errorf("(struct).%v", err)
	}
	return &StructSpec{Properties: properties}, nil
}

func (s *StructSpec) ValidateSpec() error {
	// 不能直接校验 Properties
	for k, v := range s.Properties {
		if v.Identifier == "" {
			return fmt.Errorf("(struct)[%d].identifier err: identifier is empty", k)
		}
		if v.Name == "" {
			return fmt.Errorf("(struct).name err: name is empty")
		}
		if v.DataType == nil {
			return fmt.Errorf("(struct).dataType err: dataType is empty")
		}
		if v.DataType.Type == "struct" {
			return fmt.Errorf("(struct).dataType.type  err: struct wrap struct, not support")
		}
		err := v.DataType.ValidateSpec()
		if err != nil {
			return fmt.Errorf("(struct).dataType.%v", err)
		}
	}
	return nil
}

// 字符串类型
type TextSpec struct {
	Length string
}

func NewTextSpec(bs []byte) (Validator, error) {
	spec := &TextSpec{}
	err := json.Unmarshal(bs, spec)
	if err != nil {
		return nil, fmt.Errorf("(text) err: %v", err)
	}
	return spec, nil
}

func (s *TextSpec) ValidateSpec() error {
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

func NewBooleanSpec(bs []byte) (Validator, error) {
	spec := &BooleanSpec{}
	err := json.Unmarshal(bs, spec)
	if err != nil {
		return nil, fmt.Errorf("(bool) err: %v", err)
	}
	return spec, nil
}

func (s *BooleanSpec) ValidateSpec() error {
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
	Specs map[string]string
}

func NewEnumSpec(bs []byte) (Validator, error) {
	var specs map[string]string
	err := json.Unmarshal(bs, &specs)
	if err != nil {
		return nil, fmt.Errorf("(enum) err: %v", err)
	}
	enumSpec := &EnumSpec{
		Specs: specs,
	}
	return enumSpec, nil
}

func (s *EnumSpec) ValidateSpec() error {
	for k, v := range s.Specs {
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

func (s *EnumSpec) ValidateValue(interface{}) error {
	return nil
}
