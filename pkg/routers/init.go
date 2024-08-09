package routers

import (
	"os"

	"my_project/pkg/context"
	"my_project/pkg/routers/api"
	"my_project/pkg/routers/web"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	r fiber.Router
}

func Routes(app *fiber.App) {
	app.Use(setBaseContext)

	// TODO: Figure out a configuration strategy
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

func setBaseContext(ctx fiber.Ctx) error {
	baseCtx := context.NewBaseContext(ctx)
	baseCtx.SetContext("myCtx", baseCtx)
	return baseCtx.Next()
}
