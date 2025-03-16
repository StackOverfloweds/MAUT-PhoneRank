package middleware

import (
	"github.com/StackOverfloweds/MAUT-PhoneRank/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

/*
JWTMiddleware - Middleware for JWT authentication.
This function checks the Authorization header for a valid JWT token.
If the token is missing, invalid, or expired, it returns a 401 Unauthorized response.
If valid, it extracts the user ID and role from the token and stores them in Fiber's context.
This allows authenticated routes to access user information.
*/
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{"Error": "Missing or invalid token"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return helpers.GetJWTSecret, nil
		})

		if err != nil {
			return c.Status(401).JSON(fiber.Map{"Error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"Error": "Invalid token"})
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}
