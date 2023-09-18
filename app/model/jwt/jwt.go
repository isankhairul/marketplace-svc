package model_jwt

import (
	jwtgo "github.com/golang-jwt/jwt"
	"gorm.io/datatypes"
)

type PayloadUser struct {
	// ID        int    `json:"id"`
	UserID     string         `json:"user_id"`
	UserUid    string         `json:"user_uid"`
	Name       string         `json:"name"`
	Email      string         `json:"email"`
	Telephone  string         `json:"telephone"`
	Vip        bool           `json:"vip"`
	Apotek     datatypes.JSON `json:"apotek"`
	Hospitals  []Hospital     `json:"hospitals"`
	UserType   string         `json:"user_type"`
	Profession string         `json:"profession"`
}

type ClaimsJWT struct {
	Exp      int         `json:"exp"`
	Iat      int         `json:"iat"`
	AuthTime int         `json:"auth_time"`
	Data     PayloadUser `json:"data"`
	jwtgo.StandardClaims
}

type Hospital struct {
	Key             int    `json:"key"`
	Value           string `json:"value"`
	Role            string `json:"role"`
	InstitutionType int    `json:"institution_type"`
}

func DataTokenMapToResponse(data ClaimsJWT) *PayloadUser {
	return &PayloadUser{
		UserID:     data.Data.UserID,
		UserUid:    data.Data.UserUid,
		Name:       data.Data.Name,
		Email:      data.Data.Email,
		Apotek:     data.Data.Apotek,
		Telephone:  data.Data.Telephone,
		Hospitals:  data.Data.Hospitals,
		Vip:        data.Data.Vip,
		UserType:   data.Data.UserType,
		Profession: data.Data.Profession,
	}
}
