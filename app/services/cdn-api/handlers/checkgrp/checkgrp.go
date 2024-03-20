// Package checkgrp maintains the group of handlers for health checking.
package checkgrp

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/data/sqldb"
	"github.com/testvergecloud/testApi/foundation/logger"
	"github.com/testvergecloud/testApi/foundation/web"

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

// readiness checks if the database is ready and if not will return a 500 status.
// Do not respond by just returning an error because further up in the call
// stack it will interpret that as a non-trusted error.
func (h *handlers) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	status := "ok"
	statusCode := http.StatusOK
	if err := sqldb.StatusCheck(ctx, h.db); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
		h.log.Info(ctx, "readiness failure", "status", status)
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	return web.Respond(ctx, w, data, statusCode)
}

// liveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (h *handlers) liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status     string `json:"status,omitempty"`
		Build      string `json:"build,omitempty"`
		Host       string `json:"host,omitempty"`
		Name       string `json:"name,omitempty"`
		PodIP      string `json:"podIP,omitempty"`
		Node       string `json:"node,omitempty"`
		Namespace  string `json:"namespace,omitempty"`
		GOMAXPROCS int    `json:"GOMAXPROCS,omitempty"`
	}{
		Status: "up",
		// Build:      h.build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: runtime.GOMAXPROCS(0),
	}

	// This handler provides a free timer loop.

	return web.Respond(ctx, w, data, http.StatusOK)
}

func (h *handlers) ginReadiness(c *gin.Context) error {
	status := "ok"
	statusCode := http.StatusOK
	if err := sqldb.StatusCheck(c, h.db); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
		h.log.Info(c, "readiness failure", "status", status)
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

func (h *handlers) ginLiveness(c *gin.Context) error {
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
