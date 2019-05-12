package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alexedwards/scs"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var hmacSampleSecret []byte
var sessionManager = scs.NewCookieManager(viper.GetString("SessionSecret"))

func init() {
	hmacSampleSecret = []byte(viper.GetString("AuthSecret"))
}

func createToken(id string, name string, email string) (tokenString string, err error) {
	expiryTime := time.Now().Add(time.Hour * 24 * time.Duration(viper.GetInt("AuthExpiryDays")))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":                 id,
		"name":               name,
		"email":              email,
		"AuthExpiryDateTime": expiryTime,
	})

	tokenString, err = token.SignedString(hmacSampleSecret)
	return
}

func parseToken(tokenString string) (id string, user string, email string, valid bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Expired
		if time.Now().After(claims["AuthExpiryDateTime"].(time.Time)) {
			log.Println("Token expired.", time.Now(), claims["AuthexpiryDateTime"])
			return "", "", "", false
		}
		return claims["id"].(string), claims["user"].(string), claims["email"].(string), true
	} else {
		log.Println(err)
		return "", "", "", false
	}
	return "", "", "", false
}
