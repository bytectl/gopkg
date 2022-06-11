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
	ValidateValue(value interface{}) error
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
			return fmt.Errorf("properties[%d,(%s)].%v", k, property.Identifier, err)
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
	Value      struct {
		OutputData map[string]*Property
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
func (s *Event) ValidateValue(value interface{}) error {
	if value == nil {
		return fmt.Errorf("value err: value is empty")
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
	Value      struct {
		InputData  map[string]*Property
		OutputData map[string]*Property
	}
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

func (s *Service) ValidateValue(value interface{}) error {

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
	if s.AccessMode != "" && strings.Compare(s.AccessMode, "r") != 0 && strings.Compare(s.AccessMode, "rw") != 0 {
		return fmt.Errorf("accessMode err: accessMode(%s) is invalid", s.AccessMode)
	}
	err := s.DataType.ValidateSpec()
	if err != nil {
		return fmt.Errorf("dataType.%v", err)
	}
	return nil
}

func (s *Property) ValidateValue(value interface{}) error {
	return s.DataType.ValidateValue(value)
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
	"float":  NewFloatSpec,
	"double": NewFloatSpec,
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
	spec, err := newValidator(bs)
	if err != nil {
		return fmt.Errorf("specs.%v", err)
	}
	err = spec.ValidateSpec()
	if err != nil {
		return fmt.Errorf("specs.%v", err)
	}
	return nil
}
func (s *DataType) ValidateValue(value interface{}) error {
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
	spec, err := newValidator(bs)
	if err != nil {
		return fmt.Errorf("%v.%v", s.Type, err)
	}
	err = spec.ValidateValue(value)
	if err != nil {
		return fmt.Errorf("%v.%v", s.Type, err)
	}
	return nil
}

type EmptySpec struct{}

func NewEmptySpec(bs []byte) (Validator, error) {
	return &EmptySpec{}, nil
}
func (s *EmptySpec) ValidateSpec() error {
	fmt.Println("note: empty validateSpec.....")
	return nil
}
func (s *EmptySpec) ValidateValue(value interface{}) error {
	fmt.Println("note: empty validateValue.....")
	return nil
}

// 数值类型
type DigitalSpec struct {
	Max      string
	Min      string
	Step     string
	Unit     string
	UnitName string
	Value    struct {
		Max  int
		Min  int
		Step int
	}
}

func NewDigitalSpec(bs []byte) (Validator, error) {
	spec := &DigitalSpec{}
	err := json.Unmarshal(bs, spec)
	if err != nil {
		return nil, fmt.Errorf("(digital) err: %v", err)
	}
	max, err := strconv.ParseUint(spec.Max, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("(digital).max err: %v", err)
	}
	min, err := strconv.ParseUint(spec.Min, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("(digital).min err: %v", err)
	}
	step, err := strconv.ParseUint(spec.Step, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("(digital).step err: %v", err)
	}
	spec.Value.Max = int(max)
	spec.Value.Min = int(min)
	spec.Value.Step = int(step)
	return spec, nil
}
func (s *DigitalSpec) ValidateSpec() error {
	if s.Value.Min > s.Value.Max {
		return fmt.Errorf("(float).min err: min is larger than max")
	}
	if s.Value.Step > s.Value.Max-s.Value.Min {
		return fmt.Errorf("(digital).step err: step is too large")
	}
	return nil
}
func (s *DigitalSpec) ValidateValue(value interface{}) error {
	intValue, ok := value.(int)
	if !ok {
		return fmt.Errorf("(digital).value err: %v is not int", value)
	}
	if intValue < s.Value.Min || intValue > s.Value.Max {
		return fmt.Errorf("(digital).value err: value is out of range [%v, %v]", s.Value.Min, s.Value.Max)
	}
	return nil
}

// 数值类型
type FloatSpec struct {
	Max      string
	Min      string
	Step     string
	Unit     string
	UnitName string
	Value    struct {
		Max  float64
		Min  float64
		Step float64
	}
}

func NewFloatSpec(bs []byte) (Validator, error) {
	spec := &FloatSpec{}
	err := json.Unmarshal(bs, spec)
	if err != nil {
		return nil, fmt.Errorf("(float) err: %v", err)
	}
	max, err := strconv.ParseFloat(spec.Max, 64)
	if err != nil {
		return nil, fmt.Errorf("(float).max err: %v", err)
	}
	min, err := strconv.ParseFloat(spec.Min, 64)
	if err != nil {
		return nil, fmt.Errorf("(float).min err: %v", err)
	}
	step, err := strconv.ParseFloat(spec.Step, 64)
	if err != nil {
		return nil, fmt.Errorf("(float).step err: %v", err)
	}
	spec.Value.Max = max
	spec.Value.Min = min
	spec.Value.Step = step
	return spec, nil
}
func (s *FloatSpec) ValidateSpec() error {
	if s.Value.Min > s.Value.Max {
		return fmt.Errorf("(float).min err: min is larger than max")
	}
	if s.Value.Step > s.Value.Max-s.Value.Min {
		return fmt.Errorf("(float).step err: step is greater than max")
	}
	return nil
}
func (s *FloatSpec) ValidateValue(value interface{}) error {
	floatValue, ok := value.(float64)
	if !ok {
		return fmt.Errorf("(float) err: value is not float64")
	}
	if floatValue > s.Value.Max || floatValue < s.Value.Min {
		return fmt.Errorf("(float) err: value is out of range [%v, %v]", s.Value.Min, s.Value.Max)
	}
	return nil
}

// 字符串类型
type TextSpec struct {
	Length string
	Value  struct {
		Length int
	}
}

func NewTextSpec(bs []byte) (Validator, error) {
	spec := &TextSpec{}
	err := json.Unmarshal(bs, spec)
	if err != nil {
		return nil, fmt.Errorf("(text) err: %v", err)
	}
	length, err := strconv.ParseUint(spec.Length, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("(text).length err: %v", err)
	}
	spec.Value.Length = int(length)
	return spec, nil
}

func (s *TextSpec) ValidateSpec() error {
	const (
		maxLength = 10240
		MinLength = 1
	)
	if s.Value.Length > maxLength || s.Value.Length < MinLength {
		err := fmt.Errorf("length out of range [%v, %v]", MinLength, maxLength)
		return fmt.Errorf("(text).length err: %v", err)
	}
	return nil
}

func (s *TextSpec) ValidateValue(value interface{}) error {
	stringValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("(text).value err: %v is not string", value)
	}
	if len(stringValue) > s.Value.Length {
		return fmt.Errorf("(text).value err: %v is too long then %d", value, s.Value.Length)
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

func (s *BooleanSpec) ValidateValue(value interface{}) error {
	_, ok := value.(bool)
	if !ok {
		return fmt.Errorf("(bool).value err: %v  is not bool", value)
	}
	return nil
}

// 枚举类型
type EnumSpec struct {
	Specs map[string]string
	Value struct {
		Specs map[int]string
	}
}

func NewEnumSpec(bs []byte) (Validator, error) {
	var specs map[string]string
	err := json.Unmarshal(bs, &specs)
	if err != nil {
		return nil, fmt.Errorf("(enum) err: %v", err)
	}
	var vspecs map[int]string = make(map[int]string)
	for k, v := range specs {
		if v == "" {
			return nil, fmt.Errorf("(enum).%v err: %v is empty", k, k)
		}
		ivalue, err := strconv.ParseUint(k, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("(enum).%v err: %v is no enum", k, k)
		}
		vspecs[int(ivalue)] = v
	}
	enumSpec := &EnumSpec{
		Specs: specs,
	}
	enumSpec.Value.Specs = vspecs
	return enumSpec, nil
}

func (s *EnumSpec) ValidateSpec() error {
	return nil
}

func (s *EnumSpec) ValidateValue(value interface{}) error {
	intValue, ok := value.(int)
	if !ok {
		return fmt.Errorf("(enum).value err: %v is not int", value)
	}
	if _, ok := s.Value.Specs[intValue]; !ok {
		return fmt.Errorf("(enum).value err: %v is not enum", value)
	}
	return nil
}

// 数组类型
type ArraySpec struct {
	Size  string
	Item  *DataType
	Value struct {
		Size int
	}
}

func NewArraySpec(bs []byte) (Validator, error) {
	spec := &ArraySpec{}
	err := json.Unmarshal(bs, spec)
	if err != nil {
		return nil, fmt.Errorf("(array) err: %v", err)
	}
	size, err := strconv.ParseUint(spec.Size, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("(array).size err: %v", err)
	}
	spec.Value.Size = int(size)
	return spec, nil
}
func (s *ArraySpec) ValidateSpec() error {
	const (
		maxSize = 512
		MinSize = 1
	)
	if s.Value.Size > maxSize || s.Value.Size < MinSize {
		err := fmt.Errorf("size out of range [%v, %v]", MinSize, maxSize)
		return fmt.Errorf("(array).size err: %v", err)
	}
	err := s.Item.ValidateSpec()
	if err != nil {
		return fmt.Errorf("(array).item.%v", err)
	}
	return nil
}
func (s *ArraySpec) ValidateValue(value interface{}) error {
	arrayValue, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("(array).value err: %v is not array", value)
	}
	if len(arrayValue) > int(s.Value.Size) {
		return fmt.Errorf("(array).value err: %v is too long then %d", value, s.Value.Size)
	}
	for _, v := range arrayValue {
		err := s.Item.ValidateValue(v)
		if err != nil {
			return fmt.Errorf("(array).value err: %v", err)
		}
	}
	return nil
}

// 结构体类型
type StructSpec struct {
	// Identifier Name dataType
	Properties []*Property
	Value      struct {
		Properties map[string]*Property
	}
}

func NewStructSpec(bs []byte) (Validator, error) {
	var properties []*Property
	err := json.Unmarshal(bs, &properties)
	if err != nil {
		return nil, fmt.Errorf("(struct).%v", err)
	}
	structSpec := &StructSpec{
		Properties: properties,
	}
	structSpec.Value.Properties = propertiesToMap(structSpec.Properties)
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

func (s *StructSpec) ValidateValue(value interface{}) error {
	mapValue, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("(struct).value err: %v is not map", value)
	}
	for k, v := range mapValue {
		property, ok := s.Value.Properties[k]
		if !ok {
			return fmt.Errorf("(struct).value err: %v is not found", k)
		}
		err := property.DataType.ValidateValue(v)
		if err != nil {
			return fmt.Errorf("(struct).value err: %v", err)
		}
	}
	return nil
}

func propertiesToMap(ps []*Property) map[string]*Property {
	paramsMap := make(map[string]*Property)
	for _, v := range ps {
		paramsMap[v.Identifier] = v
	}
	return paramsMap
}
