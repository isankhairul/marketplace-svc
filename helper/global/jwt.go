package global

import (
	"context"
	"errors"
	"fmt"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/message"
	"strings"

	"github.com/go-kit/kit/auth/jwt"
	jwtgo "github.com/golang-jwt/jwt"
)

const (
	bearer string = "bearer"
)

type JWTPayload struct {
	UID            string            `json:"sub,omitempty"`
	Phone          string            `json:"phone,omitempty"`
	Email          string            `json:"email,omitempty"`
	Name           string            `json:"full_name,omitempty"`
	Groups         []string          `json:"groups,omitempty"`
	ResourceAccess map[string]Access `json:"resource_access,omitempty"`
	Data           UserDoctor        `json:"data,omitempty"`
	jwtgo.StandardClaims
}

type Access struct {
	Roles []string `json:"roles"`
}

type UserDoctor struct {
	ID        int        `json:"id,omitempty"`
	UserID    int        `json:"user_id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	Telephone string     `json:"telephone,omitempty"`
	Vip       bool       `json:"vip,omitempty"`
	Apotek    string     `json:"apotek,omitempty"`
	Hospitals []Hospital `json:"hospitals,omitempty"`
}

type Hospital struct {
	Key             int    `json:"key,omitempty"`
	Value           string `json:"value,omitempty"`
	Role            string `json:"role,omitempty"`
	InstitutionType int    `json:"institution_type,omitempty"`
}

func GetJWTInfoFromContext(ctx context.Context, cfg *config.JwtConfig) (*JWTPayload, error) {
	if !cfg.SkipValidation {
		return JWTInfoToStruct(fmt.Sprint(ctx.Value(jwt.JWTContextKey)))
	}
	return nil, errors.New(message.ErrNoAuth.Message)
}

func JWTInfoToStruct(jwtToken string) (*JWTPayload, error) {
	var token, _, err = new(jwtgo.Parser).ParseUnverified(jwtToken, &JWTPayload{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTPayload); ok {
		return claims, nil
	}
	return nil, errors.New(message.ErrNoAuth.Message)
}

func ExtractTokenFromAuthHeader(val string) (token string, ok bool) {
	authHeaderParts := strings.Split(val, " ")

	if len(authHeaderParts) > 1 {
		return authHeaderParts[1], true
	}

	return "", false
}
