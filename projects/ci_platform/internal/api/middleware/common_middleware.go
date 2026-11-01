package middleware

import "github.com/gofiber/fiber/v2"

// execute the next method in the stack after provided middleware
// (it allows middleware to handle group routes)
// "Group handlers can also be used as a routing path but they must have
// Next added to them so that the flow can continue."
// https://docs.gofiber.io/guide/grouping#group-handlers
func NewGroup(
	middleware func(*fiber.Ctx) error,
) func(*fiber.Ctx) error {
	return func(context *fiber.Ctx) error {
		if err := middleware(context); err != nil {
			return err
		}
		return context.Next()
	}
}
