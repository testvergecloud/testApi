package mid

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/web/auth"

	"github.com/google/uuid"
)

// Set of error variables for handling auth errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(a *auth.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := a.Authenticate(c, c.GetHeader("authorization"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("authenticate: failed: %s", err)})
			c.Abort()
			return
		}

		if claims.Subject == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorize: you are not authorized for that action, no claims"})
			c.Abort()
			return
		}

		subjectID, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidID})
			c.Abort()
			return
		}

		c.Set("userID", subjectID)
		c.Set("claims", claims)

		c.Next()
	}
}

// Authorize executes the specified role and does not extract any domain data.
func Authorize(a *auth.Auth, rule string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := getClaims(c)
		if err := a.Authorize(c, claims, uuid.UUID{}, rule); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)})
			c.Abort()
			return
		}

		subjectID, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidID})
			c.Abort()
			return
		}

		c.Set("userID", subjectID)
		c.Set("claims", claims)

		c.Next()
	}
}
