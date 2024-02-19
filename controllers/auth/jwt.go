package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/lalathealter/dada/models"
)


func ForgeUserJWT(u models.User) (string, error) {
  claims := jwt.StandardClaims{
    Subject: u.Name,
    ExpiresAt: signExpirationTime(),
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  secret := RetrieveSecret()
  signedToken, err := token.SignedString(secret)
  return signedToken, err
}

const MAX_ACCESS_MINUTES = 12 * 60
func signExpirationTime() int64 {
  exp := time.Now().Add(time.Minute * time.Duration(MAX_ACCESS_MINUTES)).Unix()
  return exp
}


var RetrieveSecret = func()func()[]byte {
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


func GetClaimsFromJWT(c *gin.Context) jwt.MapClaims {
	dataAny, ok := c.Get(CONTEXT_JWT_CLAIMS_KEY)
	if !ok {
		return nil
	}
	return dataAny.(jwt.MapClaims)
}


func ExtractSubjectFrom(claims jwt.MapClaims) (string, error) {
	return extractStringFromClaims(claims, "sub")
}

func extractStringFromClaims(claims jwt.MapClaims, field string) (string, error) {
	if claims == nil {
		return "", ErrNoClaimsProvided
	}

	subI, ok := claims[field]
	if !ok {
		return "", ErrDidntFindFieldInClaims{field}
	}

	sub, ok := subI.(string)
	if !ok {
		return "", ErrWrongDataTypeRequiredFromClaimsField{field}
	}

	return sub, nil
}

