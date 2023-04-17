package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler struct {
	// svc AuthServicer
}

func New( /* svc AuthServicer */ ) *AuthHandler {
	return &AuthHandler{ /*svc: svc*/ }
}

type LoginUserPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var payload *LoginUserPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return fiber.ErrBadRequest
	}

	claims := jwt.MapClaims{
		"name":  payload.Name,
		"roles": []string{"admin", "user"},
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedJWT, err := token.SignedString([]byte("secret"))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(fiber.Map{"token": signedJWT})
}

func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).
		JSON(fiber.Map{"success": true})
}

func (h *AuthHandler) SignUp(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).
		JSON(fiber.Map{"success": true})
}