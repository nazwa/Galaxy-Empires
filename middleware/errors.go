package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/stvp/rollbar"
	"gopkg.in/bluesuncorp/validator.v5"
	"net/http"
)

var (
	ErrorInternalError = errors.New("Woops! Something went wrong :(")
)

func ValidationErrorToText(e *validator.FieldError) string {
	switch e.Tag {
	case "required":
		return fmt.Sprintf("%s is required", e.Field)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", e.Field, e.Param)
	case "min":
		return fmt.Sprintf("%s must be longer than %s", e.Field, e.Param)
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field, e.Param)
	}
	return fmt.Sprintf("%s is not valid", e.Field)
}

// This method collects all errors and submits them to Rollbar
func Errors(env, token string, logger service.Logger) gin.HandlerFunc {
	rollbar.Environment = env
	rollbar.Token = token

	return func(c *gin.Context) {
		c.Next()
		// Only run if there are some errors to handle
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Find out what type of error it is
				switch e.Type {
				case gin.ErrorTypePublic:
					// Only output public errors if nothing has been written yet
					if !c.Writer.Written() {
						c.JSON(c.Writer.Status(), gin.H{"Error": e.Error()})
					}
				case gin.ErrorTypeBind:
					errs := e.Err.(*validator.StructErrors)
					list := make(map[string]string)
					for field, err := range errs.Errors {
						list[field] = ValidationErrorToText(err)
					}

					// Make sure we maintain the preset response status
					status := http.StatusBadRequest
					if c.Writer.Status() != http.StatusOK {
						status = c.Writer.Status()
					}
					c.JSON(status, gin.H{"Errors": list})

				default:
					// Log all other errors
					rollbar.RequestError(rollbar.ERR, c.Request, e.Err)
					if logger != nil {
						logger.Error(e.Err)
					}
				}

			}
			// If there was no public or bind error, display default 500 message
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": ErrorInternalError.Error()})
			}
		}
	}
}