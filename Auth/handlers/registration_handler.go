package handlers

import (
    "github.com/gofiber/fiber/v2"
    "microservices/database"
)

// RegisterRequest defines the structure of the registration request body
type RegisterRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// RegisterHandler handles registration requests
func RegisterHandler(c *fiber.Ctx) error {
    // Parse request body
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid request body",
        })
    }

    // Validate request parameters
    if req.Username == "" || req.Password == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Username and password are required",
        })
    }

    // Save user to database
    if err := database.CreateUser(req.Username, req.Password); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to register user",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Registered successfully",
    })
}

// AuthHandler handles authentication requests
func AuthHandler(c *fiber.Ctx) error {
    // Parse request body
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid request body",
        })
    }

    // Check if username and password match
    if !database.ValidateUser(req.Username, req.Password) {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Invalid username or password",
        })
    }

    // Generate JWT token
    token, err := GenerateJWT(req.Username)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to generate JWT token",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "token": token,
    })
}

func JWTTokenSuccess(c *fiber.Ctx) error {
	// If JWT middleware succeeds, this means the token is valid
	return c.SendString("Authorized")
}

func JwtErrorHandler(c *fiber.Ctx, err error) error {
	if err != nil {
		// If JWT middleware fails, return custom error message
		return c.Status(fiber.StatusUnauthorized).SendString("Failed Authorization: Check jwt token")
	}
	return nil
}
