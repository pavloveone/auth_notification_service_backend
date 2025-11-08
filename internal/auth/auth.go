package auth

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey []byte

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("Jwt secret isn't set in .env")
	}
	jwtKey = []byte(secret)
}

func HashPassword(pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}

func CheckPassHash(pass, hashedPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
	return err == nil
}

func GenerateTokens(id int) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"id":     fmt.Sprintf("%d", id),
		"expire": time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenStr, err := accessToken.SignedString(jwtKey)
	if err != nil {
		log.Printf("failed to sign access token: %v", err)
		return "", "", err
	}
	refreshTokenClaims := jwt.MapClaims{
		"id":     fmt.Sprintf("%d", id),
		"expire": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenStr, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		log.Printf("failed to sign refresh token: %v", err)
		return "", "", err
	}
	return accessTokenStr, refreshTokenStr, nil
}

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := validateToken(tokenStr)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	id, ok := (*claims)["id"].(string)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	c.Locals("id", id)
	return c.Next()
}

func validateToken(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	if exp, ok := claims["expire"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token expired")
		}
	}
	return &claims, nil
}
