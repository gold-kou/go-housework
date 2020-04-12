package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model"

	"github.com/dgrijalva/jwt-go"
)

// TokenInterface interface
type TokenInterface interface {
	GenerateToken(string) (string, error)
	VerifyHeaderToken(*http.Request) (*model.Auth, error)
}

// TokenStruct struct
type TokenStruct struct{}

// NewTokenStruct constructor
func NewTokenStruct() *TokenStruct {
	return &TokenStruct{}
}

// GenerateToken generate and returns JWT
func (t *TokenStruct) GenerateToken(userName string) (tokenString string, err error) {
	// headerのセット
	token := jwt.New(jwt.SigningMethodHS256)

	// claimsのセット
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = userName
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// 電子署名
	tokenString, err = token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyHeaderToken verify token and get auth info
func (t *TokenStruct) VerifyHeaderToken(r *http.Request) (*model.Auth, error) {
	// get jwt from header
	authHeader := r.Header.Get("Authorization")

	// Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM3MjIwNTMsImlhdCI6IjIwMjAtMDMtMDhUMTE6NDc6MzMuMTc4NjU5MyswOTowMCIsIm5hbWUiOiJ0ZXN0In0.YIyT1RJGcYbdynx1V4-6MhiosmTlHmKiyiG_GjxQeuw
	bearerToken := strings.Split(authHeader, " ")[1]

	// verify jwt
	authUser, err := verifyToken(bearerToken)
	if err != nil {
		return nil, common.NewAuthorizationError(err.Error())
	}
	return authUser, nil
}

// verifyToken verify token and return user name
func verifyToken(tokenString string) (*model.Auth, error) {
	// verify
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := common.NewAuthorizationError("unexpected signing method")
			return nil, err
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	// check the result
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, common.NewAuthorizationError("token is expired")
			}
			return nil, common.NewAuthorizationError("token is invalid")
		}
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, common.NewAuthorizationError("not found claims in token")
	}
	userName, ok := claims["name"].(string)
	if !ok {
		return nil, common.NewAuthorizationError("not found name in token")
	}

	return &model.Auth{
		UserName: userName,
	}, nil
}
