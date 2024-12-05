package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"os/signal"

	"blue-admin.com/bluetasks"
	"blue-admin.com/common"
	"blue-admin.com/configs"
	"blue-admin.com/controllers"
	"blue-admin.com/database"
	_ "blue-admin.com/docs"
	"blue-admin.com/models"
	"blue-admin.com/observe"
	"blue-admin.com/utils"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
)

var (
	env                            string
	BlueAPIRoleManagementSystemcli = &cobra.Command{
		Use:   "run",
		Short: "Run Development server ",
		Long:  `Run Blue API Role Management System development server`,
		Run: func(cmd *cobra.Command, args []string) {

			switch env {
			case "":
				fiber_run("dev")
			default:
				fiber_run(env)
			}

		},
	}

	protectedURLs = []*regexp.Regexp{
		regexp.MustCompile("^/api/v1/login"),
		regexp.MustCompile("^/api/v1/checklogin"),
		regexp.MustCompile("^/api/v1/pics"),
		regexp.MustCompile("^/lmetrics"),
		regexp.MustCompile("^/docs"),
		regexp.MustCompile("^/metrics$"),
	}
)

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

func dbsessioninjection(ctx *fiber.Ctx) error {
	db, err := database.ReturnSession()
	if err != nil {
		return ctx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	ctx.Locals("db", db)
	return ctx.Next()
}

func NextFunc(contx *fiber.Ctx) error {
	return nil
}

// this is path filter which wavies token requirement for provided paths
func authFilter(c *fiber.Ctx) bool {
	originalURL := strings.ToLower(c.OriginalURL())

	for _, pattern := range protectedURLs {
		if pattern.MatchString(originalURL) {
			c.Request().Header.Add("X-APP-TOKEN", "allowed")
			return true
		}
	}
	return false
}

func NextRoute(contx *fiber.Ctx, key string) (bool, error) {
	contx.Next()
	route_name := contx.Route().Name + "_" + strings.ToLower(contx.Route().Method)

	if key == "anonymous" && models.Endpoints_JSON[route_name] == "Anonymous" {
		return true, nil
	} else {

		//  first validating the token
		claims, err := utils.ParseJWTToken(key)
		if err != nil {
			fmt.Println(err)
		}

		// check if the token have the desired role for the route
		role_test := utils.CheckValueExistsInSlice(claims.Roles, models.Endpoints_JSON[route_name])
		if role_test {
			return true, nil
		}
		return false, nil
	}
}

func fiber_run(env string) {

	prefork := false
	if env == "prod" {
		prefork = true
	}
	//  Loading Configuration
	configs.AppConfig.SetEnv(env)

	//  Staring global tracer
	// tp := observe.InitTracer()
	// defer func() {
	// 	if err := tp.Shutdown(context.Background()); err != nil {
	// 		log.Printf("Error shutting down tracer provider: %v", err)
	// 	}
	// }()

	//  lodaing privilge data
	utils.GetAppFeatures()
	//  starting scheduler files
	schd := bluetasks.ScheduledTasks()
	defer schd.Stop()

	//  Creating App logger output file
	//  App should not start with out clearing log scheduler, so panic here error response
	log_file, err := bluetasks.Logfile()
	if err != nil {
		fmt.Printf("Error Creating Logfile %v\n", err)
		panic(err)
	}

	// Basic App Configs
	body_limit, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("BODY_LIMIT", "70"))
	read_buffer_size, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("READ_BUFFER_SIZE", "70"))
	rate_limit_per_second, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("RATE_LIMIT_PER_SECOND", "5000"))

	//load config file
	app := fiber.New(fiber.Config{
		Prefork: prefork,
		// Network:     fiber.NetworkTCP,
		// Immutable:   true,
		JSONEncoder:    json.Marshal,
		JSONDecoder:    json.Unmarshal,
		BodyLimit:      body_limit * 1024 * 1024,
		ReadBufferSize: read_buffer_size * 1024,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError
			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			// Send custom error page
			err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}
			// Return from handler
			return nil
		},
	})

	// recover from panic attacks middlerware
	app.Use(recover.New())

	//  rate limiting middleware
	app.Use(limiter.New(limiter.Config{
		Max:               rate_limit_per_second,
		Expiration:        1 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	// allow cross origin request
	app.Use(cors.New())

	// adding group with authenthication middleware
	SetupRoutes(app)

	// idempotency middleware
	app.Use(idempotency.New(idempotency.Config{
		Lifetime: 10 * time.Second,
	}))

	// logger middle ware with the custom file writer object
	app.Use(logger.New(logger.Config{
		Format:     "\n${cyan}-[${time}]-[${ip}] -${white}${pid} ${red}${status} ${blue}[${method}] ${white}-${path}\n [${body}]\n[${error}]\n[${resBody}]\n[${reqHeaders}]\n[${queryParams}]\n",
		TimeFormat: "15:04:05",
		TimeZone:   "Local",
		Output:     log_file,
	}))

	// prometheus middleware concrete instance
	prometheus := fiberprometheus.New("gobluefiber")
	prometheus.RegisterAt(app, "/metrics")

	// prometheus monitoring middleware
	app.Use(prometheus.Middleware)

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!\n")
	// })
	app.Static("/", "./dist")
	app.Get("/admin/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./dist/index.html")
	})
	app.Get("/admin", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./dist/index.html")
	})

	// swagger docs
	app.Get("/docs/*", swagger.HandlerDefault)
	app.Get("/docs/*", swagger.New()).Name("swagger_routes")

	// fiber native monitoring metrics endpoint
	app.Get("/lmetrics", monitor.New(monitor.Config{Title: "goBlue Metrics Page"})).Name("custom_metrics_route")

	//  Starting Apps and Conumers comes here below
	HTTP_PORT := configs.AppConfig.Get("HTTP_PORT")
	// starting on provided port
	go func(app *fiber.App) {
		app.Listen("0.0.0.0:" + HTTP_PORT)
	}(app)

	// Starting App Conumers
	// // running background consumer on specific quues
	// the provided arument is the name of the queues
	// go func() {
	// 	messages.RabbitConsumer("email", env)
	// 	messages.RabbitConsumer("esb", env)
	// }()

	c := make(chan os.Signal, 1)   // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	app.Shutdown()

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
	fmt.Println("Blue API Role Management System was successful shutdown.")
}

func init() {
	BlueAPIRoleManagementSystemcli.Flags().StringVar(&env, "env", "help", "Which environment to run for example prod or dev")
	goFrame.AddCommand(BlueAPIRoleManagementSystemcli)

}

func SetupRoutes(app *fiber.App) {

	//app logging open telemetery
	app.Use(otelfiber.Middleware())
	// app otel spanner
	app.Use(otelspanstarter)

	// database session injection to local context
	app.Use(dbsessioninjection)

	// Role Middleware
	gapp := app.Group("/api/v1", keyauth.New(keyauth.Config{
		Next:      authFilter,
		KeyLookup: "header:X-APP-TOKEN",
		Validator: NextRoute,
	}))

	gapp.Get("/role", NextFunc).Name("get_all_roles").Get("/role", controllers.GetRoles)
	gapp.Get("/role/:role_id", NextFunc).Name("get_one_roles").Get("/role/:role_id", controllers.GetRoleByID)
	gapp.Post("/role", NextFunc).Name("post_role").Post("/role", controllers.PostRole)
	gapp.Patch("/role/:role_id", NextFunc).Name("patch_role").Patch("/role/:role_id", controllers.PatchRole)
	gapp.Delete("/role/:role_id", NextFunc).Name("delete_role").Delete("/role/:role_id", controllers.DeleteRole)
	gapp.Get("/droproles", NextFunc).Name("drop_roles").Get("/droproles", controllers.GetDropDownRoles)
	gapp.Put("/role/:role_id", NextFunc).Name("activate_deactivate_role").Put("/role/:role_id", controllers.ActivateDeactivateRoles)
	gapp.Get("/role_endpoints", NextFunc).Name("roles_endpoints").Get("/role_endpoints", controllers.GetRoleEndpointsID)

	gapp.Post("/userrole/:user_id/:role_id", NextFunc).Name("add_userrole").Post("/userrole/:user_id/:role_id", controllers.AddUserRoles)
	gapp.Delete("/userrole/:user_id/:role_id", NextFunc).Name("delete_userrole").Delete("/userrole/:user_id/:role_id", controllers.DeleteUserRoles)
	gapp.Patch("/featurerole/:feature_id", NextFunc).Name("add_featurerole").Patch("/featurerole/:feature_id", controllers.AddFeatureRoles)
	gapp.Delete("/featurerole/:feature_id", NextFunc).Name("delete_featurerole").Delete("/featurerole/:feature_id", controllers.DeleteFeatureRoles)

	gapp.Get("/app", NextFunc).Name("get_all_apps").Get("/app", controllers.GetApps)
	gapp.Get("/app/:app_id", NextFunc).Name("get_one_apps").Get("/app/:app_id", controllers.GetAppByID)
	gapp.Get("/appruid/:app_uuid", NextFunc).Name("get_app_roles_uuid").Get("/appruid/:app_uuid", controllers.GetAppRoleUUID)
	gapp.Get("/approleuuid/:app_uuid", NextFunc).Name("get_app_roles_all_uuid").Get("/approleuuid/:app_uuid", controllers.GetAppRoleAllUUID)
	gapp.Post("/app", NextFunc).Name("post_app").Post("/app", controllers.PostApp)
	gapp.Patch("/app/:app_id", NextFunc).Name("patch_app").Patch("/app/:app_id", controllers.PatchApp)
	gapp.Delete("/app/:app_id", NextFunc).Name("delete_app").Delete("/app/:app_id", controllers.DeleteApp).Name("delete_app")

	gapp.Patch("/approle/:role_id", NextFunc).Name("add_roleapp").Patch("/approle/:role_id", controllers.AddRoleApps)
	gapp.Delete("/approle/:role_id", NextFunc).Name("delete_roleapp").Delete("/approle/:role_id", controllers.DeleteRoleApps)

	gapp.Get("/user", NextFunc).Name("get_all_users").Get("/user", controllers.GetUsers)
	gapp.Get("/user/:user_id", NextFunc).Name("get_one_users").Get("/user/:user_id", controllers.GetUserByID)
	gapp.Get("/useruuid", NextFunc).Name("get_one_user_uuid").Get("/useruuid", controllers.GetUserByUUID)
	gapp.Get("/appuser", NextFunc).Name("get_appusers_uuid").Get("/appusers", controllers.GetAppUsers)
	gapp.Get("/appuser/:user_id", NextFunc).Name("get_one_appuser_by_app_id").Get("appuser/:user_id", controllers.GetAppUserByID)
	gapp.Post("/user", NextFunc).Name("post_user").Post("/user", controllers.PostUser)
	gapp.Patch("/user/:user_id", NextFunc).Name("patch_user").Patch("/user/:user_id", controllers.PatchUser)
	gapp.Delete("/user/:user_id", NextFunc).Name("delete_user").Delete("/user/:user_id", controllers.DeleteUser)
	gapp.Delete("/appuser/:user_id", NextFunc).Name("delete_app_user").Delete("/appuser/:user_id", controllers.DeleteAppUser)
	gapp.Put("/user/:user_id", NextFunc).Name("activate_deactivate_user").Put("/user/:user_id", controllers.ActivateDeactivateUser)
	gapp.Put("/user", NextFunc).Name("change_reset_password").Put("/user", controllers.ChangePassword)

	gapp.Post("/roleuser/:role_id/:user_id", NextFunc).Name("add_roleuser").Post("/roleuser/:role_id/:user_id", controllers.AddRoleUsers)
	gapp.Delete("/roleuser/:role_id/:user_id", NextFunc).Name("delete_roleuser").Delete("/roleuser/:role_id/:user_id", controllers.DeleteRoleUsers)
	gapp.Post("/approleuser/:role_id/:user_id", NextFunc).Name("add_approleuser").Post("/approleuser/:role_id/:user_id", controllers.AddAppsRoleUsers)
	gapp.Delete("/approleuser/:role_id/:user_id", NextFunc).Name("delete_approleuser").Delete("/approleuser/:role_id/:user_id", controllers.DeleteAppRoleUsers)

	gapp.Get("/feature", NextFunc).Name("get_all_features").Get("/feature", controllers.GetFeatures)
	gapp.Get("/feature/:feature_id", NextFunc).Name("get_one_features").Get("/feature/:feature_id", controllers.GetFeatureByID)
	gapp.Post("/feature", NextFunc).Name("post_feature").Post("/feature", controllers.PostFeature)
	gapp.Patch("/feature/:feature_id", NextFunc).Name("patch_feature").Patch("/feature/:feature_id", controllers.PatchFeature)
	gapp.Get("/appfeatureuuid/:app_uuid", NextFunc).Name("get_app_feature_all_uuid").Get("/appfeatureuuid/:app_uuid", controllers.GetAppFeaturesAllUUID)
	gapp.Delete("/feature/:feature_id", NextFunc).Name("delete_feature").Delete("/feature/:feature_id", controllers.DeleteFeature)
	gapp.Put("/feature/:feature_id", NextFunc).Name("activate_deactivate_features").Put("/feature/:feature_id", controllers.ActivateDeactivateFeature)
	gapp.Get("/featuredrop", NextFunc).Name("drop_features").Get("/featuredrop", controllers.GetDropFeatures)

	gapp.Patch("/endpointfeature/:endpoint_id", NextFunc).Name("add_endpointfeature").Patch("/endpointfeature/:endpoint_id", controllers.AddEndpointFeatures)
	gapp.Delete("/endpointfeature/:endpoint_id", NextFunc).Name("delete_endpointfeature").Delete("/endpointfeature/:endpoint_id", controllers.DeleteEndpointFeatures)

	gapp.Get("/endpoint", NextFunc).Name("get_all_endpoints").Get("/endpoint", controllers.GetEndpoints)
	gapp.Get("/endpoint/:endpoint_id", NextFunc).Name("get_one_endpoint").Get("/endpoint/:endpoint_id", controllers.GetEndpointByID)
	gapp.Get("/appendpointuuid/:app_uuid", NextFunc).Name("get_app_endpoint_all_uuid").Get("/appendpointuuid/:app_uuid", controllers.GetAppEndpointsAllUUID)
	gapp.Post("/endpoint", NextFunc).Name("post_endpoint").Post("/endpoint", controllers.PostEndpoint)
	gapp.Patch("/endpoint/:endpoint_id", NextFunc).Name("patch_endpoint").Patch("/endpoint/:endpoint_id", controllers.PatchEndpoint)
	gapp.Delete("/endpoint/:endpoint_id", NextFunc).Name("delete_endpoint").Delete("/endpoint/:endpoint_id", controllers.DeleteEndpoint).Name("delete_endpoint")

	gapp.Get("/page", NextFunc).Name("get_all_pages").Get("/page", controllers.GetPages)
	gapp.Get("/page/:page_id", NextFunc).Name("get_one_pages").Get("/page/:page_id", controllers.GetPageByID)
	gapp.Get("/apppagesuuid/:app_uuid", NextFunc).Name("get_app_pages_all_uuid").Get("/apppagesuuid/:app_uuid", controllers.GetAppPagesAllUUID)
	gapp.Post("/page", NextFunc).Name("post_page").Post("/page", controllers.PostPage)
	gapp.Patch("/page/:page_id", NextFunc).Name("patch_page").Patch("/page/:page_id", controllers.PatchPage)
	gapp.Delete("/page/:page_id", NextFunc).Name("delete_page").Delete("/page/:page_id", controllers.DeletePage).Name("delete_page")

	gapp.Post("/rolepage/:role_id/:page_id", NextFunc).Name("add_rolepage").Post("/rolepage/:role_id/:page_id", controllers.AddRolePages)
	gapp.Delete("/rolepage/:role_id/:page_id", NextFunc).Name("delete_rolepage").Delete("/rolepage/:role_id/:page_id", controllers.DeleteRolePages)

	// adding endpoints
	gapp.Get("/checklogin", NextFunc).Name("check_login").Get("/checklogin", controllers.CheckLogin)
	gapp.Post("/login", controllers.PostLogin)

	gapp.Get("/endpointdrop", NextFunc).Name("drop_endpoints").Get("/endpointdrop", controllers.GetDropEndPoints)
	gapp.Get("/appsdrop", NextFunc).Name("drop_sppd").Get("/appsdrop", controllers.GetDropApps)

	// adding email endpoint
	gapp.Get("/email", NextFunc).Name("send_email").Get("/email", controllers.SendEmail).Name("send_email")

	gapp.Get("/jwtsalt", NextFunc).Name("get_all_jwtsalts").Get("/jwtsalt", controllers.GetJWTSalts)

	// Client matrix
	gapp.Get("/clientmatrix/:app_uuid", NextFunc).Name("get_client_matrix").Get("/clientmatrix/:app_uuid", controllers.GetClientMatrix)
	gapp.Get("/clientmatrixpath/:app_uuid", NextFunc).Name("get_client_matrix").Get("/clientmatrixpath/:app_uuid", controllers.GetClientMatrixPath)

	// dashboard
	gapp.Get("/dashboard", NextFunc).Name("dashboard_one").Get("/dashboard", controllers.GetDashBoardGrouped)
	gapp.Get("/dashboardends", NextFunc).Name("dashboard_two").Get("/dashboardends", controllers.GetAppEndpoitnsGroupedBy)
	gapp.Get("/dashboardfeat", NextFunc).Name("dashboard_three").Get("/dashboardfeat", controllers.GetAppFeaturesGroupedBy)
	gapp.Get("/dashboardpages", NextFunc).Name("dashboard_four").Get("/dashboardpages", controllers.GetAppPages)
	gapp.Get("/dashboardroles", NextFunc).Name("dashboard_five").Get("/dashboardroles", controllers.GetAppRoles)
	gapp.Get("/dashboardrolespage", NextFunc).Name("dashboard_six").Get("/dashboardrolespage", controllers.GetAppPagesInRoles)

}
