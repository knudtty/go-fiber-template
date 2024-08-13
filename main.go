package main

import (
	"log"

	"my_project/app/state"
	"my_project/pkg/routers"
	"my_project/platform/database"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
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

	//// Initialize global middlewares
	//app.Use(recover.New())
	//app.Use(compress.New(compress.Config{
	//	Level: compress.LevelBestSpeed, // 1
	//}))
	//app.Use(logger.New())
	//app.Use(encryptcookie.New(encryptcookie.Config{
	//	Key: cookiePassword,
	//}))

	//// Static file server
	//app.Static("/", "./public", fiber.Static{Compress: true})
	//app.Static("/", "./public/assets/images", fiber.Static{Compress: true})

	routers.Routes(app, &state.AppState{
		DB: db,
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
