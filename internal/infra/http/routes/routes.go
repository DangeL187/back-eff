package routes

import (
	echoV4 "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v5"
	"github.com/swaggo/echo-swagger"

	"back/internal/app"
	_ "back/internal/docs"
)

func ConvertEchoV4HandlerToV5(h echoV4.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		req := c.Request()
		res := c.Response()
		e := echoV4.New()
		v4Ctx := e.NewContext(req, res)
		return h(v4Ctx)
	}
}

func RegisterRoutes(e *echo.Echo, app *app.App) {
	subs := e.Group("/subscriptions")

	subs.POST("", app.SubsHandler.Create)
	subs.GET("/:id", app.SubsHandler.Get)
	subs.PUT("/:id", app.SubsHandler.Update)
	subs.DELETE("/:id", app.SubsHandler.Delete)
	subs.GET("", app.SubsHandler.List)
	subs.GET("/costs", app.SubsHandler.GetSubscriptionsCostForPeriod)

	e.GET("/swagger/*", ConvertEchoV4HandlerToV5(echoSwagger.WrapHandler))
}
