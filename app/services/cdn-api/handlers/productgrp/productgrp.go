// Package productgrp maintains the group of handlers for product access.
package productgrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/core/crud/product"
	"github.com/testvergecloud/testApi/business/core/crud/user"
	wb "github.com/testvergecloud/testApi/business/web"
	"github.com/testvergecloud/testApi/business/web/mid"
	"github.com/testvergecloud/testApi/business/web/page"
	wf "github.com/testvergecloud/testApi/foundation/web"
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
func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewProduct
	if err := wf.Decode(r, &app); err != nil {
		return wb.NewTrustedError(err, http.StatusBadRequest)
	}

	prd, err := h.product.Create(ctx, toCoreNewProduct(ctx, app))
	if err != nil {
		return fmt.Errorf("create: app[%+v]: %w", app, err)
	}

	return wf.Respond(ctx, w, toAppProduct(prd), http.StatusCreated)
}

// update updates a product in the system.
func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdateProduct
	if err := wf.Decode(r, &app); err != nil {
		return wb.NewTrustedError(err, http.StatusBadRequest)
	}

	prd := mid.GetProduct(ctx)

	updPrd, err := h.product.Update(ctx, prd, toCoreUpdateProduct(app))
	if err != nil {
		return fmt.Errorf("update: productID[%s] app[%+v]: %w", prd.ID, app, err)
	}

	return wf.Respond(ctx, w, toAppProduct(updPrd), http.StatusOK)
}

// delete removes a product from the system.
func (h *handlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	prd := mid.GetProduct(ctx)

	if err := h.product.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: productID[%s]: %w", prd.ID, err)
	}

	return wf.Respond(ctx, w, nil, http.StatusNoContent)
}

// query returns a list of products with paging.
func (h *handlers) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, err := page.Parse(r)
	if err != nil {
		return err
	}

	filter, err := parseFilter(r)
	if err != nil {
		return err
	}

	orderBy, err := parseOrder(r)
	if err != nil {
		return err
	}

	prds, err := h.product.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.product.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return wf.Respond(ctx, w, wb.NewPageDocument(toAppProducts(prds), total, page.Number, page.RowsPerPage), http.StatusOK)
}

// queryByID returns a product by its ID.
func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return wf.Respond(ctx, w, toAppProduct(mid.GetProduct(ctx)), http.StatusOK)
}

// create adds a new product to the system.
func (h *handlers) ginCreate(c *gin.Context) error {
	var app AppNewProduct
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	ctx := c.Request.Context()
	prd, err := h.product.Create(ctx, toCoreNewProduct(ctx, app))
	if err != nil {
		return fmt.Errorf("create: app[%+v]: %w", app, err)
	}

	c.JSON(http.StatusCreated, toAppProduct(prd))
	return nil
}

// update updates a product in the system.
func (h *handlers) ginUpdate(c *gin.Context) error {
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
func (h *handlers) ginDelete(c *gin.Context) error {
	ctx := c.Request.Context()
	prd := mid.GetProduct(ctx)

	if err := h.product.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: productID[%s]: %w", prd.ID, err)
	}

	c.JSON(http.StatusNoContent, nil)
	return nil
}

// query returns a list of products with paging.
func (h *handlers) ginQuery(c *gin.Context) error {
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
func (h *handlers) ginQueryByID(c *gin.Context) error {
	c.JSON(http.StatusOK, toAppProduct(mid.GetProduct(c.Request.Context())))
	return nil
}
