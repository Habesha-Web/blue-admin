package tests

import (
	"fmt"
	"net/http"
	"strings"

	"blue-admin.com/common"
	"blue-admin.com/database"
	"blue-admin.com/manager"
	"blue-admin.com/observe"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel/attribute"
)

var (
	TestApp    *fiber.App
	group_path = "/api/v1"
	groupPath  = "/api/v1"
)

func setupUserTestApp() {
	godotenv.Load(".test.env")
	TestApp = fiber.New()
	TestApp.Use(dbSessionInjection)
	TestApp.Use(otelspanstarter)
	gapp := TestApp.Group("/api/v1")

	manager.SetupRoutes(gapp.(*fiber.Group))
}

func nextFunc(contx *fiber.Ctx) error {
	contx.Next()
	return nil
}

func otelspanstarter(ctx *fiber.Ctx) error {
	//  creating trace context from span if they exist
	route_name := ctx.Path() + "_" + strings.ToLower(ctx.Route().Method)
	tracer, span := observe.FiberAppSpanner(ctx, fmt.Sprintf("%v-root", route_name))
	ctx.Locals("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})
	if err := ctx.Next(); err != nil {
		return err
	}
	span.SetAttributes(attribute.String("response", ctx.Response().String()))
	span.End()
	return nil
}

func dbSessionInjection(c *fiber.Ctx) error {
	db, err := database.ReturnSession()
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	c.Locals("db", db)
	return c.Next()
}

// initalaizing the app
func ReturnTestApp() {

	// loading env file
	godotenv.Load(".test.env")

	TestApp = fiber.New()

	group_path = "/api/v1"

	app := TestApp.Group(group_path)

	manager.SetupRoutes(app.(*fiber.Group))

}
