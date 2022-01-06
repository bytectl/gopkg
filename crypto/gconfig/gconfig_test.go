package gconfig

import "testing"

func TestDecryptString(t *testing.T) {
	str := "enc(3LPBm9MsFrdXcJAEFX9/JA==)"
	key := []byte("1234567890")
	decrypted, err := DecryptString(str, key)
	if err != nil {
		t.Error(err)
	}
	t.Log(decrypted)
}

func TestEncryptString(t *testing.T) {
	str := "foo"
	key := []byte("1234567890")
	encrypted := EncryptString(str, key)
	t.Log(encrypted)
}
