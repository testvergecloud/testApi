package mid

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/testvergecloud/testApi/business/core/crud/home"
	"github.com/testvergecloud/testApi/business/web/auth"
)

type ctxHomeKey string

const homeKey ctxHomeKey = "home"

// GetHome returns the home from the context.
func GetHome(c *gin.Context) home.Home {
	v, ok := c.Get(string(homeKey))
	if !ok {
		return home.Home{}
	}
	return v.(home.Home)
}

func setHome(c *gin.Context, hme home.Home) {
	c.Set(string(homeKey), hme)
}

// AuthorizeHome executes the specified role and extracts the specified
// home from the DB if a home id is specified in the call. Depending on
// the rule specified, the userid from the claims may be compared with the
// specified user id from the home.
func AuthorizeHome(a *auth.Auth, rule string, hmeCore *home.Core) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID uuid.UUID

		if id := c.Param("home_id"); id != "" {
			homeID, err := uuid.Parse(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidID})
				c.Abort()
				return
			}

			hme, err := hmeCore.QueryByID(c, homeID)
			if err != nil {
				switch {
				case errors.Is(err, home.ErrNotFound):
					c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("querybyid: homeID[%s]: %s", homeID, err)})
				}
				c.Abort()
				return
			}

			userID = hme.UserID
			setHome(c, hme)
		}

		claims := getClaims(c)
		if err := a.Authorize(c, claims, userID, rule); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)})
			c.Abort()
			return
		}

		c.Next()
	}
}
