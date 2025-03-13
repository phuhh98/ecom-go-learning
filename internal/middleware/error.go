package middleware

import (
	"ecom-go/pkg/http/response"

	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Execute the rest of the handlers

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err // Get the error object
			response.Error(c, err)
		}
	}
}
