package auth

import (
	"github.com/gofiber/fiber/v3"
	"strings"
)

func isValidToken(c fiber.Ctx) bool {
	token := c.Cookies(AccessToken)
	if token == "" {
		return false
	}

	fields := strings.Fields(token)
	if len(fields) != 2 || fields[0] != AuthType {
		return false
	}

	return true
}

func isValidPublicKey(c fiber.Ctx) bool {
	publicKey := c.Cookies(AccessPublic)
	return publicKey != ""
}

func getAccessToken(c fiber.Ctx) string {
	fields := strings.Fields(c.Cookies(AccessToken))
	return fields[1]
}

func getAccessPublicKey(c fiber.Ctx) string {
	return c.Cookies(AccessPublic)
}
