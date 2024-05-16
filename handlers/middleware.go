package handlers

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	"strings"
)

type CustomClaims struct {
	jwt.StandardClaims
	UserID uint `json:"userID"`
}

func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "No authorization header provided"})
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid authorization header format"})
	}

	tokenString := headerParts[1]
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return jwtSecret, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid or expired JWT token"})
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		fmt.Printf("JWT parsed UserID: %v\n", claims.UserID)
		if claims.UserID > 0 {
			c.Locals("userID", claims.UserID)
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "User ID in JWT is invalid"})
		}
		return c.Next()
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid or expired JWT token"})
	}
}
