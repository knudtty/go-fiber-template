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

func Routes(app *fiber.App, state *state.AppState) {
	app.Use(middleware.JWTParser())
	app.Use(setBaseContext)

	switch os.Getenv("ROUTES_AVAILABLE") {
	case "api":
		api.Routes(app, state)
		break
	case "web":
		web.Routes(app, state)
		break
	default:
		apiGroup := app.Group("/api/v1")
		api.Routes(apiGroup, state)
		web.Routes(app, state)
	}
}

func setBaseContext(c *fiber.Ctx) error {
	baseCtx := ctx.NewBaseContext(c)
	baseCtx.SetContext("myCtx", baseCtx)
	return baseCtx.Next()
}
