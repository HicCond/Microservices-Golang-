package main

import (
    "github.com/gofiber/fiber/v2"
    "microservices/handlers"
    "microservices/database"
    jwtware "github.com/gofiber/contrib/jwt"
)

func main() {
    // Initialize Fiber app
    app := fiber.New()

    // Initialize MySQL database
    database.InitDB()

    // Define routes
    app.Post("/api/register", handlers.RegisterHandler)
    app.Post("/api/auth", handlers.AuthHandler)

    // JWT middleware with custom error handler
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{
            Key: handlers.JwtSecret,
        },
		ErrorHandler: handlers.JwtErrorHandler,
	}))

	// Route to check JWT token
	app.Get("/api/auth/check", handlers.JWTTokenSuccess)

    // Start the server
    app.Listen(":8081")
}
