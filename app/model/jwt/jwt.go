package model_jwt

import (
	jwtgo "github.com/golang-jwt/jwt"
)

type ClaimsJWT struct {
	Exp            int      `json:"exp"`
	Iat            int      `json:"iat"`
	Jti            string   `json:"jti"`
	Iss            string   `json:"iss"`
	Sub            string   `json:"sub"`
	Typ            string   `json:"typ"`
	Azp            string   `json:"azp"`
	SessionState   string   `json:"session_state"`
	Acr            string   `json:"acr"`
	AllowedOrigins []string `json:"allowed-origins"`
	Scope          string   `json:"scope"`
	Sid            string   `json:"sid"`
	UserIDLegacy   string   `json:"user_id_legacy"`
	MemberID       string   `json:"member_id"`
	FullName       string   `json:"full_name"`
	GroupID        int      `json:"group_id"`
	Phone          string   `json:"phone"`
	Topic          []int    `json:"topic"`
	Avatar         string   `json:"avatar"`
	ID             string   `json:"id"`
	CustomerID     int      `json:"customer_id"`
	ContactID      string   `json:"contact_id"`
	Email          string   `json:"email"`
	Group          string   `json:"group"`
	jwtgo.StandardClaims
}
