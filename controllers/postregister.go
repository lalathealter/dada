package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lalathealter/dada/models"
)

func (wr *Wrapper) HandleRegistration(c *gin.Context) {
	var user models.User

  if err := c.ShouldBindJSON(&user); err != nil {
    c.AbortWithError(http.StatusBadRequest, err)
		return
	}

  if err := wr.validateName(user.Name); err != nil {
    c.AbortWithError(http.StatusBadRequest, err)
    return
  }

	if err := wr.users.SaveNewUser(user); err != nil {
    c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (wr *Wrapper) validateName(name string) error {
  if !isJustLatin(name) {
    return ErrNonLatinCharactersInName{name}
	}

  if !wr.users.IsUnique(name) {
    return ErrNameNotUnique{name}
  }

  return nil
}

type ErrNonLatinCharactersInName struct {
	Name string
}

func (e ErrNonLatinCharactersInName) Error() string {
	return fmt.Sprintf("the provided name [%v] has non-latin characters", e.Name)
}

type ErrNameNotUnique struct {
	Name string
}

func (e ErrNameNotUnique) Error() string {
	return fmt.Sprintf("the provided name [%v] is not unique", e.Name)
}
