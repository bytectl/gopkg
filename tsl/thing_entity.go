package tsl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

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
	if thingEntity.Timestamp <= 0 {
		err = fmt.Errorf("thingEntity.timestamp <=0 , invalid")
	}
	return err
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
func validateEntityParams(specData map[string]*Property, data []byte) error {
	var err error
	if data == nil || len(string(data)) == 0 || strings.Compare(string(data), "{}") == 0 {
		return nil
	}
	paramMap := make(map[string]interface{})
	decoder := json.NewDecoder(bytes.NewReader(data))
	// 使用json number
	decoder.UseNumber()
	if err := decoder.Decode(&paramMap); err != nil {
		return err
	}
	if len(specData) == 0 && len(paramMap) > 0 {
		return fmt.Errorf("specData is empty, but params is not empty")
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
