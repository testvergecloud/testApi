// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/core/crud/user"
	wb "github.com/testvergecloud/testApi/business/web"
	"github.com/testvergecloud/testApi/business/web/auth"
	"github.com/testvergecloud/testApi/business/web/mid"
	"github.com/testvergecloud/testApi/business/web/page"
	"github.com/testvergecloud/testApi/foundation/validate"
	wf "github.com/testvergecloud/testApi/foundation/web"

	"github.com/golang-jwt/jwt/v4"
)

type handlers struct {
	user *user.Core
	auth *auth.Auth
}

func new(user *user.Core, auth *auth.Auth) *handlers {
	return &handlers{
		user: user,
		auth: auth,
	}
}

// create adds a new user to the system.
func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewUser
	if err := wf.Decode(r, &app); err != nil {
		return wb.NewTrustedError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewUser(app)
	if err != nil {
		return wb.NewTrustedError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return wb.NewTrustedError(err, http.StatusConflict)
		}
		return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	return wf.Respond(ctx, w, toAppUser(usr), http.StatusCreated)
}

// update updates a user in the system.
func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdateUser
	if err := wf.Decode(r, &app); err != nil {
		return wb.NewTrustedError(err, http.StatusBadRequest)
	}

	uu, err := toCoreUpdateUser(app)
	if err != nil {
		return wb.NewTrustedError(err, http.StatusBadRequest)
	}

	usr := mid.GetUser(ctx)

	updUsr, err := h.user.Update(ctx, usr, uu)
	if err != nil {
		return fmt.Errorf("update: userID[%s] uu[%+v]: %w", usr.ID, uu, err)
	}

	return wf.Respond(ctx, w, toAppUser(updUsr), http.StatusOK)
}

// delete removes a user from the system.
func (h *handlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	usr := mid.GetUser(ctx)

	if err := h.user.Delete(ctx, mid.GetUser(ctx)); err != nil {
		return fmt.Errorf("delete: userID[%s]: %w", usr.ID, err)
	}

	return wf.Respond(ctx, w, nil, http.StatusNoContent)
}

// query returns a list of users with paging.
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

	users, err := h.user.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.user.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return wf.Respond(ctx, w, wb.NewPageDocument(toAppUsers(users), total, page.Number, page.RowsPerPage), http.StatusOK)
}

// queryByID returns a user by its ID.
func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return wf.Respond(ctx, w, toAppUser(mid.GetUser(ctx)), http.StatusOK)
}

// token provides an API token for the authenticated user.
func (h *handlers) token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	kid := wf.Param(r, "kid")
	if kid == "" {
		return validate.NewFieldsError("kid", errors.New("missing kid"))
	}

	email, pass, ok := r.BasicAuth()
	if !ok {
		return auth.NewAuthError("must provide email and password in Basic auth")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return auth.NewAuthError("invalid email format")
	}

	usr, err := h.user.Authenticate(ctx, *addr, pass)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return wb.NewTrustedError(err, http.StatusNotFound)
		case errors.Is(err, user.ErrAuthenticationFailure):
			return auth.NewAuthError(err.Error())
		default:
			return fmt.Errorf("authenticate: %w", err)
		}
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: usr.Roles,
	}

	token, err := h.auth.GenerateToken(kid, claims)
	if err != nil {
		return fmt.Errorf("generatetoken: %w", err)
	}

	return wf.Respond(ctx, w, toToken(token), http.StatusOK)
}

// create adds a new user to the system.
func (h *handlers) ginCreate(c *gin.Context) error {
	var app AppNewUser
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	nc, err := toCoreNewUser(app)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	usr, err := h.user.Create(c.Request.Context(), nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return err
		}
		return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	c.JSON(http.StatusCreated, toAppUser(usr))
	return nil
}

// update updates a user in the system.
func (h *handlers) ginUpdate(c *gin.Context) error {
	var app AppUpdateUser
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	uu, err := toCoreUpdateUser(app)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	ctx := c.Request.Context()
	usr := mid.GetUser(ctx)

	updUsr, err := h.user.Update(ctx, usr, uu)
	if err != nil {
		return fmt.Errorf("update: userID[%s] uu[%+v]: %w", usr.ID, uu, err)
	}

	c.JSON(http.StatusOK, toAppUser(updUsr))
	return nil
}

// delete removes a user from the system.
func (h *handlers) ginDelete(c *gin.Context) error {
	ctx := c.Request.Context()
	usr := mid.GetUser(ctx)

	if err := h.user.Delete(ctx, mid.GetUser(ctx)); err != nil {
		return fmt.Errorf("delete: userID[%s]: %w", usr.ID, err)
	}

	c.JSON(http.StatusNoContent, nil)
	return nil
}

// query returns a list of users with paging.
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
	users, err := h.user.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.user.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	c.JSON(http.StatusOK, wb.NewPageDocument(toAppUsers(users), total, page.Number, page.RowsPerPage))
	return nil
}

// queryByID returns a user by its ID.
func (h *handlers) ginQueryByID(c *gin.Context) error {
	c.JSON(http.StatusOK, toAppUser(mid.GetUser(c.Request.Context())))
	return nil
}

// token provides an API token for the authenticated user.
func (h *handlers) ginToken(c *gin.Context) error {
	kid := c.Param("kid")
	if kid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing kid"})
		return validate.NewFieldsError("kid", errors.New("missing kid"))
	}

	email, pass, ok := c.Request.BasicAuth()
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "must provide email and password in Basic auth"})
		return auth.NewAuthError("must provide email and password in Basic auth")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return auth.NewAuthError("invalid email format")
	}

	usr, err := h.user.Authenticate(c.Request.Context(), *addr, pass)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return wb.NewTrustedError(err, http.StatusNotFound)
		case errors.Is(err, user.ErrAuthenticationFailure):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return auth.NewAuthError(err.Error())
		default:
			return fmt.Errorf("authenticate: %w", err)
		}
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: usr.Roles,
	}

	token, err := h.auth.GenerateToken(kid, claims)
	if err != nil {
		return fmt.Errorf("generatetoken: %w", err)
	}

	c.JSON(http.StatusOK, toToken(token))
	return nil
}
