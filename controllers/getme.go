package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lalathealter/dada/models"
	"github.com/lalathealter/dada/controllers/auth"
)

func (wr *Wrapper) HandleViewSelf(c *gin.Context) {
  claims := auth.GetClaimsFromJWT(c)
  name, err := auth.ExtractSubjectFrom(claims)
  if err != nil {
    c.AbortWithError(http.StatusBadRequest, err)
    return
  }

  info := wr.users.GetInfo(models.User{Name: name})
  c.JSON(http.StatusOK, info)
}


