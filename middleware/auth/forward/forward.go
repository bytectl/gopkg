package forward

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type authKey struct{}

type AuthInfo struct {
	Token      string
	Method     string
	Protocol   string
	Host       string
	RequestUri string
	SourceIp   string
}

const (
	bearerWord       string = "Bearer"
	authorizationKey string = "Authorization"
	methodKey        string = "X-Forwarded-Method"
	protocolKey      string = "X-Forwarded-Protocol"
	hostKey          string = "X-Forwarded-Host"
	requestUriKey    string = "X-Forwarded-Uri"
	sourceIpKey      string = "X-Forwarded-For"
	reason           string = "UNAUTHORIZED"
)

// Server is a server auth middleware. Get the forward auth data from Header.
func Server() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			token := ""
			if header, ok := transport.FromServerContext(ctx); ok {
				auths := strings.SplitN(header.RequestHeader().Get(authorizationKey), " ", 2)
				if len(auths) == 2 && strings.EqualFold(auths[0], bearerWord) {
					token = auths[1]
				} else {
					return nil, errors.Unauthorized(reason, "token is empty")
				}
				ctx = NewContext(ctx, &AuthInfo{
					Token:      token,
					RequestUri: header.RequestHeader().Get(requestUriKey),
					Method:     header.RequestHeader().Get(methodKey),
					Protocol:   header.RequestHeader().Get(protocolKey),
					Host:       header.RequestHeader().Get(hostKey),
					SourceIp:   header.RequestHeader().Get(sourceIpKey),
				})
			}
			return handler(ctx, req)
		}
	}
}

// NewContext put auth info into context
func NewContext(ctx context.Context, authInfo *AuthInfo) context.Context {
	return context.WithValue(ctx, authKey{}, authInfo)
}

// FromContext extract auth info from context
func FromContext(ctx context.Context) (authInfo *AuthInfo, ok bool) {
	authInfo, ok = ctx.Value(authKey{}).(*AuthInfo)
	return
}
