package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("my_secret_key")

// 生成JWT
func generateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 校验JWT
func validateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["username"].(string), nil
	}
	return "", fmt.Errorf("invalid token")
}

func main() {
	// 生成token
	token, err := generateJWT("alice")
	if err != nil {
		fmt.Println("生成JWT失败:", err)
		return
	}
	fmt.Println("生成的JWT:", token)

	// 校验token
	username, err := validateJWT(token)
	if err != nil {
		fmt.Println("校验失败:", err)
		return
	}
	fmt.Println("校验通过，用户名:", username)
}
