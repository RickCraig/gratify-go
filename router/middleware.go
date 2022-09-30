package router

import (
	"github.com/gin-gonic/gin"
)

// Pulls the user authentication token from the database
// and ensures the user has access to the endpoint.
//
// 	{ token: "12345abcd" }
//
// As a bonus it also attaches the user object to the context
//
//	{ user: { _id: "abcd123", email: "rick@me.com" ... } }
//
// Returns handler function to be used as middleware
func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: Add the authentication check in here
		ctx.Next()
	}
}
