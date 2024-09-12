package utils

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/golang-jwt/jwt/v5"
    "p2pbot/internal/JWTConfig"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
 return func(c echo.Context) error {
  // log the request
  Logger.LogInfo().Fields(map[string]interface{}{
   "method": c.Request().Method,
   "uri":    c.Request().URL.Path,
   "user_agent":    c.Request().UserAgent(),
   "query":  c.Request().URL.RawQuery,
   "client_ip": c.RealIP(),
  }).Msg("Request")

  // call the next middleware/handler
  err := next(c)
  if err != nil {
   Logger.LogError().Fields(map[string]interface{}{
    "error": err.Error(),
   }).Msg("Response")
   return err
  }

  Logger.LogInfo().Fields(map[string]interface{}{
      "status": c.Response().Status,
      "client_ip":     c.RealIP(),
      "size":          c.Response().Size,
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
        
        Logger.LogInfo().Fields(map[string]interface{}{
            "email": claims.Email,
            "client_ip": c.RealIP(),
        }).Msg("User authenticated")
    
        return next(c)
    }
}
