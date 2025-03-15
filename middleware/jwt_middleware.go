package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// secret key for the jwt
var secretKey = []byte("secret_key")

func JWTMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {
		// get the token from the header
		tokenString := c.Get("Authorization")

		// check if token is empty
		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{"Error": "Missing or invalid token"})
		}
		// parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			return c.Status(401).JSON(fiber.Map{"Error": "Invalid token"})
		}

		// set user id in context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"Error": "Invalid token"})
		}
		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])
		return c.Next()
	}
}
