package accounts

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyError struct{}

func (m *MyError) Error() string {
	return "Problem generating authentication token"
}

var sampleSecretKey = []byte("8978D868DFBEFD831287536AB667")
var refreshTokenKey = []byte(os.Getenv("mysuperSecretKey"))

// Generates a signed JWT token
func GenerateJwt(email string, id int) (string, error) {
	fmt.Printf("PRIVATE KEY: %v\n", sampleSecretKey)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	claims["authorised"] = true
	claims["email"] = email
	claims["sub"] = strconv.Itoa(id)

	// the last part is sign the token with the private key
	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		return "", &MyError{}
	}

	return tokenString, nil
}

func GenerateRefreshToken(email string, id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour * 24 * 7)
	claims["email"] = email
	claims["sub"] = id

	// the last part is sign the token with the private key
	tokenString, err := token.SignedString(refreshTokenKey)

	if err != nil {
		return "", &MyError{}
	}

	return tokenString, nil
}

func ParseJWTToken(tokenString string) (string, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})

	if err != nil {
		return "", err
	}

	fmt.Printf("THESE ARE THE CLAIMS: %v\n", claims)
	val, _ := claims["sub"].(string)
	fmt.Printf("THE VAL: %v\n", val)
	return val, nil
}
