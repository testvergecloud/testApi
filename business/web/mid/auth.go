package mid

import (
	"context"
	"errors"
	"net/http"

	wb "github.com/testvergecloud/testApi/business/web"
	"github.com/testvergecloud/testApi/business/web/auth"
	wf "github.com/testvergecloud/testApi/foundation/web"

	"github.com/google/uuid"
)

// Set of error variables for handling auth errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(a *auth.Auth) wf.MidHandler {
	m := func(handler wf.Handler) wf.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims, err := a.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				return auth.NewAuthError("authenticate: failed: %s", err)
			}

			if claims.Subject == "" {
				return auth.NewAuthError("authorize: you are not authorized for that action, no claims")
			}

			subjectID, err := uuid.Parse(claims.Subject)
			if err != nil {
				return wb.NewTrustedError(ErrInvalidID, http.StatusBadRequest)
			}

			ctx = setUserID(ctx, subjectID)
			ctx = setClaims(ctx, claims)

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// Authorize executes the specified role and does not extract any domain data.
func Authorize(a *auth.Auth, rule string) wf.MidHandler {
	m := func(handler wf.Handler) wf.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims := getClaims(ctx)
			if err := a.Authorize(ctx, claims, uuid.UUID{}, rule); err != nil {
				return auth.NewAuthError("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
