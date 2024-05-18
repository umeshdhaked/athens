package jwt

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("my-jwt-secret-key")

func CreateToken(id int64, mobile string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			constants.JwtTokenUserID: id,
			constants.JwtTokenMobile: mobile,
			constants.JwtTokenRole:   role,
			"exp":                    time.Now().Add(time.Hour * 24).Unix(),
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

func RefreshToken(ctx *gin.Context) (*dtos.LoginSuccessResponse, error) {
	jwtToken := ctx.Request.Header["Authorization"][0]

	if er := VerifyToken(jwtToken); er != nil {
		log.Println("INVALID_TOKEN")
		return nil, errors.Join(er, errors.New("INVALID_TOKEN"))
	}
	claims, _ := DecodeToken(jwtToken)
	exp := claims["exp"].(float64)
	currTime := time.Now().Unix()
	mobile := claims[constants.JwtTokenMobile].(string)
	role := claims[constants.JwtTokenRole].(string)
	id := claims[constants.JwtTokenUserID].(float64)
	if int64(exp) < currTime {
		return nil, errors.New("TOKEN_EXPIRED")
	}
	if int64(exp)-currTime <= 7200 { // you can refresh token only when 2hr is left to existing one expiry.
		tkn, err := CreateToken(int64(id), mobile, role)
		if err != nil {
			return nil, errors.Join(err, errors.New("internal server error"))
		}
		return &dtos.LoginSuccessResponse{MobileNumber: mobile, LoginToken: tkn}, nil
	} else {
		return nil, errors.New("TOKEN_REFRESH_NOT_ALLOWED")
	}
}
