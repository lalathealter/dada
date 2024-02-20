package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type ErrUnexpectedSigningMethod struct {
	Token *jwt.Token
}

func (err ErrUnexpectedSigningMethod) Error() string {
	return fmt.Sprintf("Unexpected signing method: %v", err.Token.Header["alg"])
}

var ErrNoTokenProvided = errors.New("No authorization token was provided; Denied access")
var ErrInvalidToken = errors.New("The provided authorization token is invalid; Denied access")

var ErrNoClaimsProvided = errors.New("No claims were provided with the token; Token is invalid")

type ErrDidntFindFieldInClaims struct {
	Field string
}

func (err ErrDidntFindFieldInClaims) Error() string {
	return fmt.Sprintf("No field like %v was provided by the claims of the given authorization token", err.Field)
}

type ErrWrongDataTypeRequiredFromClaimsField struct {
	Field string
}

func (err ErrWrongDataTypeRequiredFromClaimsField) Error() string {
	return fmt.Sprintf("Could not rightly extract data from the field of %v of claims of the given authorization token;", err.Field)
}
