// Package vproductgrp maintains the group of handlers for detailed product data.
package vproductgrp

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/core/views/vproduct"
	wb "github.com/testvergecloud/testApi/business/web"
	"github.com/testvergecloud/testApi/business/web/page"
)

type handlers struct {
	vProduct *vproduct.Core
}

func new(vProduct *vproduct.Core) *handlers {
	return &handlers{
		vProduct: vProduct,
	}
}

// Query returns a list of products with paging.
func (h *handlers) Query(c *gin.Context) error {
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
	prds, err := h.vProduct.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.vProduct.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	c.JSON(http.StatusOK, wb.NewPageDocument(toAppProducts(prds), total, page.Number, page.RowsPerPage))
	return nil
}
