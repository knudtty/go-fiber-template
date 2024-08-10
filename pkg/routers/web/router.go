package web

import (
	"log"

	"my_project/app/controllers"
	ctx "my_project/pkg/context"
	"my_project/pkg/middleware"
	"my_project/templates"

	"github.com/gofiber/fiber/v2"
)

func WrapWeb(f func(*ctx.WebCtx) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		webCtx, ok := c.UserContext().Value("myCtx").(*ctx.WebCtx)
		if !ok {
			log.Fatal("webCtx not found")
		}
		return f(webCtx)
	}
}

func Routes(app fiber.Router) {
	app.Use(setWebContext)
	publicRoutes(app)
	privateRoutes(app)
}

func setWebContext(c *fiber.Ctx) error {
	baseCtx, ok := c.UserContext().Value("myCtx").(*ctx.Base)
	if !ok {
		log.Fatal("myCtx base not found for web")
	}
	webCtx := ctx.NewWebContext(baseCtx)
	baseCtx.SetContext("myCtx", webCtx)
	return baseCtx.Next()
}

func publicRoutes(app fiber.Router) {
	app.Get("/", WrapWeb(func(c *ctx.WebCtx) error {
		return c.Render(templates.Home())
	}))

	app.Get("/login", WrapWeb(func(c *ctx.WebCtx) error {
		return c.Render(templates.Login())
	}))
	authGroup := app.Group("/auth")
	authGroup.Get("/session", WrapWeb(controllers.CreateSession))
    authGroup.Get("/redirect", WrapWeb(controllers.AuthRedirectFromProvider))
}

func privateRoutes(app fiber.Router) {
    app.Use(middleware.JWTProtected())
}
