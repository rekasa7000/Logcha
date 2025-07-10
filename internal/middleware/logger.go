package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func LoggerMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${ip} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		Output: log.Writer(),
	})
}

func RequestLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		
		// Process request
		err := c.Next()
		
		// Log request details
		log.Printf(
			"[%s] %d - %s %s - %s - %v",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Response().StatusCode(),
			c.Method(),
			c.Path(),
			c.IP(),
			time.Since(start),
		)
		
		return err
	}
}