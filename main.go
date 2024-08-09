package main

import (
	"log"

	"my_project/pkg/routers"
	"my_project/platform/cache"
	"my_project/platform/database"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	database.Init()
	cache.Init()

	routers.Routes(app)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000", fiber.ListenConfig{
		EnablePrefork:     true,
		EnablePrintRoutes: true,
	}))
}
