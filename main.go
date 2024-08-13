package main

import (
	"log"
	"os"

	"my_project/app/state"
	"my_project/pkg/routers"
	"my_project/platform/database"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		Prefork:           true,
		EnablePrintRoutes: true,
	})

	db, _ := database.GetDbConnection()

	// Initialize global middlewares
	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))
	app.Use(logger.New())
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_PASSWORD"),
	}))

	// Static file server
	app.Static("/", "./public", fiber.Static{Compress: true})

	routers.Routes(app, &state.AppState{
		DB: db,
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
