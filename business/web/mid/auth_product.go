package mid

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/core/crud/product"
	"github.com/testvergecloud/testApi/business/web/auth"

	"github.com/google/uuid"
)

type ctxProductKey int

const productKey ctxProductKey = 1

// GetProduct returns the product from the context.
func GetProduct(ctx context.Context) product.Product {
	v, ok := ctx.Value(productKey).(product.Product)
	if !ok {
		return product.Product{}
	}
	return v
}

func setProduct(ctx context.Context, prd product.Product) context.Context {
	return context.WithValue(ctx, productKey, prd)
}

// AuthorizeProduct executes the specified role and extracts the specified
// product from the DB if a product id is specified in the call. Depending on
// the rule specified, the userid from the claims may be compared with the
// specified user id from the product.
func AuthorizeProduct(a *auth.Auth, rule string, prdCore *product.Core) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID uuid.UUID

		if id := c.Param("product_id"); id != "" {
			var err error
			productID, err := uuid.Parse(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid ID: %s", err)})
				c.Abort()
				return
			}

			prd, err := prdCore.QueryByID(c, productID)
			if err != nil {
				switch {
				case errors.Is(err, product.ErrNotFound):
					c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
					c.Abort()
					return
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("querybyid: productID[%s]: %w", productID, err)})
					c.Abort()
					return
				}
			}

			userID = prd.UserID
			c.Request = c.Request.WithContext(setProduct(c.Request.Context(), prd))
		}

		claims := getClaims(c.Request.Context())

		if err := a.Authorize(c.Request.Context(), claims, userID, rule); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)})
			c.Abort()
			return
		}

		c.Next()
	}
}
