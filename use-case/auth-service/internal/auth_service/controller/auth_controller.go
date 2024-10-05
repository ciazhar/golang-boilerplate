package controller

import (
	"context"
	"errors"
	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/auth_service/model"
	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/auth_service/service"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"time"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(
	authService *service.AuthService,
) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// RegisterUser Register User
func (c *AuthController) RegisterUser(ctx *fiber.Ctx) error {
	var user model.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Response{
			Error: "Invalid input",
			Data:  err.Error(),
		})
	}

	newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.authService.RegisterUser(newCtx, user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Could not register user",
			Data:  err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.Response{
		Message: "User registered successfully",
	})
}

// Login User
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var body model.LoginRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Response{
			Error: "Invalid input",
			Data:  err.Error(),
		})
	}

	newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	login, err := c.authService.Login(newCtx, body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Could not login user",
			Data:  err.Error(),
		})
	}

	return ctx.JSON(response.Response{
		Message: "User logged in successfully",
		Data:    login,
	})
}

// RefreshToken Refresh Token Handler
func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {

	token, err := extractToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Error: "No token provided",
			Data:  err.Error(),
		})
	}

	newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newToken, err := c.authService.RefreshToken(newCtx, token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Could not refresh token",
			Data:  err.Error(),
		})
	}

	return ctx.JSON(response.Response{
		Message: "Token refreshed successfully",
		Data:    newToken,
	})
}

// Protected route example
func (c *AuthController) Protected(ctx *fiber.Ctx) error {

	token, err := extractToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Error: "No token provided",
			Data:  err.Error(),
		})
	}

	newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the token exists in Redis
	userId, err := c.authService.Protected(newCtx, token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Could not validate token",
			Data:  err.Error(),
		})
	}

	return ctx.JSON(response.Response{
		Message: "Token validated successfully",
		Data:    model.ProtectedResponse{UserId: userId},
	})
}

// Logout Handler
func (c *AuthController) Logout(ctx *fiber.Ctx) error {

	token, err := extractToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Error: "No token provided",
			Data:  err.Error(),
		})
	}

	newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.authService.Logout(newCtx, token); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Could not logout",
			Data:  err.Error(),
		})
	}

	return ctx.JSON(response.Response{
		Message: "User logged out successfully",
	})
}

// Revoke Handler
func (c *AuthController) Revoke(ctx *fiber.Ctx) error {

	token, err := extractToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Error: "No token provided",
			Data:  err.Error(),
		})
	}

	newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.authService.Revoke(newCtx, token); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Could not revoke tokens",
			Data:  err.Error(),
		})
	}

	return ctx.JSON(response.Response{
		Message: "Tokens revoked successfully",
	})
}

func extractToken(ctx *fiber.Ctx) (string, error) {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no token provided")
	}
	if len(authHeader) < len("Bearer ") {
		return "", errors.New("invalid token format")
	}
	return authHeader[len("Bearer "):], nil
}
