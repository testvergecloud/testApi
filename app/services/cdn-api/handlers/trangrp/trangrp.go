// Package trangrp maintains the group of handlers for transaction example.
package trangrp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/core/crud/product"
	"github.com/testvergecloud/testApi/business/core/crud/user"
)

type handlers struct {
	user    *user.Core
	product *product.Core
}

func new(user *user.Core, product *product.Core) *handlers {
	return &handlers{
		user:    user,
		product: product,
	}
}

// create adds a new user and product at the same time under a single transaction.
func (h *handlers) create(c *gin.Context) error {
	ctx := c.Request.Context()
	h, err := h.executeUnderTransaction(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	var app AppNewTran
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	np, err := toCoreNewProduct(app.Product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	nu, err := toCoreNewUser(app.User)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	usr, err := h.user.Create(ctx, nu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	np.UserID = usr.ID

	prd, err := h.product.Create(ctx, np)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusCreated, toAppProduct(prd))
	return nil
}
