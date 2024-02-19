package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lalathealter/dada/models"
)

func (wr *Wrapper) HandleRegistration(c *gin.Context) {
	var user models.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.Error(err)
		return
	}

	if !isJustLatin(user.Name) {
    c.Status(http.StatusBadRequest)
		c.Error(ErrNonLatinCharactersInName{user.Name})
		return
	}

  if len(user.Name) > MAX_USERNAME_LEN {
    c.Status(http.StatusBadRequest)
    c.Error(ErrNameTooLong{user.Name})
    return 
  }

  if !wr.users.IsUnique(user) {
    c.Status(http.StatusBadRequest)
    c.Error(ErrNameNotUnique{user.Name})
    return
  }

	if err := wr.users.SaveNewUser(user); err != nil {
    c.Status(http.StatusInternalServerError)
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

const MAX_USERNAME_LEN = 32
type ErrNameTooLong struct {
  Name string
}
func (err ErrNameTooLong) Error() string {
  return fmt.Sprintf("the provided name [%v] is longer than maximumavailable username length of %v", err.Name, MAX_USERNAME_LEN)
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
