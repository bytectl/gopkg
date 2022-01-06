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
		if s, ok := value.(string); ok {
			if strings.HasPrefix(s, EncPrefix) {
				config[k], _ = DecryptString(s, key)
			}
		} else if m, ok := value.(map[string]interface{}); ok {
			DecryptConfigMap(m, key)
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
