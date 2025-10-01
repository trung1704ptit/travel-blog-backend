package middleware

import "github.com/labstack/echo/v4"

// CORS will handle the CORS middleware
func CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		origin := c.Request().Header.Get("Origin")

		// Allow specific origins for development
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:5173",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:3001",
			"http://127.0.0.1:5173",
			"http://127.0.0.1:8080",
		}

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Response().Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// For development, allow all origins
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		}

		// Set other CORS headers
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		c.Response().Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Request().Method == "OPTIONS" {
			return c.NoContent(204)
		}

		return next(c)
	}
}
