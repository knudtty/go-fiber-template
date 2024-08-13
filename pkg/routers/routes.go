package routers

import (
	"os"

	"my_project/app/state"
	ctx "my_project/pkg/context"
	"my_project/pkg/middleware"
	"my_project/pkg/routers/api"
	"my_project/pkg/routers/web"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	r fiber.Router
}

func Routes(app *fiber.App, as *state.AppState) {
	app.Use(middleware.JWTParser())
	app.Use(setBaseContext(as))

	switch os.Getenv("ROUTES_AVAILABLE") {
	case "api":
		api.Routes(app)
		break
	case "web":
		web.Routes(app)
		break
	default:
		apiGroup := app.Group("/api/v1")
		api.Routes(apiGroup)
		web.Routes(app)
	}
}

func setBaseContext(as *state.AppState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		baseCtx := ctx.NewBaseContext(c)
		baseCtx.AppState = as
		baseCtx.SetContext("myCtx", baseCtx)
		return baseCtx.Next()
	}
}
