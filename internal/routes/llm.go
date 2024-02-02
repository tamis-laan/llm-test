package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetRouteLLM(router fiber.Router) {

	router.Post("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello world!!")
	})

}
