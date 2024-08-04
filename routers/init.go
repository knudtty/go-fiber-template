package routers

import (
	"os"

	"github.com/caritaspay/caritas/routers/api"
	"github.com/caritaspay/caritas/routers/context"
	"github.com/caritaspay/caritas/routers/web"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	r fiber.Router
}

func Routes(app *fiber.App) {
	app.Use(setBaseContext)

	// TODO: Figure out a configuration strategy
	if len(os.Getenv("API_ONLY")) > 0 {
		api.Routes(app)
	} else if len(os.Getenv("WEBSITE_ONLY")) > 0 {
		web.Routes(app)
	} else {
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
