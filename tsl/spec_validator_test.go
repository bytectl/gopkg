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
					"profile": {
						"productKey": "switch",
						"deviceName": "开关"
					},

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
			`),
			want: nil,
		},
		{
			name: "switch2",
			spec: []byte(`
				{
				"events": [
					{
					"desc": "",
					"identifier": "error",
					"method": "thing.event.error.post",
					"name": "故障",
					"outputData": [
						{
						"dataType": {
							"specs": {
							"max": "1000",
							"min": "0",
							"step": "1",
							"unit": "0",
							"unitName": "无"
							},
							"type": "int"
						},
						"identifier": "errorCode",
						"name": "错误码"
						}
					],
					"required": true,
					"type": "error"
					}
				],
				"profile": {
					"deviceName": "开关",
					"productKey": "xxxx"
				},
				"properties": [
					{
					"accessMode": "rw",
					"dataType": {
						"specs": {
						"0": "关",
						"1": "开"
						},
						"type": "enum"
					},
					"desc": "",
					"identifier": "switch",
					"name": "开关",
					"required": false
					},
					{
					"accessMode": "rw",
					"dataType": {
						"specs": {
						"max": "1440",
						"min": "0",
						"step": "1",
						"unit": "m",
						"unitName": "分"
						},
						"type": "int"
					},
					"desc": "",
					"identifier": "countDown",
					"name": "倒计时",
					"required": false
					},
					{
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
					},
					{
					"accessMode": "rw",
					"dataType": {
						"specs": [
						{
							"dataType": {
							"specs": {
								"length": "16"
							},
							"type": "text"
							},
							"identifier": "ip",
							"name": "ip"
						},
						{
							"dataType": {
							"specs": {
								"length": "16"
							},
							"type": "text"
							},
							"identifier": "mac",
							"name": "mac"
						}
						],
						"type": "struct"
					},
					"desc": "",
					"identifier": "info",
					"name": "信息",
					"required": true
					},
					{
					"accessMode": "rw",
					"dataType": {
						"specs": {
						"item": {
							"type": "text",
							"specs": {
								"length": "16"
							}
						},
						"size": "2"
						},
						"type": "array"
					},
					"desc": "",
					"identifier": "ports",
					"name": "ports",
					"required": true
					}
				],
				"schema": "http://localhost:8000/v1/thing/model/schema",
				"services": [
					{
					"callType": "async",
					"desc": "重启服务",
					"identifier": "reset",
					"inputData": [
						{
						"dataType": {
							"specs": {
							"max": "30",
							"min": "0",
							"step": "1",
							"unit": "s",
							"unitName": "秒"
							},
							"type": "int"
						},
						"identifier": "countDown",
						"name": "倒计时"
						}
					],
					"method": "thing.service.reset",
					"name": "重启",
					"outputData": [
						{
						"dataType": {
							"specs": {
							"0": "失败",
							"1": "成功"
							},
							"type": "bool"
						},
						"identifier": "success",
						"name": "是否成功"
						}
					],
					"required": true
					}
				]
				}			
			`),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateSpec(tt.spec); err != tt.want {
				t.Errorf("ValidateSpec() error = %v, want %v", err, tt.want)
			}
		})
	}

}
