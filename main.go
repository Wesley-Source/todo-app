package main

import (
	"log"
	"os"
	"todo-app/app/middleware"
	"todo-app/app/routes"
	"todo-app/config/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize and connect to the databases
	database.ConnectDatabase()
	middleware.ConnectSessionsDB()

	// Load environment variables
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatalln(err)
	}

	engine := html.New("./app/views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Static files directory for icons, scripts, thumbnails etc.
	app.Static("/", "./app/public")

	// Routes configuration
	app.Get("/", routes.Index)

	app.Get("/login", middleware.AuthMiddleware, routes.LoginGet)
	app.Post("/login", middleware.AuthMiddleware, routes.LoginPost)

	app.Post("/register", middleware.AuthMiddleware, routes.RegisterPost)
	app.Get("/register", middleware.AuthMiddleware, routes.RegisterGet)

	app.Get("/logout", middleware.AuthMiddleware, routes.LogoutGet)
	app.Get("/todo", middleware.AuthMiddleware, routes.TodoGet)

	app.Post("/list_add", middleware.AuthMiddleware, routes.ListAddPost)
	app.Post("/task_add", middleware.AuthMiddleware, routes.TaskAddPost)

	log.Fatalln(app.Listen(os.Getenv("PORT")))
}
