// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
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
func (h *handlers) create(c *gin.Context) error {
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
func (h *handlers) update(c *gin.Context) error {
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
	usr := mid.GetUser(c)

	updUsr, err := h.user.Update(ctx, usr, uu)
	if err != nil {
		return fmt.Errorf("update: userID[%s] uu[%+v]: %w", usr.ID, uu, err)
	}

	c.JSON(http.StatusOK, toAppUser(updUsr))
	return nil
}

// delete removes a user from the system.
func (h *handlers) delete(c *gin.Context) error {
	ctx := c.Request.Context()
	usr := mid.GetUser(c)

	if err := h.user.Delete(ctx, mid.GetUser(c)); err != nil {
		return fmt.Errorf("delete: userID[%s]: %w", usr.ID, err)
	}

	c.JSON(http.StatusNoContent, nil)
	return nil
}

// query returns a list of users with paging.
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
func (h *handlers) queryByID(c *gin.Context) error {
	c.JSON(http.StatusOK, toAppUser(mid.GetUser(c)))
	return nil
}

// token provides an API token for the authenticated user.
func (h *handlers) token(c *gin.Context) error {
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
