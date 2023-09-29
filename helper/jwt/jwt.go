package jwt

import (
	"encoding/json"
	"errors"
	jwtgo "github.com/golang-jwt/jwt"
	model_jwt "marketplace-svc/app/model/jwt"
	"marketplace-svc/helper/config"
)

func ExtractClaims(bearerToken string, cfg config.JwtConfig) (jwtgo.MapClaims, error) {
	// checking empty bearer token
	if bearerToken == "" {
		return nil, errors.New("bearer Token is empty")
	}

	token, _ := jwtgo.ParseWithClaims(bearerToken, jwtgo.MapClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(cfg.Key), nil
	})

	// validate in krakend, only extract
	// err include: expired, not valid secret-key
	//if err != nil {
	//	return nil, err
	//}

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

func ExtractToken(bearerToken string, cfg config.JwtConfig) (*model_jwt.ClaimsJWT, error) {
	// checking empty bearer token
	if bearerToken == "" {
		return nil, errors.New("bearer Token is empty")
	}

	var claimsJWT model_jwt.ClaimsJWT
	mapClaims, errMapClaims := ExtractClaims(bearerToken, cfg)
	if errMapClaims != nil {
		return nil, errMapClaims
	}

	jsonMapClaims, err := json.Marshal(mapClaims)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonMapClaims, &claimsJWT)
	if err != nil {
		return nil, err
	}

	return &claimsJWT, nil
}

/*
func GenerateJWTMedia(cfg config.MediaSvc) (string, error) {
	jwtSecret := cfg.JWTKey
	var sampleSecretKey = []byte(jwtSecret)
	token := jwtgo.New(jwtgo.SigningMethodHS256)
	claims := token.Claims.(jwtgo.MapClaims)
	// claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["username"] = "klikmedis-svc"
	claims["fullname"] = "klikmedis svc"
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
*/

/*
func ExtractClaimsWithoutExpired(bearerToken string, cfg config.JwtConfig) (jwtgo.MapClaims, error) {
	// checking empty bearer token
	if bearerToken == "" {
		return nil, errors.New("bearer Token is empty")
	}

	token, err := jwtgo.ParseWithClaims(bearerToken, jwtgo.MapClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(cfg.Key), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("Invalid JWT token format")
			} else if ve.Errors&(jwt.ValidationErrorNotValidYet) != 0 { // exclude ValidationErrorExpired
				return nil, errors.New("JWT token is not yet valid")
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, errors.New("JWT token signature is invalid")
			} else {
				return nil, errors.New("invalid JWT token")
			}
		}
	}

	// check valid bearer token with publishing jwt-login-secret
	if token == nil {
		return nil, errors.New("bearer Token is empty")
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("failed casting to jwt MapClaims")
	}

	return claims, nil
}
*/
