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
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	if len(os.Getenv("DOCKER_ENV")) == 0 {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal(err)
		}
	}
	db, err := database.GetDbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.UsersStore.DB.Close()

	// Initialize a new Fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		BodyLimit:         4 * 1024 * 1024, // this sets limit to 4MB
		//Prefork:           true,
		//EnablePrintRoutes: true,
	})

	// Initialize global middlewares
	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelDefault,
	}))
	app.Use(logger.New())
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_PASSWORD"),
	}))
	if os.Getenv("DEVELOPMENT") != "" {
		app.Use(pprof.New())
	}

	// Static file server
	app.Static("/", "./public", fiber.Static{Compress: true})
	app.Static("/", "./deps", fiber.Static{Compress: true})

	routers.Routes(app, &state.AppState{
		DB: db,
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
