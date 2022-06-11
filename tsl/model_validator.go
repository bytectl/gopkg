package tsl

import (
	"encoding/json"
	"fmt"
)

// 校验物模型正确性
func ValidateThingSpec(spec []byte) error {
	var thing Thing
	err := json.Unmarshal(spec, &thing)
	if err != nil {
		fmt.Println("ValidateThing err:", err)
		return fmt.Errorf("ValidateThing err: %s", err)
	}
	err = thing.Validate()
	if err != nil {
		fmt.Println("ValidateThing err:", err)
		return err
	}
	return err
}
