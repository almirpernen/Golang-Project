package main

import (
	"log"
	"os"

	"github.com/almirpernen/database"
	"github.com/almirpernen/handlers"
	"github.com/gofiber/fiber/v2"
)

func envPortOr(defaultPort string) string {
	if port, exists := os.LookupEnv("PORT"); exists {
		return ":" + port
	}
	return ":" + defaultPort
}

func main() {
	database.ConnectDb()

	app := fiber.New()

	setupRoutes(app)

	port := envPortOr("3000")

	err := app.Listen(port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func setupRoutes(app *fiber.App) {

	app.Post("/signup", handlers.Signup)
	app.Post("/signin", handlers.Signin)
	app.Post("/refresh", handlers.JWTMiddleware, handlers.RefreshToken)

	app.Post("/bortzhurnal", handlers.JWTMiddleware, handlers.CreatePost)
	app.Get("/bortzhurnal", handlers.ListPosts)
	app.Get("/bortzhurnal/:id", handlers.GetPost)
	app.Delete("/bortzhurnal/:id", handlers.JWTMiddleware, handlers.DeletePost)
	app.Put("/bortzhurnal/:id", handlers.JWTMiddleware, handlers.UpdatePost)
	app.Post("/bortzhurnal/:id/like", handlers.JWTMiddleware, handlers.LikePost)
	app.Post("/bortzhurnal/:id/unlike", handlers.JWTMiddleware, handlers.UnlikePost)

	app.Get("/users", handlers.JWTMiddleware, handlers.ListUsers)
	app.Get("/users/:id", handlers.JWTMiddleware, handlers.GetUsers)
	app.Delete("/users/:id", handlers.JWTMiddleware, handlers.DeleteUser)
	app.Post("/users/:id/follow", handlers.JWTMiddleware, handlers.FollowUser)
	app.Post("/users/:id/unfollow", handlers.JWTMiddleware, handlers.UnfollowUser)

	app.Post("/feedback/:id", handlers.JWTMiddleware, handlers.CreateComment)
	app.Get("/feedback", handlers.ListComments)
	app.Get("/feedback/:id", handlers.GetComment)
	app.Put("/feedback/:id", handlers.JWTMiddleware, handlers.UpdateComment)
	app.Delete("/feedback/:id", handlers.JWTMiddleware, handlers.DeleteComment)
	app.Post("/feedback/:id/like", handlers.JWTMiddleware, handlers.LikeComment)
	app.Post("/feedback/:id/unlike", handlers.JWTMiddleware, handlers.UnlikeComment)

	app.Get("/test", handlers.TestApi)
}
