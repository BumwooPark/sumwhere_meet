package middlewares

import (
	"context"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

const ContextDBName = "DB"

func ContextDB(service string, db *xorm.Engine) echo.MiddlewareFunc {
	db.ShowExecTime()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			session := db.NewSession()
			defer session.Close()

			func(session interface{}, ctx context.Context) {
				if s, ok := session.(interface{ SetContext(context.Context) }); ok {
					s.SetContext(ctx)
				}
			}(session, req.Context())

			c.SetRequest(req.WithContext(context.WithValue(req.Context(), ContextDBName, session)))

			switch req.Method {
			case "POST", "PUT", "DELETE", "PATCH":
				if err := session.Begin(); err != nil {
					log.Println(err)
				}
				if err := next(c); err != nil {
					session.Rollback()
					return err
				}
				if c.Response().Status >= 500 {
					session.Rollback()
					return nil
				}
				if err := session.Commit(); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			default:
				return next(c)
			}

			return nil
		}
	}
}
