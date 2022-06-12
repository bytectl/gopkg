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
		return fmt.Errorf("ValidateSpec err: %s", err)
	}
	err = thing.ValidateSpec()
	if err != nil {
		return err
	}
	return err
}
