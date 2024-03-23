package mid

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/testvergecloud/testApi/business/core/crud/user"
	"github.com/testvergecloud/testApi/business/web/auth"
)

type ctxUserKey int

const (
	userIDKey ctxUserKey = iota + 1
	userKey
)

// GetUserID returns the claims from the context.
func GetUserID(c *gin.Context) uuid.UUID {
	v, ok := c.Get(string(userIDKey))
	if !ok {
		return uuid.UUID{}
	}
	return v.(uuid.UUID)
}

// GetUser returns the user from the context.
func GetUser(c *gin.Context) user.User {
	v, ok := c.Get(string(userKey))
	if !ok {
		return user.User{}
	}
	return v.(user.User)
}

func setUserID(c *gin.Context, userID uuid.UUID) {
	c.Set(string(userIDKey), userID)
}

func setUser(c *gin.Context, usr user.User) {
	c.Set(string(userKey), usr)
}

// AuthorizeUser executes the specified role and extracts the specified user
// from the DB if a user id is specified in the call. Depending on the rule
// specified, the userid from the claims may be compared with the specified
// user id.
func AuthorizeUser(a *auth.Auth, rule string, usrCore *user.Core) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID uuid.UUID

		if id := c.Param("user_id"); id != "" {
			var err error
			userID, err = uuid.Parse(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
				c.Abort()
				return
			}

			usr, err := usrCore.QueryByID(c, userID)
			if err != nil {
				c.JSON(http.StatusNoContent, gin.H{"error": "User not found"})
				c.Abort()
				return
			}

			setUser(c, usr)
		}

		claims := getClaims(c)
		if err := a.Authorize(c, claims, userID, rule); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized for that action"})
			c.Abort()
			return
		}

		c.Next()
	}
}
