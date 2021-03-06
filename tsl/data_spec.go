package tsl

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/util/grand"
)

// 校验接口
type Validator interface {
	// 校验物模型定义规范
	ValidateSpec() error
	// 校验物模型实例值
	ValidateValue(value interface{}) error
	// 转换为物模型实例
	ToEntityString() string
	// 随机值生成
	Random() interface{}
}

// 数据类型
type DataType struct {
	// 数据类型
	Type string
	// 数据类型定义
	Specs json.RawMessage
	Value struct {
		Specs Validator
	}
}

// 数据类型注册表
var TypeSpecRegister = map[string]func([]byte) (Validator, error){
	"int":    NewDigitalSpec,
	"long":   NewDigitalSpec,
	"float":  NewFloatSpec,
	"double": NewFloatSpec,
	"text":   NewTextSpec,
	"enum":   NewEnumSpec,
	"bool":   NewBooleanSpec,
	"array":  NewArraySpec,
	"struct": NewStructSpec,
	"date":   NewDateSpec,
}

func (s *DataType) init() error {
	if s.Value.Specs != nil {
		return nil
	}
	bs := s.Specs
	if len(bs) == 0 {
		return fmt.Errorf("spec is empty")
	}
	if TypeSpecRegister[s.Type] == nil {
		return fmt.Errorf("type %s is not supported", s.Type)
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
	if spec == nil {
		return fmt.Errorf("specs is empty")
	}
	s.Value.Specs = spec
	return nil
}

func (s *DataType) ValidateSpec() error {
	err := s.init() // 初始化
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	err = s.Value.Specs.ValidateSpec()
	if err != nil {
		return fmt.Errorf("specs.%v", err)
	}
	return nil
}

func (s *DataType) ValidateValue(value interface{}) error {
	err := s.init() // 初始化
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	err = s.Value.Specs.ValidateValue(value)
	if err != nil {
		return fmt.Errorf("%v.%v", s.Type, err)
	}
	return nil
}

func (s *DataType) ToEntityString() string {
	err := s.init() // 初始化
	if err != nil {
		fmt.Println("DataType.ToEntityString,", err)
		return ""
	}
	return s.Value.Specs.ToEntityString()
}
func (s *DataType) Random() interface{} {
	err := s.init() // 初始化
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return s.Value.Specs.Random()
}

func (s *DataType) GenerateGoType() string {
	err := s.init() // 初始化
	if err != nil {
		fmt.Printf("GenerateGoType:%v", err)
		return "unknown"
	}
	switch s.Type {
	case "int":
		return "int"
	case "long":
		return "int64"
	case "float":
		return "float64"
	case "double":
		return "float64"
	case "text":
		return "string"
	case "enum":
		return "int"
	case "bool":
		return "int8"
	case "array":
		return "[]interface{}"
	case "struct":
		return "map[string]interface{}"
	case "date":
		return "int64"
	default:
		return "unknown"
	}
}

// 默认值
func (s *DataType) DefaultValueString() string {
	err := s.init() // 初始化
	if err != nil {
		fmt.Printf("GenerateGoType:%v", err)
		return "nil"
	}
	switch s.Type {
	case "int":
		return "0"
	case "long":
		return "0"
	case "float":
		return "0.0"
	case "double":
		return "0.0"
	case "text":
		return `""`
	case "enum":
		return "0"
	case "bool":
		return "0"
	case "array":
		return "nil"
	case "struct":
		return "nil"
	case "date":
		return "0"
	}
	return "nil"
}

type DateSpec struct{}

func NewDateSpec(bs []byte) (Validator, error) {
	return &DateSpec{}, nil
}
func (s *DateSpec) ValidateSpec() error {
	fmt.Println("note: empty validateSpec.....")
	return nil
}
func (s *DateSpec) ValidateValue(value interface{}) error {
	_, ok := value.(string)
	if !ok {
		return fmt.Errorf("(date).value err: %v is not string", value)
	}
	return nil
}
func (s *DateSpec) ToEntityString() string {
	return "unix时间戳(ms)"
}
func (s *DateSpec) Random() interface{} {
	return time.Now().UnixMilli()
}

// 数值类型
type DigitalSpec struct {
	Max      string
	Min      string
	Step     string
	Unit     string
	UnitName string
	Value    struct {
		Max  int64
		Min  int64
		Step uint64
	}
}

func NewDigitalSpec(bs []byte) (Validator, error) {
	spec := &DigitalSpec{}
	err := json.Unmarshal(bs, spec)
	if err != nil {
		return nil, fmt.Errorf("(digital) err: %v", err)
	}
	max, err := strconv.ParseInt(spec.Max, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("(digital).max err: %v", err)
	}
	min, err := strconv.ParseInt(spec.Min, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("(digital).min err: %v", err)
	}
	step := uint64(0)
	if len(spec.Step) != 0 {
		step, err = strconv.ParseUint(spec.Step, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("(digital).step err: %v", err)
		}
	}
	spec.Value.Max = max
	spec.Value.Min = min
	spec.Value.Step = step
	return spec, nil
}
func (s *DigitalSpec) ValidateSpec() error {
	if s.Value.Min > s.Value.Max {
		return fmt.Errorf("(float).min err: min is larger than max")
	}
	if s.Value.Step != 0 && s.Value.Step > uint64(s.Value.Max-s.Value.Min) {
		return fmt.Errorf("(digital).step err: step is too large")
	}
	return nil
}
func (s *DigitalSpec) ValidateValue(value interface{}) error {
	vNumber, ok := value.(json.Number)
	if !ok {
		return fmt.Errorf("(digital).value err: %v type %T is not json.Number", value, value)
	}
	int64Value, err := vNumber.Int64()
	if err != nil {
		return fmt.Errorf("(digital).value err: %v", err)
	}
	if int64Value < s.Value.Min || int64Value > s.Value.Max {
		return fmt.Errorf("(digital).value err: value(%v) is out of range [%v, %v]", int64Value, s.Value.Min, s.Value.Max)
	}
	return nil
}

func (s *DigitalSpec) ToEntityString() string {
	spec := fmt.Sprintf("range: %v-%v %v(%v),step: %v", s.Value.Min, s.Value.Max, s.UnitName, s.Unit, s.Value.Step)
	return spec
}

func (s *DigitalSpec) Random() interface{} {
	return int64(rand.Intn(int(s.Value.Max-s.Value.Min+1))) + s.Value.Min
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
	step := float64(0)
	if len(spec.Step) != 0 {
		step, err = strconv.ParseFloat(spec.Step, 64)
		if err != nil {
			return nil, fmt.Errorf("(float).step err: %v", err)
		}
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
	vNumber, ok := value.(json.Number)
	if !ok {
		return fmt.Errorf("(float).value err: %v type %T is not json.Number", value, value)
	}
	floatValue, err := vNumber.Float64()
	if err != nil {
		return fmt.Errorf("(float).value err: %v", err)
	}
	if floatValue > s.Value.Max || floatValue < s.Value.Min {
		return fmt.Errorf("(float) err: value(%v) is out of range [%v, %v]", floatValue, s.Value.Min, s.Value.Max)
	}
	return nil
}

func (s *FloatSpec) ToEntityString() string {
	spec := fmt.Sprintf("range: %v-%v(unit:%v),step: %v", s.Value.Min, s.Value.Max, s.Unit, s.Value.Step)
	return spec
}

func (s *FloatSpec) Random() interface{} {
	return rand.Float64()*(s.Value.Max-s.Value.Min) + s.Value.Min
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
		err := fmt.Errorf("length(%v) out of range [%v, %v]", s.Value.Length, MinLength, maxLength)
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

func (s *TextSpec) ToEntityString() string {
	spec := fmt.Sprintf("max-length: %v", s.Value.Length)
	return spec
}

func (s *TextSpec) Random() interface{} {
	n := rand.Intn(s.Value.Length + 1)
	return grand.Letters(n)
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
	vNumber, ok := value.(json.Number)
	if !ok {
		return fmt.Errorf("(bool).value err: %v type %T is not json.Number", value, value)
	}
	boolValue, err := vNumber.Int64()
	if err != nil {
		return fmt.Errorf("(bool).value err: %v", err)
	}
	if boolValue != 0 && boolValue != 1 {
		return fmt.Errorf("(bool).value err: %v  is not bool", value)
	}
	return nil
}

func (s *BooleanSpec) ToEntityString() string {
	spec := fmt.Sprintf("0-%v,1-%v", s.FalseValue, s.TrueValue)
	return spec
}

func (s *BooleanSpec) Random() interface{} {
	return rand.Intn(2)
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
	vNumber, ok := value.(json.Number)
	if !ok {
		return fmt.Errorf("(enum).value err: %v type %T is not json.Number", value, value)
	}
	enumValue, err := vNumber.Int64()
	if err != nil {
		return fmt.Errorf("(enum).value err: %v", err)
	}
	if _, ok := s.Value.Specs[int(enumValue)]; !ok {
		return fmt.Errorf("(enum).value err: %+v is not defined enum", value)
	}
	return nil
}

func (s *EnumSpec) ToEntityString() string {
	specs := []string{}
	for k, v := range s.Specs {
		specs = append(specs, fmt.Sprintf("%v-%v", k, v))
	}
	return strings.Join(specs, ",")
}
func (s *EnumSpec) Random() interface{} {
	n := rand.Intn(len(s.Specs))
	for k := range s.Specs {
		if 0 == n {
			i, err := strconv.Atoi(k)
			if err != nil {
				return 0
			}
			return i
		}
		n--
	}
	return 0
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
		err := fmt.Errorf("size(%v) out of range [%v, %v]", s.Value.Size, MinSize, maxSize)
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

func (s *ArraySpec) Random() interface{} {
	arrayValue := []interface{}{}
	for i := 0; i < s.Value.Size; i++ {
		arrayValue = append(arrayValue, s.Item.Random())
	}
	return arrayValue
}

func (s *ArraySpec) ToEntityString() string {
	var items []interface{}
	str := fmt.Sprintf("%v,%v,size:%v", s.Item.Type, s.Item.ToEntityString(), s.Value.Size)
	if s.Item.Type == "struct" {

		vm := map[string]interface{}{}
		err := json.Unmarshal([]byte(s.Item.ToEntityString()), &vm)
		if err != nil {
			fmt.Println("Unmarshal item,err: ", err)
		}
		items = append(items, vm)
	} else {
		items = append(items, str)
	}
	bs, _ := json.Marshal(items)
	return string(bs)
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
	return structSpec, nil
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

func (s *StructSpec) ToEntityString() string {
	m := propertyToEntityMap(s.Properties)
	bs, _ := json.Marshal(m)
	return string(bs)
}

func (s *StructSpec) Random() interface{} {
	m := map[string]interface{}{}
	for _, v := range s.Properties {
		m[v.Identifier] = v.Random()
	}
	return m
}
