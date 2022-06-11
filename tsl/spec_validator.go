package tsl

import (
	"encoding/json"
	"fmt"
)

// 校验物模型正确性
func ValidateSpec(spec []byte) error {
	var thing Thing
	err := json.Unmarshal(spec, &thing)
	if err != nil {
		fmt.Println("ValidateSpec err:", err)
		return fmt.Errorf("ValidateSpec err: %s", err)
	}
	err = thing.ValidateSpec()
	if err != nil {
		fmt.Println("ValidateSpec err:", err)
		return err
	}
	return err
}
