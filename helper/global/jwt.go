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
	Name           string            `json:"full_name,omitempty"`
	Groups         []string          `json:"groups,omitempty"`
	ResourceAccess map[string]Access `json:"resource_access,omitempty"`
	Scope          string            `json:"scope,omitempty"`
	SID            string            `json:"sid,omitempty"`
	MemberID       string            `json:"member_id,omitempty"`
	UserIDLegacy   string            `json:"user_id_legacy,omitempty"`
	Phone          string            `json:"phone,omitempty"`
	GroupID        int               `json:"group_id,omitempty"`
	Topic          []int             `json:"topic,omitempty"`
	ID             string            `json:"id,omitempty"`
	Avatar         string            `json:"avatar,omitempty"`
	ContactId      string            `json:"contact_id,omitempty"`
	CustomerID     int               `json:"customer_id,omitempty"`
	Group          string            `json:"group,omitempty"`
	jwtgo.StandardClaims
}

type Access struct {
	Roles []string `json:"roles"`
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
	if len(authHeaderParts) != 2 || !strings.EqualFold(authHeaderParts[0], bearer) {
		return "", false
	}

	return authHeaderParts[1], true
}
