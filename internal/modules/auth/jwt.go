package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: IMPROVE AND CREATE A NEW PACKAGE FOR JWT MANAGEMENT, THIS IS JUST A TEMPORARY SOLUTION
// ADD A VALIDATION FUNCTION TO VALIDATE THE TOKEN AND RETURN THE CLAIMS, AND A FUNCTION TO REFRESH THE TOKEN
// ADD A FUNCTION TO EXTRACT THE TOKEN FROM THE REQUEST HEADER AND VALIDATE IT, RETURNING THE CLAIMS...

var secret = []byte("change-this-secret")

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, role string) (string, error) {

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
