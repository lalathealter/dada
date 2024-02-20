package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const CONTEXT_JWT_KEY = "jwt-token"
const CONTEXT_JWT_CLAIMS_KEY = "jwt-claims"

func ValidateJWT(c *gin.Context) {
  tokenStr := GetCookieJWT(c)
  if tokenStr == "" {
		c.AbortWithError(http.StatusUnauthorized, ErrNoTokenProvided)
    return
  }

	claims := jwt.MapClaims{}
	tokenObj, err := jwt.ParseWithClaims(tokenStr, claims, prepForParsingJWT)
	if err != nil || !tokenObj.Valid {
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
