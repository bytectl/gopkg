package basic

import (
	"context"
	"encoding/base64"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

const (
	// authorizationKey
	authorizationKey string = "Authorization"
	// reason holds the error reason.
	reason string = "UNAUTHORIZED"
)

var (
	ErrWrongContext = errors.Unauthorized(reason, "Wrong context for middleware")
)

// Client is an client basic auth middleware.
func Client(user, password string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if clientContext, ok := transport.FromClientContext(ctx); ok {
				clientContext.RequestHeader().Set(authorizationKey, authorizationHeader(user, password))
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
}
