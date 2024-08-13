package web

import (
	"log"

	"my_project/app/controllers"
	"my_project/app/state"
	ctx "my_project/pkg/context"
	"my_project/pkg/middleware"
	"my_project/templates"

	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router, as *state.AppState) {
	app.Use(setWebContext)
	publicRoutes(app, as)
	privateRoutes(app, as)
}

func WrapWeb(f func(*ctx.WebCtx, *state.AppState) error, state *state.AppState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		webCtx, ok := c.UserContext().Value("myCtx").(*ctx.WebCtx)
		if !ok {
			log.Fatal("webCtx not found")
		}
		return f(webCtx, state)
	}
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

func publicRoutes(app fiber.Router, as *state.AppState) {
	app.Get("/", WrapWeb(func(c *ctx.WebCtx, as *state.AppState) error {
		return c.Render(templates.Home())
	}, as))

	app.Get("/login", WrapWeb(func(c *ctx.WebCtx, as *state.AppState) error {
		return c.Render(templates.Login())
	}, as))
	authGroup := app.Group("/auth")
	authGroup.Get("/session", WrapWeb(controllers.CreateSession, as))
	authGroup.Get("/redirect", WrapWeb(controllers.AuthRedirectFromProvider, as))
}

func privateRoutes(app fiber.Router, as *state.AppState) {
	app.Use(WrapWeb(middleware.AuthenticatedUser, as))
	app.Get("/userinfo", WrapWeb(controllers.UserInfo, as))
}
