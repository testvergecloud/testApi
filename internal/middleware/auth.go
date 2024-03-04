package middleware

import (
	"context"
	"go-starter/config"
	"log"
	"net/http"
	"net/url"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	jwtvalidator "github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

type CustomClaims struct {
	Permissions []string `json:"permissions"`
}

// Validate errors out if `ShouldReject` is true.
func (c *CustomClaims) Validate(ctx context.Context) error {
	return nil
}

var (

	//config is loaded manually because we don't have a constructor here
	cfg, _ = config.LoadConfig(".")

	// The issuer of our token.
	issuer = cfg.Auth0Domain

	// The audience of our token.
	audience = []string{cfg.Auth0Audience}

	//In case of HS256 algorithm, signingKey must be copied from auth0 panel and be pasted here
	// The signing key for the token.
	// signingKey = []byte("secret")

	// keyFunc = func(ctx context.Context) (interface{}, error) {
	// 	return signingKey, nil
	// }

	customClaims = func() jwtvalidator.CustomClaims {
		return &CustomClaims{}
	}
)

func CheckJWT() gin.HandlerFunc {
	issuerURL, err := url.Parse(issuer)
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := jwtvalidator.New(
		provider.KeyFunc,
		jwtvalidator.RS256,
		issuer,
		audience,
		jwtvalidator.WithCustomClaims(customClaims),
		jwtvalidator.WithAllowedClockSkew(30*time.Second),
	)
	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(ctx *gin.Context) {
		encounteredError := true
		//converting gin handler to http handler
		var nextHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			ctx.Request = r
			ctx.Next()
		}

		middleware.CheckJWT(nextHandler).ServeHTTP(ctx.Writer, ctx.Request)

		if encounteredError {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				map[string]string{"message": "JWT is invalid."},
			)
		}
	}
}

func (c CustomClaims) HasPermissions(expectedPermission string) bool {
	for _, b := range c.Permissions {
		if b == expectedPermission {
			return true
		}
	}
	return false
}

func ValidatePermissions(expectedPermission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Context().Value(jwtmiddleware.ContextKey{}).(*jwtvalidator.ValidatedClaims)
		claims := token.CustomClaims.(*CustomClaims)
		if claims.HasPermissions(expectedPermission) {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	}
}
