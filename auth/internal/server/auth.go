package server

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))
var jwtKid = []byte(os.Getenv("JWT_KID"))

func generateJWT(userID int, email string) (map[string]interface{}, error) {
	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	jti := fmt.Sprintf("%d", userID)

	accessTokenClaims := &jwt.StandardClaims{
		Audience:  "http://localhost:8080",
		Issuer:    "https://krakend.io",
		Subject:   email,
		Id:        jti,
		ExpiresAt: expirationTime,
	}

	refreshTokenClaims := &jwt.StandardClaims{
		Audience:  "http://localhost:8080",
		Issuer:    "https://krakend.io",
		Subject:   email,
		Id:        jti + "_refresh", // Slightly different JTI for refresh token
		ExpiresAt: expirationTime,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken.Header["kid"] = jwtKid
	accessToken.Header["alg"] = jwt.SigningMethodHS256.Alg()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshToken.Header["kid"] = jwtKid
	refreshToken.Header["alg"] = jwt.SigningMethodHS256.Alg()

	responsePayload := map[string]interface{}{
		"access_token": map[string]interface{}{
			"aud": accessTokenClaims.Audience,
			"iss": accessTokenClaims.Issuer,
			"sub": accessTokenClaims.Subject,
			"jti": accessTokenClaims.Id,
			"exp": accessTokenClaims.ExpiresAt,
		},
		"refresh_token": map[string]interface{}{
			"aud": refreshTokenClaims.Audience,
			"iss": refreshTokenClaims.Issuer,
			"sub": refreshTokenClaims.Subject,
			"jti": refreshTokenClaims.Id,
			"exp": refreshTokenClaims.ExpiresAt,
		},
		"exp": expirationTime,
	}

	return responsePayload, nil
}

func validateJWT(tokenStr string) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}
