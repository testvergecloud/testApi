// Package checkgrp maintains the group of handlers for health checking.
package checkgrp

import (
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/data/sqldb"
	"github.com/testvergecloud/testApi/foundation/logger"

	"github.com/jmoiron/sqlx"
)

type handlers struct {
	build string
	log   *logger.Logger
	db    *sqlx.DB
}

func new(build string, log *logger.Logger, db *sqlx.DB) *handlers {
	return &handlers{
		build: build,
		db:    db,
		log:   log,
	}
}

func (h *handlers) Readiness(c *gin.Context) error {
	status := "ok"
	statusCode := http.StatusOK
	ctx := c.Request.Context()
	if err := sqldb.StatusCheck(ctx, h.db); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
		h.log.Info(ctx, "readiness failure", "status", status)
		return c.AbortWithError(statusCode, err)
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	c.JSON(statusCode, data)

	return nil
}

func (h *handlers) Liveness(c *gin.Context) error {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := gin.H{
		"status":     "up",
		"build":      "", // You can set this if needed.
		"host":       host,
		"name":       os.Getenv("KUBERNETES_NAME"),
		"podIP":      os.Getenv("KUBERNETES_POD_IP"),
		"node":       os.Getenv("KUBERNETES_NODE_NAME"),
		"namespace":  os.Getenv("KUBERNETES_NAMESPACE"),
		"GOMAXPROCS": runtime.GOMAXPROCS(0),
	}

	c.JSON(http.StatusOK, data)

	return nil
}
