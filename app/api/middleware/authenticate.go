package middleware

import (
	"context"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/global"
	helperjwt "marketplace-svc/helper/jwt"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
)

type AuthenticateError interface {
	error() error
}

func Authenticate(cfg config.JwtConfig) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.SkipValidation {
				h.ServeHTTP(w, r)
				return
			}
			token, _ := global.ExtractTokenFromAuthHeader(r.Header.Get("Authorization"))
			//if !ok {
			//	base.ResponseWriter(w, http.StatusUnauthorized, base.SetDefaultResponse(r.Context(), message.UnauthorizedError))
			//	return
			//}
			// use model_jwt.ClaimsJWT or global.JWTPayload ?
			payload, _ := helperjwt.ExtractToken(token, cfg)

			ctx := context.WithValue(r.Context(), jwt.JWTClaimsContextKey, payload)
			ctx = context.WithValue(ctx, jwt.JWTContextKey, token)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthzToContext use in the ServerBefore method
func AuthzToContext(cfg *config.JwtConfig) httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		if cfg.SkipValidation {
			return ctx
		}
		if token, ok := global.ExtractTokenFromAuthHeader(r.Header.Get("Authorization")); ok {
			payload, err := helperjwt.ExtractToken(token, *cfg)
			if err == nil && payload != nil {
				return context.WithValue(ctx, jwt.JWTContextKey, payload)
			}
		}
		return ctx
	}
}

func IsAuthContext(ctx context.Context) (*global.JWTInfo, bool) {
	payload := ctx.Value(jwt.JWTClaimsContextKey)
	if user, ok := payload.(*global.JWTInfo); ok && user != nil {
		return user, true
	}
	return nil, false
}
