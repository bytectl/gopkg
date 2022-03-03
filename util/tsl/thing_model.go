package tsl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

//	数据类型
type DataType struct {
	Type  string          `json:"type"`
	Specs json.RawMessage `json:"specs"`
}

// 数值类型
type DigitalSpec struct {
	Max      string `json:"max"`
	Min      string `json:"min"`
	Step     string `json:"step"`
	Unit     string `json:"unit"`
	UnitName string `json:"unitName"`
}

// 数组类型
type ArraySpec struct {
	Size  string     `json:"size"`
	Items []DataType `json:"items"`
}

// 结构体类型
type StructSpec struct {
	Identifier string   `json:"identifier"`
	Name       string   `json:"name"`
	DataType   DataType `json:"dataType"`
}

// 属性
type Property struct {
	AccessMode string   `json:"accessMode"`
	Identifier string   `json:"identifier"`
	Name       string   `json:"name"`
	Desc       string   `json:"desc"`
	Required   bool     `json:"required"`
	DataType   DataType `json:"dataType"`
}

type Event struct {
	Identifier string     `json:"identifier"`
	Name       string     `json:"name"`
	Desc       string     `json:"desc"`
	Method     string     `json:"method"`
	Type       string     `json:"type"`
	OutputData []Property `json:"outputData"`
}

type Service struct {
	Identifier string     `json:"identifier"`
	Name       string     `json:"name"`
	Desc       string     `json:"desc"`
	Method     string     `json:"method"`
	CallType   string     `json:"callType"`
	Required   bool       `json:"required"`
	OutputData []Property `json:"outputData"`
	InputData  []Property `json:"inputData"`
}

// 物模型
type ThingModel struct {
	Events     []Event    `json:"events"`
	Services   []Service  `json:"services"`
	Properties []Property `json:"properties"`
}

func NewTingModel(thingModel string) (*ThingModel, error) {
	var tm ThingModel
	err := json.Unmarshal([]byte(thingModel), &tm)
	if err != nil {
		fmt.Println("NewTingModel err:", err)
		return nil, err
	}
	tm.addDefault()
	return &tm, err
}

// 添加默认服务和事件
func (tm *ThingModel) addDefault() {
	// 添加属性上报事件
	tm.Events = append(tm.Events, Event{
		Identifier: "post",
		Name:       "属性上报",
		Desc:       "属性上报",
		Method:     "thing.event.property.post",
		Type:       "info",
		OutputData: tm.Properties,
	})

	// 添加属性设置服务
	serviceSetProperties := []Property{}
	for _, v := range tm.Properties {
		if v.AccessMode == "rw" {
			serviceSetProperties = append(serviceSetProperties, v)
		}
	}
	tm.Services = append(tm.Services, Service{
		Identifier: "set",
		Name:       "属性设置",
		Desc:       "属性设置",
		Method:     "thing.service.property.set",
		CallType:   "sync",
		Required:   true,
		InputData:  serviceSetProperties,
	})
	tm.Services = append(tm.Services, Service{
		Identifier: "get",
		Name:       "属性获取",
		Desc:       "属性获取",
		Method:     "thing.service.property.get",
		CallType:   "sync",
		Required:   true,
		OutputData: tm.Properties,
	})
}

func (tm *ThingModel) propertiesToMap(ps []Property) map[string]Property {
	paramsMap := make(map[string]Property)
	for _, v := range ps {
		paramsMap[v.Identifier] = v
	}
	return paramsMap
}

// 服务校验
func (tm *ThingModel) ServiceValidate(service string, inputParams map[string]interface{}, outputParams map[string]interface{}) (bool, error) {
	result := false
	if service == "" {
		return result, fmt.Errorf("service is empty")
	}
	// if params == nil {
	// 	return result, fmt.Errorf("params is nil")
	// }

	// serviceInfo, ok := tm.mapProperties[service]
	// if !ok {
	// 	return result, fmt.Errorf("service not found")
	// }

	// if serviceInfo.Required {
	// 	for _, v := range serviceInfo.InputData {
	// 		if _, ok := params[v.Identifier]; !ok {
	// 			return result, fmt.Errorf("param %s is required", v.Identifier)
	// 		}
	// 	}
	// }

	return true, nil

}

// 属性校验
func (tm *ThingModel) ValidateProperties(properties string) (bool, error) {
	var (
		paramMap = make(map[string]interface{})
		result   = false
	)

	decoder := json.NewDecoder(bytes.NewReader([]byte(properties)))
	// 使用json number
	decoder.UseNumber()
	if err := decoder.Decode(&paramMap); err != nil {
		return result, err
	}
	mapProperties := tm.propertiesToMap(tm.Properties)
	// 遍历属性模型
	for k, v := range paramMap {
		if v == nil {
			continue
		}
		property, ok := mapProperties[k]
		if !ok {
			return result, fmt.Errorf("property %s is not found in thing model", k)
		}
		preTypeErr := fmt.Errorf("property %s is not a %v,value: %v", k, property.DataType.Type, v)

		switch property.DataType.Type {
		case "long", "date", "int":
			vNumber, ok := v.(json.Number)
			if !ok {
				return result, preTypeErr
			}
			if _, err := vNumber.Int64(); err != nil {
				fmt.Printf("%v\n", err)
				return result, preTypeErr
			}
			fallthrough
		case "float", "double":
			vNumber, ok := v.(json.Number)
			if !ok {
				return result, preTypeErr
			}
			value, err := vNumber.Float64()
			if err != nil {
				fmt.Printf("%v\n", err)
				return result, preTypeErr
			}
			var spec DigitalSpec
			bs, _ := json.Marshal(property.DataType.Specs)
			json.Unmarshal(bs, &spec)
			if spec.Max != "" {
				max, _ := strconv.ParseFloat(spec.Max, 64)
				if value > max {
					return result, fmt.Errorf("property %s is out of max value %v", k, spec.Max)
				}
			}
			if spec.Min != "" {
				min, _ := strconv.ParseFloat(spec.Min, 64)
				if value < min {
					return result, fmt.Errorf("property %s is out of min value %v", k, spec.Min)
				}
			}
		case "enum":
			vNumber, ok := v.(json.Number)
			if !ok {
				return result, preTypeErr
			}
			value, err := vNumber.Int64()
			if err != nil {
				fmt.Printf("%s\n", err)
				return result, preTypeErr
			}
			spec := make(map[string]interface{})
			json.Unmarshal(property.DataType.Specs, &spec)
			if _, ok := spec[strconv.Itoa(int(value))]; !ok {
				return result, preTypeErr
			}
		case "bool":
			if _, ok := v.(bool); !ok {
				return result, preTypeErr
			}
		case "array":
			if _, ok := v.([]interface{}); !ok {
				return result, preTypeErr
			}
		case "struct":
			if _, ok := v.(map[string]interface{}); !ok {
				return result, preTypeErr
			}
		case "text":
			if _, ok := v.(string); !ok {
				return result, preTypeErr
			}
		default:
			return result, fmt.Errorf("property %s type %s is not supported", k, property.DataType.Type)
		}
	}
	return true, nil
}
