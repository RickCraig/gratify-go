package router

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/victoriam-go/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Middleware struct {
	Database database.Database
}

var publicPaths []string = []string{
	"/auth/login",
	"/auth/register",
	"/auth/forgot",
	"/auth/reset",
	"/user/by-token",
}

func contains(p string) bool {
	for _, public := range publicPaths {
		r, _ := regexp.Compile(public)
		if r.MatchString(p) {
			return true
		}
	}
	return false
}

type Address struct {
	Address1 string `bson:"address1"`
	Address2 string `bson:"address2"`
	City     string `bson:"city"`
	County   string `bson:"County"`
	Country  string `bson:"country"`
	Phone    string `bson:"phone"`
}

type Customer struct {
	Address Address `bson:"address"`
}

type User struct {
	Id        string    `bson:"_id,omitempty"`
	Email     string    `bson:"email"`
	Customer  Customer  `bson:"customer"`
	CreatedAt time.Time `bson:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty"`
}

type ErrorJSON struct {
	Reason string `json:"reason"`
}

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
func (m *Middleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the route is private of public
		if contains(c.Request.URL.Path) {
			fmt.Println("This is a public endpoint, not checking for auth")
			c.Next()
		} else {
			fmt.Println("This is a private endpoint, auth will be checked")
			// This path requires authorisation to use
			col, ctx, cancel := m.Database.GetCollection("users")
			defer cancel()

			// Get the key passed in the headers
			authKey := c.Request.Header.Get("Authorization")
			if authKey == "" {
				c.JSON(
					http.StatusBadRequest,
					ErrorJSON{Reason: "Authorization header is required for private endpoints"},
				)
				return
			}

			// Check for the key assigned to the user
			// set during login, and saved to browser
			// localstorage if the user chose to allow
			// cookies
			var user User
			err := col.FindOne(ctx, bson.D{{Key: "token", Value: authKey}}).Decode(&user)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					// No user with the accessToken passed exists
					// This should send the user to the login page
					c.JSON(
						http.StatusForbidden,
						ErrorJSON{Reason: "The Authorization token sent does not exist"},
					)
				} else {
					fmt.Println("An error occured while authorising", err)
					c.Status(http.StatusInternalServerError)
				}
				c.Abort()
				return
			}

			// Set the user for use in the api endpoints
			c.Set("user-id", user.Id)
			c.Next()
		}
	}
}

// Sets all the CORS headers for Gin
func (m *Middleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
