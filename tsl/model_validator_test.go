package tsl

import (
	"testing"
)

func TestValidate(t *testing.T) {
	var tests = []struct {
		name string
		spec []byte
		want error
	}{
		{
			name: "empty",
			spec: []byte(`
				{
					"schema": "http://localhost:8000/v1/thing/model/schema",
					"profile": {
						"productKey": "empty",
						"deviceName": "无"
					},
					"services": [],
					"properties": [],
					"events": []
				}
			`),
			want: nil,
		},
		{
			name: "switch",
			spec: []byte(`
				{
					"schema": "http://localhost:8000/v1/thing/model/schema",
					"profile": {
						"productKey": "switch",
						"deviceName": "开关"
					},

					"properties": [{
						"identifier": "switch",
						"dataType": {
							"specs": {
								"0": "关",
								"1": "开",
								"test": "测试"
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
			`),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateThingSpec(tt.spec); err != tt.want {
				t.Errorf("ValidateThingSpec() error = %v, want %v", err, tt.want)
			}
		})
	}

}
