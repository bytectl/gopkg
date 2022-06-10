package tsl

import "encoding/json"

// 物模型实体
type ThingEntity struct {
	ID        string          `json:"id"`      // 消息ID，String类型的数字，取值范围0~4294967295，且每个消息ID在当前设备中具有唯一性。
	Version   string          `json:"version"` // 协议版本号，目前协议版本号唯一取值为1.0。
	Params    json.RawMessage `json:"params,omitempty"`
	Method    string          `json:"method"`
	Timestamp string          `json:"timestamp"`
	Data      json.RawMessage `json:"data,omitempty"`
	// Sys        interface{} `json:"sys,omitempty"` //	扩展功能的参数，其下包含各功能字段。
}
