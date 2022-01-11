package gconfig

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gogf/gf/crypto/gaes"
)

const (
	EncPrefix = "enc("
	EncSuffix = ")"
)

func DecryptConfigMap(config map[string]interface{}, key []byte) {
	for k, value := range config {
		switch v1 := value.(type) {
		case string:
			source, err := DecryptString(v1, key)
			if err != nil {
				fmt.Println(k, ":", v1, ",decrypt error:", err)
				continue
			}
			config[k] = source
		case map[string]interface{}:
			DecryptConfigMap(v1, key)
		case []interface{}:
			for i, item := range v1 {
				switch v2 := item.(type) {
				case string:
					source, err := DecryptString(v2, key)
					if err != nil {
						fmt.Println(k, ":", v2, ",decrypt error:", err)
						continue
					}
					v1[i] = source
				case map[string]interface{}:
					DecryptConfigMap(v2, key)
				}
			}
		}
	}
}

func EncryptString(source string, key []byte) string {
	key = gaes.PKCS5Padding(key, 16)
	encrypted, err := gaes.Encrypt([]byte(source), key)
	if err != nil {
		fmt.Println(source, ",encrypt error:", err)
		return source
	}
	return EncPrefix + base64.StdEncoding.EncodeToString(encrypted) + EncSuffix

}
func DecryptString(encrypted string, key []byte) (string, error) {

	if !strings.HasPrefix(encrypted, EncPrefix) {
		return encrypted, nil
	}
	key = gaes.PKCS5Padding(key, 16)
	encrypted = strings.TrimPrefix(encrypted, EncPrefix)
	encrypted = strings.TrimSuffix(encrypted, EncSuffix)
	decoded, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return encrypted, err
	}
	decrypted, err := gaes.Decrypt(decoded, key)
	if err != nil {
		return encrypted, err
	}
	return string(decrypted), nil
}
