package main

import (
	"log"

	"my_project/pkg/routers"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
        Prefork: true,
        EnablePrintRoutes: true,
	})

	routers.Routes(app)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
