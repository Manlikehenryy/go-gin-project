package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manlikehenryy/go-gin-project/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IsAuthenticated(c *gin.Context) {
    // Retrieve the JWT token from the cookie
    cookie, err := c.Cookie("jwt")
    if err != nil {
        helpers.SendError(c, http.StatusUnauthorized, "Unauthorized: No JWT token provided")
        c.Abort() // Abort the request pipeline if authentication fails
        return
    }

    // Parse the JWT token from the cookie
    userIdStr, err := helpers.ParseJwt(cookie)
    if err != nil {
        helpers.SendError(c, http.StatusUnauthorized, "Unauthorized: Invalid JWT token")
        c.Abort() // Abort the request pipeline if token parsing fails
        return
    }

    // Convert the user ID string to primitive.ObjectID
    userId, err := primitive.ObjectIDFromHex(userIdStr)
    if err != nil {
        helpers.SendError(c, http.StatusUnauthorized, "Unauthorized: Invalid user ID")
        c.Abort() // Abort the request pipeline if user ID conversion fails
        return
    }

    // Store the user ID in the request context
    c.Set("userId", userId)

    // Continue to the next middleware or handler
    c.Next()
}