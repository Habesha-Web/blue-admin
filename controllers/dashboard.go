package controllers

import (
	"fmt"
	"net/http"

	"blue-admin.com/common"
	"blue-admin.com/models"
	"blue-admin.com/observe"
	"blue-admin.com/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type endpoints struct {
	Name       string `json:"name"`
	RoutePaths string `json:"route_path"`
	Method     string `json:"method"`
}

type features struct {
	Name      string      `json:"name"`
	Endpoints []endpoints `json:"endpoints"`
}
type roles struct {
	Name     string     `json:"name"`
	Features []features `json:"features"`
}

type AppMeta struct {
	Name  string    `json:"name"`
	UUID  uuid.UUID `json:"uuid"`
	Roles []roles   `json:"roles"`
}
type Role string
type Featuers []string
type EndPoints []string
type AppFeaturesMeta map[Role]Featuers
type AppEndpointsMeta map[Role]EndPoints

// Get App Summary
// @Summary Get Roles Grouped By APP ID
// @Description Get App summary by App ID
// @Tags Dashboard Meta
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param app_id query string true "App ID"
// @Success 200 {object} common.ResponseHTTP{data=AppMeta}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /dashboard [get]
func GetDashBoardGrouped(contx *fiber.Ctx) error {
	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  getting the uuid from query string
	app_uuid, _ := uuid.Parse(contx.Query("app_id"))

	var app models.App
	var app_get AppMeta
	if res := db.WithContext(tracer.Tracer).Model(&models.App{}).Preload(clause.Associations).Preload("Roles.Features").Preload("Roles.Features.Endpoints").Where("uuid = ?", app_uuid).First(&app); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    "nil",
		})
	}
	mapstructure.Decode(app, &app_get)

	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got.",
		Data:    &app_get,
	})
}

// Get App Endpoints
// @Summary Get App Endpoints Grouped By Role
// @Description Get App Endpoints summary
// @Tags Dashboard Meta
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param app_id query string true "App ID"
// @Success 200 {object} common.ResponseHTTP{data=AppEndpointsMeta}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /dashboardends [get]
func GetAppEndpoitnsGroupedBy(contx *fiber.Ctx) error {
	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	var app models.App
	var app_get = make(map[string][]string)
	app_uuid, _ := uuid.Parse(contx.Query("app_id"))

	if res := db.WithContext(tracer.Tracer).Model(&models.App{}).Preload(clause.Associations).Preload("Roles.Features").Preload("Roles.Features.Endpoints").Where("uuid = ?", app_uuid).First(&app); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    "nil",
		})
	}

	for _, value := range app.Roles {
		key := value.Name
		for _, value := range value.Features {
			for _, value := range value.Endpoints {

				app_get[key] = append(app_get[key], value.RoutePath)
			}
		}

	}

	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got.",
		Data:    &app_get,
	})
}

// Get App Features Grouped by Role
// @Summary Get App Features Grouped By Role
// @Description Get App Featuers summary
// @Tags Dashboard Meta
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param app_id query string true "App ID"
// @Success 200 {object} common.ResponseHTTP{data=AppFeaturesMeta}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /dashboardfeat [get]
func GetAppFeaturesGroupedBy(contx *fiber.Ctx) error {
	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	var app models.App
	var app_get = make(map[string][]string)
	app_uuid, _ := uuid.Parse(contx.Query("app_id"))

	if res := db.WithContext(tracer.Tracer).Model(&models.App{}).Preload(clause.Associations).Preload("Roles.Features").Preload("Roles.Features.Endpoints").Where("uuid = ?", app_uuid).First(&app); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    "nil",
		})
	}

	for _, value := range app.Roles {
		key := value.Name
		for _, value := range value.Features {
			app_get[key] = append(app_get[key], value.Name)
		}

	}

	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got.",
		Data:    &app_get,
	})
}

// Get App Pages
// @Summary Get App Pages
// @Description Get App Pages summary
// @Tags Dashboard Meta
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param app_id query string true "App ID"
// @Success 200 {object} common.ResponseHTTP{data=AppFeaturesMeta}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /dashboardpages [get]
func GetAppPages(contx *fiber.Ctx) error {
	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	var app models.App
	var app_pages = make(map[string][]string)
	app_uuid, _ := uuid.Parse(contx.Query("app_id"))

	if res := db.WithContext(tracer.Tracer).Model(&models.App{}).Preload(clause.Associations).Preload("Roles.Pages").Where("uuid = ?", app_uuid).First(&app); res.Error != nil {
		fmt.Println(res.Error.Error())
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    "nil",
		})
	}

	key := app.Name
	for _, value := range app.Roles {
		for _, value := range value.Pages {
			app_pages[key] = append(app_pages[key], value.Name)
		}

	}
	app_page := utils.UniqueSlice(app_pages[key])
	app_pages[key] = app_page
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got.",
		Data:    &app_pages,
	})

}

// Get App Roles
// @Summary Get App Roles
// @Description Get App Roles summary
// @Tags Dashboard Meta
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param app_id query string true "App ID"
// @Success 200 {object} common.ResponseHTTP{data=AppFeaturesMeta}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /dashboardroles [get]
func GetAppRoles(contx *fiber.Ctx) error {
	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)
	var app models.App
	var app_roles = make(map[string][]string)
	app_uuid, _ := uuid.Parse(contx.Query("app_id"))

	if res := db.WithContext(tracer.Tracer).Model(&models.App{}).Preload(clause.Associations).Preload("Roles").Where("uuid = ?", app_uuid).First(&app); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    "nil",
		})
	}

	key := app.Name
	for _, value := range app.Roles {
		app_roles[key] = append(app_roles[key], value.Name)
	}
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got.",
		Data:    &app_roles,
	})

}

// Get Page Found for Every role in App
// @Summary Get Page Roles for Secfic App by ID
// @Description Get Roles by page by ID
// @Tags Dashboard Meta
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param app_id query string true "App ID"
// @Success 200 {object} common.ResponseHTTP{data=models.PageGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /dashboardrolespage [get]
func GetAppPagesInRoles(contx *fiber.Ctx) error {
	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	app_uuid, _ := uuid.Parse(contx.Query("app_id"))
	var role_pages = make(map[string][]string)
	var app models.App
	if res := db.WithContext(tracer.Tracer).Model(&models.App{}).Preload(clause.Associations).Preload("Roles.Pages").Where("uuid = ?", app_uuid).First(&app); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    "nil",
		})
	}
	for _, value := range app.Roles {
		key := value.Name
		for _, value := range value.Pages {
			role_pages[key] = append(role_pages[key], value.Name)
		}
	}
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got.",
		Data:    &role_pages,
	})
}
