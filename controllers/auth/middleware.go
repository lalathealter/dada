package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Auth struct {
	Token string `json:"token" binding:"required"`
}

const CONTEXT_JWT_KEY = "jwt-token"
const CONTEXT_JWT_CLAIMS_KEY = "jwt-claims"

func ValidateJWT(c *gin.Context) {
	var au Auth
	if err := c.ShouldBindJSON(&au); err != nil {
		c.AbortWithError(http.StatusUnauthorized, ErrNoTokenProvided)
		return
	}

	claims := jwt.MapClaims{}
	tokenObj, err := jwt.ParseWithClaims(au.Token, claims, prepForParsingJWT)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, ErrCouldntParseJWToken)
		return
	} else if !tokenObj.Valid {
		c.AbortWithError(http.StatusBadRequest, ErrInvalidToken)
		return
	}

	c.Set(CONTEXT_JWT_KEY, tokenObj)
	c.Set(CONTEXT_JWT_CLAIMS_KEY, claims)
	c.Next()
}

func prepForParsingJWT(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, ErrUnexpectedSigningMethod{token}
	}
	return RetrieveSecret(), nil
}
