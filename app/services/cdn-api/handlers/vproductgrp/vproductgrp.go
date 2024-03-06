// Package vproductgrp maintains the group of handlers for detailed product data.
package vproductgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/testvergecloud/testApi/business/core/views/vproduct"
	wb "github.com/testvergecloud/testApi/business/web"
	"github.com/testvergecloud/testApi/business/web/page"
	wf "github.com/testvergecloud/testApi/foundation/web"
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
func (h *handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

	prds, err := h.vProduct.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.vProduct.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return wf.Respond(ctx, w, wb.NewPageDocument(toAppProducts(prds), total, page.Number, page.RowsPerPage), http.StatusOK)
}
