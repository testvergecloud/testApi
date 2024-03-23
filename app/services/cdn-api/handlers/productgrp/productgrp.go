// Package productgrp maintains the group of handlers for product access.
package productgrp

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/core/crud/product"
	"github.com/testvergecloud/testApi/business/core/crud/user"
	wb "github.com/testvergecloud/testApi/business/web"
	"github.com/testvergecloud/testApi/business/web/mid"
	"github.com/testvergecloud/testApi/business/web/page"
)

// Set of error variables for handling product group errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

type handlers struct {
	product *product.Core
	user    *user.Core
}

func new(product *product.Core, user *user.Core) *handlers {
	return &handlers{
		product: product,
		user:    user,
	}
}

// create adds a new product to the system.
func (h *handlers) create(c *gin.Context) error {
	var app AppNewProduct
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	ctx := c.Request.Context()
	prd, err := h.product.Create(ctx, toCoreNewProduct(c, app))
	if err != nil {
		return fmt.Errorf("create: app[%+v]: %w", app, err)
	}

	c.JSON(http.StatusCreated, toAppProduct(prd))
	return nil
}

// update updates a product in the system.
func (h *handlers) update(c *gin.Context) error {
	var app AppUpdateProduct
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	ctx := c.Request.Context()
	prd := mid.GetProduct(ctx)

	updPrd, err := h.product.Update(ctx, prd, toCoreUpdateProduct(app))
	if err != nil {
		return fmt.Errorf("update: productID[%s] app[%+v]: %w", prd.ID, app, err)
	}

	c.JSON(http.StatusOK, toAppProduct(updPrd))
	return nil
}

// delete removes a product from the system.
func (h *handlers) delete(c *gin.Context) error {
	ctx := c.Request.Context()
	prd := mid.GetProduct(ctx)

	if err := h.product.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: productID[%s]: %w", prd.ID, err)
	}

	c.JSON(http.StatusNoContent, nil)
	return nil
}

// query returns a list of products with paging.
func (h *handlers) query(c *gin.Context) error {
	page, err := page.Parse(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	filter, err := parseFilter(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	orderBy, err := parseOrder(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	ctx := c.Request.Context()
	prds, err := h.product.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.product.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	c.JSON(http.StatusOK, wb.NewPageDocument(toAppProducts(prds), total, page.Number, page.RowsPerPage))
	return nil
}

// queryByID returns a product by its ID.
func (h *handlers) queryByID(c *gin.Context) error {
	c.JSON(http.StatusOK, toAppProduct(mid.GetProduct(c.Request.Context())))
	return nil
}
