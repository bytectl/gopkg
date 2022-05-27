package tsl

import (
	"encoding/json"
	"testing"
)

func TestValidateThingModel(t *testing.T) {
	model := `{
    "properties": [{
        "identifier": "switch",
        "dataType": {
            "specs": {
                "0": "关",
                "1": "开"
            },
            "type": "enum"
        },
        "name": "开关",
        "accessMode": "rw",
        "required": false,
        "desc": ""
    }, {
        "identifier": "countDown",
        "dataType": {
            "specs": {
                "unit": "m",
                "min": "0",
                "unitName": "分",
                "max": "1440",
                "step": "1"
            },
            "type": "int"
        },
        "name": "倒计时",
        "accessMode": "rw",
        "required": false,
        "desc": ""
    }, {
        "accessMode": "rw",
        "dataType": {
            "specs": {
                "length": "256"
            },
            "type": "text"
        },
        "desc": "",
        "identifier": "order",
        "name": "定时",
        "required": false
    }],
    "services": [],
    "events": []
}
`
	entity := `
	{
      "id": "123",
      "method": "thing.event.property.post",
      "params": {
        "countDown": 0,
        "order": "test",
        "switch": 0
      },
      "timestamp": "123",
      "version": "1.0"
    }
	`
	tm, err := NewThingModel(model)
	if err != nil {
		t.Error(err)
	}
	vm := tm.ToValidateModel()
	en := ThingEntity{}
	err = json.Unmarshal([]byte(entity), &en)
	if err != nil {
		t.Error(err)
	}
	p, _ := json.Marshal(en.Params)
	_, err = vm.EventValidate("post", string(p))
	if err != nil {
		t.Error(err)
	}

}
