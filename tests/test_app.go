package tests

import (
	"blue-admin.com/manager"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	TestApp   *fiber.App
	groupPath = "/api/v1"
)

func setupUserTestApp() {
	godotenv.Load(".test.env")
	TestApp = fiber.New()
	manager.SetupRoutes(TestApp)
}
