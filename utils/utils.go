package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"

	"blue-admin.com/configs"
	"github.com/golang-jwt/jwt/v5"
)

// var key = config.Config("TOKEN_SALT")

type UserClaim struct {
	jwt.RegisteredClaims
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	UUID  string   `json:"uuid"`
}

// Combine password and salt then hash them using the SHA-512
func HashFunc(password string) string {

	// var salt []byte
	// get salt from env variable
	salt := []byte(configs.AppConfig.Get("SECRETE_SALT"))

	// Convert password string to byte slice
	var pwdByte = []byte(password)

	// Create sha-512 hasher
	var sha512 = sha512.New()

	pwdByte = append(pwdByte, salt...)

	sha512.Write(pwdByte)

	// Get the SHA-512 hashed password
	var hashedPassword = sha512.Sum(nil)

	// Convert the hashed to hex string
	var hashedPasswordHex = hex.EncodeToString(hashedPassword)
	return hashedPasswordHex
}

func PasswordsMatch(hashedPassword, currPassword string) bool {

	var currPasswordHash = HashFunc(currPassword)
	return hashedPassword == currPasswordHash
}

// source of this token encode decode functions
// https://github.com/gurleensethi/go-jwt-tutorial/blob/main/main.go
func CreateJWTToken(email string, uuid string, roles []string, duration int) (string, error) {
	my_claim := UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		Email:            email,
		Roles:            roles,
		UUID:             uuid,
	}

	salt_a, _ := GetJWTSalt()
	exp := time.Now().UTC().Add(time.Duration(duration) * time.Minute)
	my_claim.ExpiresAt = jwt.NewNumericDate(exp)
	my_claim.Issuer = "Blue Admin"
	my_claim.Subject = "UI Authentication Token"
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, my_claim)
	signedString, err := token.SignedString([]byte(salt_a))
	if err != nil {
		return "", fmt.Errorf("error creating signed string: %v", err)
	}

	return signedString, nil
}

func ParseJWTToken(jwtToken string) (UserClaim, error) {
	salt_a, salt_b := GetJWTSalt()
	response_a := UserClaim{}
	response_b := UserClaim{}

	token_a, aerr := jwt.ParseWithClaims(jwtToken, &response_a, func(token *jwt.Token) (interface{}, error) {
		return []byte(salt_a), nil
	})
	token_b, berr := jwt.ParseWithClaims(jwtToken, &response_b, func(token *jwt.Token) (interface{}, error) {
		return []byte(salt_b), nil
	})

	if aerr != nil && berr != nil {
		return UserClaim{}, aerr
	}

	// check token validity, for example token might have been expired
	if !token_a.Valid {
		if !token_b.Valid {
			return UserClaim{}, fmt.Errorf("invalid token with second salt")
		}
		return response_b, nil
	}
	return response_a, nil

}

// Return Unique values in list
func UniqueSlice(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Return Unique values in list
func CheckValueExistsInSlice(slice []string, role_test string) bool {
	for _, role := range slice {
		if role == role_test || role == "superuser" {
			return true
		}
	}
	return false
}
