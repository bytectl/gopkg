package uid

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type uidKey struct{}

const (
	authUserKey string = "X-Auth-User"
	reason      string = "UNAUTHORIZED"
)

var (
	ErrWrongContext = errors.Unauthorized(reason, fmt.Sprintf("can't get authUserKey(%s) data from header", authUserKey))
)

type UserData struct {
	ID string
}

// Server is a server auth middleware. Get the user data from Header.
func Server() middleware.Middleware {

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				id := header.RequestHeader().Get(authUserKey)
				if id != "" {
					ctx = NewContext(ctx, &UserData{
						ID: id,
					})
				} else {
					return nil, ErrWrongContext
				}
			}
			return handler(ctx, req)
		}
	}
}

// NewContext put user info into context
func NewContext(ctx context.Context, info *UserData) context.Context {
	return context.WithValue(ctx, uidKey{}, info)
}

// FromContext extract user info from context
func FromContext(ctx context.Context) (info *UserData, ok bool) {
	info, ok = ctx.Value(uidKey{}).(*UserData)
	return
}
