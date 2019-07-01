package middlewares

import (
	"context"
	"github.com/labstack/echo/v4"
)

const ContextFirebaseName = "ContextFirebase"

func ContextFireBase(client AppAdapterInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			c.SetRequest(req.WithContext(context.WithValue(req.Context(), ContextFirebaseName, client)))
			return next(c)
		}
	}
}
