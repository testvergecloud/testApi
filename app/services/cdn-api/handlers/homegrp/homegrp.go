// Package homegrp maintains the group of handlers for home access.
package homegrp

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/core/crud/home"
	wb "github.com/testvergecloud/testApi/business/web"
	"github.com/testvergecloud/testApi/business/web/mid"
	"github.com/testvergecloud/testApi/business/web/page"
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
func (h *handlers) create(c *gin.Context) error {
	var app AppNewHome
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	ctx := c.Request.Context()
	nh, err := toCoreNewHome(c, app)
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
func (h *handlers) update(c *gin.Context) error {
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
	hme := mid.GetHome(c)

	updHme, err := h.home.Update(ctx, hme, uh)
	if err != nil {
		return fmt.Errorf("update: homeID[%s] app[%+v]: %w", hme.ID, app, err)
	}

	c.JSON(http.StatusOK, toAppHome(updHme))
	return nil
}

func (h *handlers) delete(c *gin.Context) error {
	ctx := c.Request.Context()
	hme := mid.GetHome(c)

	if err := h.home.Delete(ctx, hme); err != nil {
		return fmt.Errorf("delete: homeID[%s]: %w", hme.ID, err)
	}

	c.JSON(http.StatusNoContent, nil)
	return nil
}

// query returns a list of homes with paging.
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
func (h *handlers) queryByID(c *gin.Context) error {
	c.JSON(http.StatusOK, toAppHome(mid.GetHome(c)))
	return nil
}
