package utilities

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// ParseJWT parses and validates a JWT token
func ParseJWT(jwtToken string) (jwt.MapClaims, error) {

	// Validate the token
	token, _ := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AppConfig.JWTKey), nil
	})

	// Fetch the claims from the token...
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the audience and time-sensitive claims in the JWT...
		if timeError := jwt.MapClaims.Valid(claims); timeError != nil && !jwt.MapClaims.VerifyAudience(claims, AppConfig.JWTAudience, true) {
			return nil, errors.New("JWT has expired and / or the audience claim is not valid")
		}
		return claims, nil
	}

	return nil, errors.New("Unable to parse JWT")

}
