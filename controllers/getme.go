package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/lalathealter/dada/models"
)

type Auth struct {
  Token string `json:"token" binding:"required"`
}

func (wr *Wrapper) HandleViewSelf(c *gin.Context) {
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

  name, err := extractSubjectFrom(claims)
  if err != nil {
    c.AbortWithError(http.StatusBadRequest, err)
    return
  }

  info := wr.users.GetInfo(models.User{Name: name})
  c.JSON(http.StatusOK, info)
}

func prepForParsingJWT(token *jwt.Token) (interface{}, error) {
  if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
    return nil, ErrUnexpectedSigningMethod{token}
  }
  return retrieveSecret(), nil
}

type ErrUnexpectedSigningMethod struct {
  Token *jwt.Token
}
func (err ErrUnexpectedSigningMethod) Error() string {
  return fmt.Sprintf("Unexpected signing method: %v", err.Token.Header["alg"])
}

var ErrNoTokenProvided = errors.New("No authorization token was provided; Denied access")
var ErrInvalidToken = errors.New("The provided authorization token is invalid; Denied access")
var ErrCouldntParseJWToken = errors.New("Could not parse authorization token")

func extractSubjectFrom(claims jwt.MapClaims) (string, error) {
  if claims == nil {
    return "", ErrNoTokenProvided
  }

  subI, ok := claims["sub"]
  if !ok {
    return "", ErrCouldntParseJWToken
  }

  sub, ok := subI.(string)
  if ! ok {
    return "", ErrCouldntParseJWToken
  }

  return sub, nil
}
