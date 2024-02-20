package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lalathealter/dada/controllers/auth"
	"github.com/lalathealter/dada/models"
)


func (wr *Wrapper) HandleUsernameChange(c *gin.Context) {
  var nup models.NewUsernamePackage
  if err := c.ShouldBindJSON(&nup); err != nil {
    c.AbortWithError(http.StatusBadRequest, err)
    return
  }

  claims := auth.GetClaimsFromJWT(c)
  name, err := auth.ExtractSubjectFrom(claims)
  if err != nil {
    c.AbortWithError(http.StatusBadRequest, err)
    return
  }

  err = wr.users.ChangeUsername(name, nup.Name)
  if err != nil {
    c.AbortWithError(http.StatusInternalServerError, err)
    return
  }

  c.Status(http.StatusCreated)
}
