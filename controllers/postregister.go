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
		c.Error(ErrNonLatinCharactersInName{user.Name})
		return
	}

	if wr.users.IsUnique(user.Name) {
		c.Error(ErrNameNotUnique{user.Name})
		return
	}

	if err := wr.users.SaveNewUser(user); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
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
