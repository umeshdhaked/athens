package jwt

import (
	"errors"
	"fmt"
	"github.com/fastbiztech/hastinapura/internal/pkg/models/responses"
	"github.com/gin-gonic/gin"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("my-jwt-secret-key")

func CreateToken(id string, username string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       id,
			"username": username,
			"role":     role,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func RefreshToken(ctx *gin.Context) (*responses.LoginSuccessResponse, error) {
	jwtToken := ctx.Request.Header["Token"][0]

	if er := VerifyToken(jwtToken); er != nil {
		log.Println("INVALID_TOKEN")
		return nil, errors.Join(er, errors.New("INVALID_TOKEN"))
	}
	claims, _ := DecodeToken(jwtToken)
	exp := claims["exp"].(float64)
	currTime := time.Now().Unix()
	userNme := claims["username"].(string)
	role := claims["role"].(string)
	id := claims["id"].(string)
	if int64(exp) < currTime {
		return nil, errors.New("TOKEN_EXPIRED")
	}
	if int64(exp)-currTime < 7200 {
		tkn, err := CreateToken(id, userNme, role)
		if err != nil {
			return nil, errors.Join(err, errors.New("internal server error"))
		}
		return &responses.LoginSuccessResponse{MobileNumber: userNme, LoginToken: tkn}, nil
	} else {
		return nil, errors.New("TOKEN_REFRESH_NOT_ALLOWED")
	}
}
