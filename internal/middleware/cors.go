package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	allowMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	allowHeaders = "Accept, Authorization, Content-Type, X-Request-ID"
	maxAge       = "86400" // 24h
)

func CORS(allowOrigins []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			origin := req.Header.Get(echo.HeaderOrigin)

			res.Header().Add(echo.HeaderVary, echo.HeaderOrigin)

			if origin == "" {
				return next(c)
			}

			if !isOriginAllowed(origin, allowOrigins) {
				if req.Method == http.MethodOptions {
					return c.NoContent(http.StatusForbidden)
				}
				return next(c)
			}

			if req.Method == http.MethodOptions && req.Header.Get(echo.HeaderAccessControlRequestMethod) != "" {
				res.Header().Set(echo.HeaderAccessControlAllowOrigin, origin)
				res.Header().Set(echo.HeaderAccessControlAllowMethods, allowMethods)
				res.Header().Set(echo.HeaderAccessControlAllowHeaders, allowHeaders)
				res.Header().Set(echo.HeaderAccessControlAllowCredentials, "true")
				res.Header().Set(echo.HeaderAccessControlMaxAge, maxAge)
				return c.NoContent(http.StatusNoContent)
			}

			res.Header().Set(echo.HeaderAccessControlAllowOrigin, origin)
			res.Header().Set(echo.HeaderAccessControlAllowCredentials, "true")
			res.Header().Set(echo.HeaderAccessControlExposeHeaders, "Content-Length, Content-Type, X-Request-ID")

			return next(c)
		}
	}
}

func isOriginAllowed(origin string, allowOrigins []string) bool {
	for _, allowOrigin := range allowOrigins {
		if allowOrigin == "*" || strings.EqualFold(allowOrigin, origin) {
			return true
		}
	}
	return false
}
