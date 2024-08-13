package api

import (
	"log"

	"my_project/pkg/context"

	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	app.Use(setApiContext)

	app.Get("/", WrapApi(func(ac *context.ApiCtx) error {
		m := fiber.Map{}
		m["message"] = "Hello, World!"
		return ac.JSON(m)
	}))
}

func WrapApi(f func(*context.ApiCtx) error) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		apiCtx, ok := ctx.UserContext().Value("myCtx").(*context.ApiCtx)
		if !ok {
			log.Fatal("apiCtx not found")
		}
		return f(apiCtx)
	}
}

func setApiContext(ctx *fiber.Ctx) error {
	baseCtx, ok := ctx.UserContext().Value("myCtx").(*context.Base)
	if !ok {
		log.Fatalf("myCtx base not found for api on route %s", ctx.Route().Path)
	}
	apiCtx := context.NewApiContext(baseCtx)
	baseCtx.SetContext("myCtx", apiCtx)
	return baseCtx.Next()
}
