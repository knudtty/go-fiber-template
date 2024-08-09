package web

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"my_project/pkg/context"
	"my_project/templates"
)

func WrapWeb(f func(*context.WebCtx) error) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		webCtx, ok := ctx.UserContext().Value("myCtx").(*context.WebCtx)
		if !ok {
			log.Fatal("webCtx not found")
		}
		return f(webCtx)
	}
}

func Routes(app fiber.Router) {
	app.Use(setWebContext)
	app.Get("/", WrapWeb(func(c *context.WebCtx) error {
		return c.Render(templates.Home())
	}))
}

func setWebContext(ctx fiber.Ctx) error {
	baseCtx, ok := ctx.UserContext().Value("myCtx").(*context.Base)
	if !ok {
		log.Fatal("myCtx base not found for web")
	}
	webCtx := context.NewWebContext(baseCtx)
	baseCtx.SetContext("myCtx", webCtx)
	return baseCtx.Next()
}
