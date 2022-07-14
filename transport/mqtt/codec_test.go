package mqtt

import (
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
)

func TestDefaultRequestDecoder(t *testing.T) {
	v1 := &struct {
		A string `json:"a"`
		B int64  `json:"b"`
	}{}
	err1 := DefaultRequestDecoder([]byte("{\"a\":\"1\", \"b\": 2}"), &v1)
	if err1 != nil {
		t.Errorf("expected no error, got %v", err1)
	}
	if !reflect.DeepEqual("1", v1.A) {
		t.Errorf("expected %v, got %v", "1", v1.A)
	}
	if !reflect.DeepEqual(int64(2), v1.B) {
		t.Errorf("expected %v, got %v", 2, v1.B)
	}
}

type dataWithStatusCode struct {
	A string `json:"a"`
	B int64  `json:"b"`
}

func TestDefaultResponseEncoder(t *testing.T) {

	v1 := &dataWithStatusCode{A: "1", B: 2}

	err := DefaultResponseEncoder(nil, "/sys/test/error", v1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

}

func TestDefaultResponseEncoderWithError(t *testing.T) {

	se := errors.New(511, "", "")
	DefaultErrorEncoder(nil, "/device/test/error", se)

}

func TestDefaultResponseEncoderEncodeNil(t *testing.T) {

	err := DefaultResponseEncoder(nil, "/device/test/error", nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
