package main

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"my_project/routers"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	routers.Routes(app)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000", fiber.ListenConfig{
		EnablePrefork:     true,
		EnablePrintRoutes: true,
	}))

    // TODO: global config
    // TODO: yaml config file
    // TODO: db TODO app
    // TODO: api TODO routes
    // TODO: auto migrations
    // TODO: oauth2 auth
    // TODO: structured logging
    // TODO: users, sessions, and roles in db
    // TODO: JWT or sessions
}
