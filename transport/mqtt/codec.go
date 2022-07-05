package mqtt

import (
	"bytes"
	"fmt"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
)

// SupportPackageIsVersion1 These constants should not be referenced from any other code.
const SupportPackageIsVersion1 = true

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func([]byte, interface{}) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(*bytes.Buffer, interface{}) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(*bytes.Buffer, error)

// DefaultRequestDecoder decodes the request body to object.
func DefaultRequestDecoder(data []byte, v interface{}) error {
	if len(data) == 0 {
		return nil
	}
	codec := encoding.GetCodec(json.Name)
	if err := codec.Unmarshal(data, v); err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	return nil
}

// DefaultResponseEncoder encodes the object to the mqtt reply.
func DefaultResponseEncoder(w *bytes.Buffer, v interface{}) error {
	if v == nil {
		return nil
	}
	codec := encoding.GetCodec(json.Name)
	body, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	w.Write(body)
	return nil
}

// DefaultErrorEncoder encodes the error to the mqtt response.
func DefaultErrorEncoder(w *bytes.Buffer, err error) {
	var reply struct {
		Id      string `json:"id"`
		Code    int32  `json:"code"`
		Reason  string `json:"reason,omitempty"`
		Message string `json:"message"`
	}
	se := errors.FromError(err)
	reply.Id = se.Metadata["id"]
	reply.Code = se.Code
	reply.Message = se.Message
	reply.Reason = se.Reason
	codec := encoding.GetCodec(json.Name)
	body, err := codec.Marshal(reply)
	if err != nil {
		errBody := fmt.Sprintf("%s", err)
		w.Write([]byte(errBody))
		return
	}
	w.Write(body)
}
