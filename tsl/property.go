package tsl

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

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
func (s *Property) ToEntityString() string {
	specs := []string{
		s.DataType.Type,
		s.Name,
	}
	if s.Desc != "" {
		specs = append(specs, s.Desc)
	}
	if s.DataType.Type == "struct" || s.DataType.Type == "array" {
		// 直接返回 json 字符串
		return s.DataType.ToEntityString()
	}
	specs = append(specs, s.DataType.ToEntityString())
	return strings.Join(specs, ",")
}

func (s *Property) Random() interface{} {
	return s.DataType.Random()
}

func propertiesToMap(ps []*Property) map[string]*Property {
	paramsMap := make(map[string]*Property)
	for _, v := range ps {
		paramsMap[v.Identifier] = v
	}
	return paramsMap
}

// 属性列表转换为map
func propertyToEntityMap(p []*Property) map[string]interface{} {
	m := map[string]interface{}{}

	for _, v := range p {
		str := v.ToEntityString()
		m[v.Identifier] = str
		if v.DataType.Type == "struct" {
			tm := map[string]interface{}{}
			err := json.Unmarshal([]byte(str), &tm)
			if err != nil {
				fmt.Printf("%v Unmarshal err: %v: \n", str, err)
			}
			m[v.Identifier] = tm
		} else if v.DataType.Type == "array" {
			tm := []interface{}{}
			err := json.Unmarshal([]byte(str), &tm)
			if err != nil {
				fmt.Printf("%v Unmarshal err: %v: \n", str, err)
			}
			m[v.Identifier] = tm
		}
	}
	return m
}

// 随机生成属性值转map
func propertyRandomValueToMap(p []*Property) map[string]interface{} {
	m := map[string]interface{}{}
	for _, v := range p {
		m[v.Identifier] = v.Random()
	}
	return m
}

// 随机生成属性和属性值转map
func propertyRandomAndRandomValueToMap(p []*Property) map[string]interface{} {

	// 属性数量
	count := len(p)
	seq := p
	if count != 0 {
		seq = make([]*Property, 0, count)
		// copy slice
		for i := 0; i < count; i++ {
			seq = append(seq, p[i])
		}
		// // 随机排序
		for i := 0; i < count; i++ {
			rand.Seed(time.Now().UnixNano())
			j := rand.Intn(count)
			seq[i], seq[j] = seq[j], seq[i]
		}
		n := rand.Intn(count)
		// 取n个属性
		seq = seq[:n+1]
	}

	m := map[string]interface{}{}
	for _, v := range seq {
		m[v.Identifier] = v.Random()
	}
	return m
}
