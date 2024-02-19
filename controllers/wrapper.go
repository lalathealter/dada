package controllers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/lalathealter/dada/models"
)


func InitWrapper(db *sql.DB) *Wrapper {
	wr := Wrapper{
		models.UserModel{DB: db},
	}

	return &wr
}

type Wrapper struct {
	users models.UserCollectionI
}

func wrapError(err error) ErrorWrapper {
  return ErrorWrapper{err.Error()}
}
type ErrorWrapper struct {
  Error string `json:"error"`
}

func HandleErrors(c *gin.Context) {
  c.Next()

  for _, err := range c.Errors {
    c.JSON(-1, wrapError(err.Err))
  }
}


func ObligateToUseJSON(c *gin.Context) {
  c.Writer.Header().Set("Content-Type", "application/json")
  c.Next()
}
