package global

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/transport/http"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/message"
	stdhttp "net/http"
	"reflect"
	"strconv"
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

type JWTInfo struct {
	ActorName     string `json:"-"`
	ActorUID      string `json:"-"`
	ActorAvatar   string `json:"-"`
	ActorIDLegacy string `json:"-"`
	Phone         string `json:"-"`
	Topic         []int  `json:"-"`
	Email         string `json:"-"`
	ID            string `json:"-"`
	IsVerified    bool   `json:"-"`
	CustomerID    uint64 `json:"-"`
	GroupID       int    `json:"-"`
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

func ExtractClaimsOnly(bearerToken string) (jwtgo.MapClaims, error) {
	// checking empty bearer token
	if bearerToken == "" {
		return nil, errors.New("bearer Token is empty")
	}

	token, _ := jwtgo.ParseWithClaims(bearerToken, jwtgo.MapClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(""), nil
	})

	// check valid bearer token with publishing jwt-login-secret
	if token == nil {
		return nil, errors.New("bearer Token is empty")
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return nil, errors.New("failed casting to jwt MapClaims")
	}

	return claims, nil
}

func ExtractToken(bearerToken string) (*JWTInfo, error) {
	// checking empty bearer token
	if bearerToken == "" {
		return nil, errors.New("bearer Token is empty")
	}

	var claimsJWT JWTInfo
	mapClaims, errMapClaims := ExtractClaimsOnly(bearerToken)
	if errMapClaims != nil {
		return nil, errMapClaims
	}

	var claimID string
	// handle ID if not string
	if reflect.TypeOf(mapClaims["id"]).Kind() == reflect.Float64 {
		claimID = strconv.Itoa(int(mapClaims["id"].(float64)))
	} else {
		claimID = mapClaims["id"].(string)
	}
	// check attribute avatar
	avatar, ok := mapClaims["avatar"]
	if !ok {
		claimsJWT.ActorAvatar = ""
	} else {
		claimsJWT.ActorAvatar = fmt.Sprintf("%v", avatar)
	}
	var groupID int
	var customerID int64
	if tmpGroupID, ok := mapClaims["group_id"]; ok {
		groupID = int(tmpGroupID.(float64))
	}
	if tmpCustomerID, ok := mapClaims["customer_id"]; ok {
		customerID = int64(tmpCustomerID.(float64))
	}

	claimsJWT.ActorName = fmt.Sprintf("%v", mapClaims["full_name"])
	claimsJWT.ActorUID = fmt.Sprintf("%v", mapClaims["sub"])
	claimsJWT.ActorIDLegacy = fmt.Sprintf("%v", mapClaims["user_id_legacy"])
	claimsJWT.Phone = fmt.Sprintf("%v", mapClaims["phone"])
	claimsJWT.Email = fmt.Sprintf("%v", mapClaims["email"])
	claimsJWT.ID = fmt.Sprintf("%v", claimID)
	claimsJWT.CustomerID = uint64(customerID)
	claimsJWT.GroupID = groupID

	return &claimsJWT, nil
}

func GetContextTokenFromHTTP(ctx context.Context, r *stdhttp.Request) context.Context {
	token, ok := ExtractTokenFromAuthHeader(r.Header.Get("Authorization"))
	if !ok {
		return ctx
	}
	payload, _ := ExtractToken(token)
	ctx = r.Context()
	ctx = context.WithValue(ctx, jwt.JWTContextKey, token)

	return context.WithValue(ctx, jwt.JWTClaimsContextKey, payload)
}

func HTTPToContextJWTClaims() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		return GetContextTokenFromHTTP(ctx, r)
	}
}

func HTTPHeaderToContext() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		ctx = GetContextTokenFromHTTP(ctx, r)

		// get header device_id
		deviceID := 1
		headerDeviceID, err := strconv.ParseInt(r.Header.Get("device_id"), 10, 8)
		if headerDeviceID != 0 && err == nil {
			deviceID = int(headerDeviceID)
		}
		ctx = context.WithValue(ctx, "device_id", deviceID)
		return ctx
	}
}
