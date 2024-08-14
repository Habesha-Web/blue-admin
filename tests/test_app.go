package tests

import (
	"blue-admin.com/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	TestApp    *fiber.App
	group_path string
)

func nextFunc(contx *fiber.Ctx) error {
	return contx.Next()
}

// initalaizing the app
func ReturnTestApp() {

	// loading env file
	godotenv.Load(".test.env")

	TestApp = fiber.New()

	group_path = "/api/v1"

	app := TestApp.Group(group_path)

	app.Get("/role", nextFunc).Name("get_all_roles").Get("/role", controllers.GetRoles)
	app.Get("/role/:role_id", nextFunc).Name("get_one_roles").Get("/role/:role_id", controllers.GetRoleByID)
	app.Post("/role", nextFunc).Name("post_role").Post("/role", controllers.PostRole)
	app.Patch("/role/:role_id", nextFunc).Name("patch_role").Patch("/role/:role_id", controllers.PatchRole)
	app.Delete("/role/:role_id", nextFunc).Name("delete_role").Delete("/role/:role_id", controllers.DeleteRole).Name("delete_role")

	app.Post("/userrole/:user_id/:role_id", nextFunc).Name("add_userrole").Post("/userrole/:user_id/:role_id", controllers.AddUserRoles)
	app.Delete("/userrole/:user_id/:role_id", nextFunc).Name("delete_userrole").Delete("/userrole/:user_id/:role_id", controllers.DeleteUserRoles)

	app.Post("/featurerole/:feature_id/:role_id", nextFunc).Name("add_featurerole").Post("/featurerole/:feature_id/:role_id", controllers.AddFeatureRoles)
	app.Delete("/featurerole/:feature_id/:role_id", nextFunc).Name("delete_featurerole").Delete("/featurerole/:feature_id/:role_id", controllers.DeleteFeatureRoles)

	app.Get("/app", nextFunc).Name("get_all_apps").Get("/app", controllers.GetApps)
	app.Get("/app/:app_id", nextFunc).Name("get_one_apps").Get("/app/:app_id", controllers.GetAppByID)
	app.Post("/app", nextFunc).Name("post_app").Post("/app", controllers.PostApp)
	app.Patch("/app/:app_id", nextFunc).Name("patch_app").Patch("/app/:app_id", controllers.PatchApp)
	app.Delete("/app/:app_id", nextFunc).Name("delete_app").Delete("/app/:app_id", controllers.DeleteApp).Name("delete_app")

	app.Post("/roleapp/:role_id/:app_id", nextFunc).Name("add_roleapp").Post("/roleapp/:role_id/:app_id", controllers.AddRoleApps)
	app.Delete("/roleapp/:role_id/:app_id", nextFunc).Name("delete_roleapp").Delete("/roleapp/:role_id/:app_id", controllers.DeleteRoleApps)

	app.Get("/user", nextFunc).Name("get_all_users").Get("/user", controllers.GetUsers)
	app.Get("/user/:user_id", nextFunc).Name("get_one_users").Get("/user/:user_id", controllers.GetUserByID)
	app.Post("/user", nextFunc).Name("post_user").Post("/user", controllers.PostUser)
	app.Patch("/user/:user_id", nextFunc).Name("patch_user").Patch("/user/:user_id", controllers.PatchUser)
	app.Delete("/user/:user_id", nextFunc).Name("delete_user").Delete("/user/:user_id", controllers.DeleteUser).Name("delete_user")

	app.Post("/roleuser/:role_id/:user_id", nextFunc).Name("add_roleuser").Post("/roleuser/:role_id/:user_id", controllers.AddRoleUsers)
	app.Delete("/roleuser/:role_id/:user_id", nextFunc).Name("delete_roleuser").Delete("/roleuser/:role_id/:user_id", controllers.DeleteRoleUsers)

	app.Get("/feature", nextFunc).Name("get_all_features").Get("/feature", controllers.GetFeatures)
	app.Get("/feature/:feature_id", nextFunc).Name("get_one_features").Get("/feature/:feature_id", controllers.GetFeatureByID)
	app.Post("/feature", nextFunc).Name("post_feature").Post("/feature", controllers.PostFeature)
	app.Patch("/feature/:feature_id", nextFunc).Name("patch_feature").Patch("/feature/:feature_id", controllers.PatchFeature)
	app.Delete("/feature/:feature_id", nextFunc).Name("delete_feature").Delete("/feature/:feature_id", controllers.DeleteFeature).Name("delete_feature")

	app.Post("/endpointfeature/:endpoint_id/:feature_id", nextFunc).Name("add_endpointfeature").Post("/endpointfeature/:endpoint_id/:feature_id", controllers.AddEndpointFeatures)
	app.Delete("/endpointfeature/:endpoint_id/:feature_id", nextFunc).Name("delete_endpointfeature").Delete("/endpointfeature/:endpoint_id/:feature_id", controllers.DeleteEndpointFeatures)

	app.Get("/endpoint", nextFunc).Name("get_all_endpoints").Get("/endpoint", controllers.GetEndpoints)
	app.Get("/endpoint/:endpoint_id", nextFunc).Name("get_one_endpoints").Get("/endpoint/:endpoint_id", controllers.GetEndpointByID)
	app.Post("/endpoint", nextFunc).Name("post_endpoint").Post("/endpoint", controllers.PostEndpoint)
	app.Patch("/endpoint/:endpoint_id", nextFunc).Name("patch_endpoint").Patch("/endpoint/:endpoint_id", controllers.PatchEndpoint)
	app.Delete("/endpoint/:endpoint_id", nextFunc).Name("delete_endpoint").Delete("/endpoint/:endpoint_id", controllers.DeleteEndpoint).Name("delete_endpoint")

	app.Get("/page", nextFunc).Name("get_all_pages").Get("/page", controllers.GetPages)
	app.Get("/page/:page_id", nextFunc).Name("get_one_pages").Get("/page/:page_id", controllers.GetPageByID)
	app.Post("/page", nextFunc).Name("post_page").Post("/page", controllers.PostPage)
	app.Patch("/page/:page_id", nextFunc).Name("patch_page").Patch("/page/:page_id", controllers.PatchPage)
	app.Delete("/page/:page_id", nextFunc).Name("delete_page").Delete("/page/:page_id", controllers.DeletePage).Name("delete_page")

}
