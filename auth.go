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
var expiryTimeFormat string

func init() {
	hmacSampleSecret = []byte(viper.GetString("AuthSecret"))
	expiryTimeFormat = "2006-01-02 15:04:05"
}

func createToken(id string, name string, email string) (tokenString string, err error) {
	// fmt.Println(viper.GetInt("AuthExpiryDays"))
	// fmt.Println(time.Hour * 24 * time.Duration(viper.GetInt("AuthExpiryDays")))
	expiryTime := time.Now().Add(time.Hour * 24 * time.Duration(viper.GetInt("AuthExpiryDays")))
	// fmt.Println(expiryTimeFormat)
	// fmt.Println(expiryTime.Format(expiryTimeFormat))
	// fmt.Printf("expiryTime: %T, %v\n", expiryTime, expiryTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":                 id,
		"name":               name,
		"email":              email,
		"AuthExpiryDateTime": expiryTime.Format(expiryTimeFormat),
	})

	tokenString, err = token.SignedString(hmacSampleSecret)
	return
}

func parseToken(tokenString string) (id string, name string, email string, valid bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Expired
		timeoutTime, err := time.Parse(expiryTimeFormat, claims["AuthExpiryDateTime"].(string))
		if err != nil {
			log.Println(err)
			return "", "", "", false
		}
		// fmt.Println(claims["AuthExpiryDateTime"].(string))
		// fmt.Println(timeoutTime, err)
		if time.Now().After(timeoutTime) {
			log.Println("Token expired. Now / ExpiryTime: ", time.Now(), claims["AuthexpiryDateTime"])
			return "", "", "", false
		}
		// fmt.Println(claims)
		return claims["id"].(string), claims["name"].(string), claims["email"].(string), true
	} else {
		log.Println(err)
		return "", "", "", false
	}
	return "", "", "", false
}
