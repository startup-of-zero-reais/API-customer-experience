package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type (
	// JwtService is a interface to access the jwt data
	JwtService interface {
		GenerateToken(subject, id string) (string, error)
		ValidateToken(token string) (bool, error)
		DecodedToken(key string) interface{}
	}

	// JwtServiceImpl is a implementation of JwtService
	JwtServiceImpl struct {
		decodedToken map[string]interface{}
	}
)

func NewJwtService() JwtService {
	return &JwtServiceImpl{
		decodedToken: map[string]interface{}{},
	}
}

func (j *JwtServiceImpl) GenerateToken(subject, id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   subject,
		Id:        id,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JwtServiceImpl) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("sessão expirada")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		j.decodedToken["id"] = claims.Subject
		j.decodedToken["email"] = claims.Id

		return true, nil
	}

	return false, nil
}

func (j *JwtServiceImpl) DecodedToken(key string) interface{} {
	return j.decodedToken[key]
}
