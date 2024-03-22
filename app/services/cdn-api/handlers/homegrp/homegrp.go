// Package homegrp maintains the group of handlers for home access.
package homegrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/core/crud/home"
	wb "github.com/testvergecloud/testApi/business/web"
	"github.com/testvergecloud/testApi/business/web/mid"
	"github.com/testvergecloud/testApi/business/web/page"
	wf "github.com/testvergecloud/testApi/foundation/web"
)

// Set of error variables for handling home group errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

type handlers struct {
	home *home.Core
}

func new(home *home.Core) *handlers {
	return &handlers{
		home: home,
	}
}

// create adds a new home to the system.
func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewHome
	if err := wf.Decode(r, &app); err != nil {
		return wb.NewTrustedError(err, http.StatusBadRequest)
	}

	nh, err := toCoreNewHome(ctx, app)
	if err != nil {
		return wb.NewTrustedError(err, http.StatusBadRequest)
	}

	hme, err := h.home.Create(ctx, nh)
	if err != nil {
		return fmt.Errorf("create: hme[%+v]: %w", app, err)
	}

	return wf.Respond(ctx, w, toAppHome(hme), http.StatusCreated)
}

// update updates a home in the system.
func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdateHome
	if err := wf.Decode(r, &app); err != nil {
		return wb.NewTrustedError(err, http.StatusBadRequest)
	}

	uh, err := toCoreUpdateHome(app)
	if err != nil {
		return wb.NewTrustedError(err, http.StatusBadRequest)
	}

	hme := mid.GetHome(ctx)

	updHme, err := h.home.Update(ctx, hme, uh)
	if err != nil {
		return fmt.Errorf("update: homeID[%s] app[%+v]: %w", hme.ID, app, err)
	}

	return wf.Respond(ctx, w, toAppHome(updHme), http.StatusOK)
}

// delete deletes a home from the system.
func (h *handlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	hme := mid.GetHome(ctx)

	if err := h.home.Delete(ctx, hme); err != nil {
		return fmt.Errorf("delete: homeID[%s]: %w", hme.ID, err)
	}

	return wf.Respond(ctx, w, nil, http.StatusNoContent)
}

// query returns a list of homes with paging.
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

	homes, err := h.home.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.home.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return wf.Respond(ctx, w, wb.NewPageDocument(toAppHomes(homes), total, page.Number, page.RowsPerPage), http.StatusOK)
}

// queryByID returns a home by its ID.
func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return wf.Respond(ctx, w, toAppHome(mid.GetHome(ctx)), http.StatusOK)
}

// create adds a new home to the system.
func (h *handlers) ginCreate(c *gin.Context) error {
	var app AppNewHome
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	ctx := c.Request.Context()
	nh, err := toCoreNewHome(ctx, app)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	hme, err := h.home.Create(ctx, nh)
	if err != nil {
		return fmt.Errorf("create: hme[%+v]: %w", app, err)
	}

	c.JSON(http.StatusCreated, toAppHome(hme))
	return nil
}

// update updates a home in the system.
func (h *handlers) ginUpdate(c *gin.Context) error {
	var app AppUpdateHome
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	uh, err := toCoreUpdateHome(app)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	ctx := c.Request.Context()
	hme := mid.GetHome(ctx)

	updHme, err := h.home.Update(ctx, hme, uh)
	if err != nil {
		return fmt.Errorf("update: homeID[%s] app[%+v]: %w", hme.ID, app, err)
	}

	c.JSON(http.StatusOK, toAppHome(updHme))
	return nil
}

func (h *handlers) ginDelete(c *gin.Context) error {
	ctx := c.Request.Context()
	hme := mid.GetHome(ctx)

	if err := h.home.Delete(ctx, hme); err != nil {
		return fmt.Errorf("delete: homeID[%s]: %w", hme.ID, err)
	}

	c.JSON(http.StatusNoContent, nil)
	return nil
}

// query returns a list of homes with paging.
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
	homes, err := h.home.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.home.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	c.JSON(http.StatusOK, wb.NewPageDocument(toAppHomes(homes), total, page.Number, page.RowsPerPage))
	return nil
}

// queryByID returns a home by its ID.
func (h *handlers) ginQueryByID(c *gin.Context) error {
	c.JSON(http.StatusOK, toAppHome(mid.GetHome(c.Request.Context())))
	return nil
}
