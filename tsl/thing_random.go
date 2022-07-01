package tsl

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// 随机生成物模型实体数据
func (s *Thing) RandomAll() ([]byte, error) {
	rand.Seed(time.Now().UnixNano())
	var m struct {
		Events   []*EntityRequest `json:"events"`
		Services []*EntityRequest `json:"services"`
	}

	s.init() // initialize
	for _, v := range s.Value.Services {
		e, err := v.Random(true)
		if err != nil {
			return nil, err
		}
		m.Services = append(m.Services, &EntityRequest{
			ID:        e.ID,
			Version:   e.Version,
			Method:    e.Method,
			Params:    e.Params,
			Timestamp: e.Timestamp,
		})
	}
	for _, v := range s.Value.Events {
		e, err := v.Random(true)
		if err != nil {
			return nil, err
		}
		m.Events = append(m.Events, &EntityRequest{
			ID:        e.ID,
			Version:   e.Version,
			Method:    e.Method,
			Params:    e.Params,
			Timestamp: e.Timestamp,
		})
	}
	bs, err := json.MarshalIndent(m, "", "  ")
	return bs, err
}

func (s *Thing) Random(method string, generateAllProperty bool) ([]byte, error) {
	var entity *ThingEntity
	rand.Seed(time.Now().UnixNano())
	s.init() // initialize
	tmethod, err := NewThingMethod(method)
	if err != nil {
		return nil, fmt.Errorf("method: (%s) no found", method)
	}
	if tmethod.IsService() {
		tm := s.Value.Services[tmethod.Action]
		if tm == nil {
			return nil, fmt.Errorf("service.%s no found", tmethod.Action)
		}
		entity, err = tm.Random(generateAllProperty)
		if err != nil {
			return nil, err
		}
	} else {
		tm := s.Value.Events[tmethod.Action]
		if tm == nil {
			return nil, fmt.Errorf("event.%s no found", tmethod.Action)
		}
		entity, err = tm.Random(generateAllProperty)
		if err != nil {
			return nil, err
		}
	}

	bs, err := json.MarshalIndent(entity, "", "  ")
	return bs, err
}

func (s *Event) Random(generateAllProperty bool) (*ThingEntity, error) {
	s.init() // initialize
	inputData := make(map[string]interface{})
	inputBytes, _ := json.Marshal(inputData)

	tmethod, err := NewThingMethod(s.Method)
	if err != nil {
		return nil, fmt.Errorf("Event.Random, err: %v", err)
	}
	outputData := propertyRandomValueToMap(s.OutputData)
	if tmethod.IsProperty && generateAllProperty == false {
		// 随机生成属性和属性值propertyRandomAndRandomValueToMap
		outputData = propertyRandomAndRandomValueToMap(s.OutputData)
	}
	outputBytes, err := json.Marshal(outputData)
	if err != nil {
		return nil, fmt.Errorf("Event.Random, err: %v", err)
	}
	return &ThingEntity{
		ID:        fmt.Sprintf("%d", rand.Int31()),
		Version:   "1.0",
		Timestamp: time.Now().UnixMilli(),
		Params:    outputBytes,
		Data:      inputBytes,
		Method:    s.Method,
	}, nil
}

func (s *Service) Random(generateAllProperty bool) (*ThingEntity, error) {
	s.init() // initialize
	inputData := propertyRandomValueToMap(s.InputData)
	outputData := propertyRandomValueToMap(s.OutputData)
	tmethod, err := NewThingMethod(s.Method)
	if err != nil {

		return nil, fmt.Errorf("Service.Random, err: %v", err)
	}
	if tmethod.IsProperty && tmethod.IsSet && generateAllProperty == false {
		// 随机生成属性和属性值propertyRandomAndRandomValueToMap
		inputData = propertyRandomAndRandomValueToMap(s.InputData)
	}
	inputBytes, err := json.Marshal(inputData)
	if err != nil {

		return nil, fmt.Errorf("Service.Random, err: %v", err)
	}
	outputBytes, err := json.Marshal(outputData)
	if err != nil {
		return nil, fmt.Errorf("Service.Random, err: %v", err)
	}
	return &ThingEntity{
		ID:        fmt.Sprintf("%d", rand.Int31()),
		Version:   "1.0",
		Timestamp: time.Now().UnixMilli(),
		Params:    inputBytes,
		Data:      outputBytes,
		Method:    s.Method,
	}, nil
}
