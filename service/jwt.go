package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/farseer810/file-manager/model"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func ParseUserFromJwt(ctx *gin.Context) *model.User {
	tokenString, err := ctx.Cookie(TokenCookieKey)
	if err != nil {
		return nil
	}

	token, err := jwt.Parse(tokenString, jwtSignatureKeyFunc)
	if err != nil || !token.Valid {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	userInfoJson, ok := claims["userInfo"].(string)
	if !ok {
		return nil
	}

	var user model.User
	err = jsoniter.Unmarshal([]byte(userInfoJson), &user)
	if err != nil {
		return nil
	}
	return &user
}
