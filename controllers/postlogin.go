package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
  "time"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/lalathealter/dada/models"
)



func (wr *Wrapper) HandleLogin(c *gin.Context) {
  
  var userCred models.User
  if err := c.ShouldBindJSON(&userCred); err != nil {
    c.AbortWithError(http.StatusBadRequest, err)
    return
  }

  
  if !wr.users.HasValidLogin(userCred) {
    c.AbortWithError(http.StatusUnauthorized, ErrInvalidLogin)
    return
  }

  tokenStr, err := forgeUserJWT(userCred)
  if err != nil {
    c.AbortWithError(http.StatusInternalServerError, err)
    return
  }

  c.JSON(http.StatusCreated, gin.H{
    "token": tokenStr,
  })
}

var ErrInvalidLogin = errors.New("The provided login pair is invalid")

func forgeUserJWT(u models.User) (string, error) {
  claims := jwt.StandardClaims{
    Subject: u.Name,
    ExpiresAt: signExpirationTime(),
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  secret := retrieveSecret()
  signedToken, err := token.SignedString(secret)
  return signedToken, err
}

const MAX_ACCESS_MINUTES = 12 * 60
func signExpirationTime() int64 {
  exp := time.Now().Add(time.Minute * time.Duration(MAX_ACCESS_MINUTES)).Unix()
  return exp
}


var retrieveSecret = func()func()[]byte {
  JWT_SECRET := loadSecret()
  byted := []byte(JWT_SECRET)
  return func () []byte  {
    return byted
  }
}()

func loadSecret() string {
  sec := os.Getenv("DADA_JWT_SECRET")
  if sec == "" {
    byteSec := make([]byte, 32)
    rand.Read(byteSec)

    sec = base64.StdEncoding.EncodeToString(byteSec)
    os.Setenv("DADA_JWT_SECRET", sec)
    fmt.Println("Setting a new JWT_SECRET to the os environment")
  }

  return sec
}
