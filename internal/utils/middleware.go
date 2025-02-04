package utils

import (
	"net/http"
	"net/url"
	"os"
	"p2pbot/internal/JWTConfig"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type CustomClaims struct {
	Scope string `json:"scope"`
	Email string `json:"email"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	if c.Email == "" {
		return jwtmiddleware.ErrJWTInvalid
	}
	return nil
}

func CheckJWT(next echo.HandlerFunc) echo.HandlerFunc {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing issuer URL")
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating JWT validator")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Error().Err(err).Msg("Failed to validate JWT")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}
	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(ctx echo.Context) error {
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			ctx.SetRequest(r)
			next(ctx)
		}

		middleware.CheckJWT(handler).ServeHTTP(ctx.Response(), ctx.Request())

		if encounteredError {
			ctx.JSON(
				http.StatusUnauthorized,
				map[string]string{"message": "JWT is invalid."},
			)
		}

		return nil
	}
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// log the request
		log.Info().Fields(map[string]interface{}{
			"method":     c.Request().Method,
			"uri":        c.Request().URL.Path,
			"user_agent": c.Request().UserAgent(),
			"query":      c.Request().URL.RawQuery,
			"client_ip":  c.RealIP(),
		}).Msg("Request")

		// call the next middleware/handler
		err := next(c)
		if err != nil {
			log.Error().Fields(map[string]interface{}{
				"error": err.Error(),
			}).Msg("Response")
			return err
		}

		log.Info().Fields(map[string]interface{}{
			"status":    c.Response().Status,
			"client_ip": c.RealIP(),
			"size":      c.Response().Size,
		}).Msg("Response")

		return nil
	}
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		if !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		claims := token.Claims.(*JWTConfig.JWTCustomClaims)

		c.Set("email", claims.Email)

		log.Info().Fields(map[string]interface{}{
			"email":     claims.Email,
			"client_ip": c.RealIP(),
		}).Msg("User authenticated")

		return next(c)
	}
}

func ExtractEmail(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		claims, ok := ctx.Request().Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		if !ok {
			ctx.JSON(
				http.StatusInternalServerError,
				map[string]string{"message": "Failed to get JWT claims."},
			)
			return nil
		}

		customClaims, ok := claims.CustomClaims.(*CustomClaims)
		if !ok {
			ctx.JSON(
				http.StatusInternalServerError,
				map[string]string{"message": "Failed to cast custom claims."},
			)
			return nil
		}

		if customClaims.Email == "" {
			ctx.JSON(
				http.StatusInternalServerError,
				map[string]string{"message": "Failed to get JWT claims"},
			)
			return nil
		}
		ctx.Set("email", customClaims.Email)

		log.Info().Fields(map[string]interface{}{
			"email":     customClaims.Email,
			"client_ip": ctx.RealIP(),
		}).Msg("User authenticated")

		return next(ctx)
	}
}
