package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"time"
	"atorgold/response"
)

func RateLimiter(max int, expiration time.Duration) func(fiber.Ctx) error {
	return limiter.New(limiter.Config{Max: max, Expiration: expiration, LimitReached: limitReachedFunc, KeyGenerator: func(c fiber.Ctx) string {
		remoteIp := c.IP()
		if c.Get("X-NginX-Proxy") == "true" {
			remoteIp = c.Get("X-Real-IP")
		}

		return remoteIp
	}})
}

func limitReachedFunc(c fiber.Ctx) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(response.ErrorResponse{
		Message: fiber.ErrTooManyRequests.Message,
		Status:  fiber.StatusTooManyRequests,
	})
}

