package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lalathealter/dada/models"
	"github.com/lalathealter/dada/controllers/auth"
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

  tokenStr, err := auth.ForgeUserJWT(userCred)
  if err != nil {
    c.AbortWithError(http.StatusInternalServerError, err)
    return
  }

  c.JSON(http.StatusCreated, gin.H{
    "token": tokenStr,
  })
}

var ErrInvalidLogin = errors.New("The provided login pair is invalid")

