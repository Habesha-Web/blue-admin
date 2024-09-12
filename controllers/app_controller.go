package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"blue-admin.com/common"
	"blue-admin.com/models"
	"blue-admin.com/observe"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetAppis a function to get a Apps by ID
// @Summary Get Apps
// @Description Get Apps
// @Tags Apps
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.AppGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /app [get]
func GetApps(contx *fiber.Ctx) error {

	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  parsing Query Prameters
	Page, _ := strconv.Atoi(contx.Query("page"))
	Limit, _ := strconv.Atoi(contx.Query("size"))
	//  checking if query parameters  are correct
	if Page == 0 || Limit == 0 {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Not Allowed, Bad request",
			Data:    nil,
		})
	}

	//  querying result with pagination using gorm function
	result, err := common.PaginationPureModel(db, models.App{}, []models.App{}, uint(Page), uint(Limit), tracer.Tracer)
	if err != nil {
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Failed to get all App.",
			Data:    "something",
		})
	}

	// returning result if all the above completed successfully
	return contx.Status(http.StatusOK).JSON(result)
}

// GetAppByID is a function to get a Apps by ID
// @Summary Get App by ID
// @Description Get app by ID
// @Tags Apps
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param app_id path int true "App ID"
// @Success 200 {object} common.ResponseHTTP{data=models.AppGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /app/{app_id} [get]
func GetAppByID(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  parsing Query Prameters
	id, err := strconv.Atoi(contx.Params("app_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Preparing and querying database using Gorm
	var apps_get models.AppGet
	var apps models.App
	if res := db.WithContext(tracer.Tracer).Model(&models.App{}).Preload(clause.Associations).Where("id = ?", id).First(&apps); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "App not found",
				Data:    nil,
			})
		}
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving App",
			Data:    nil,
		})
	}

	// filtering response data according to filtered defined struct
	mapstructure.Decode(apps, &apps_get)

	//  Finally returing response if All the above compeleted successfully
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one app.",
		Data:    &apps_get,
	})
}

// GetAppRoleUUID is a function to get a Apps by ID
// @Summary Get App Roles by UUID
// @Description Get app roles by UUID
// @Tags Apps
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param app_uuid path string true "App UUID"
// @Success 200 {object} common.ResponseHTTP{data=[]models.RolePut}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /appruid/{app_uuid} [get]
func GetAppRoleUUID(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  parsing Query Prameters
	uuid := contx.Params("app_uuid")
	if uuid == "" {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "No uuid",
			Data:    nil,
		})
	}

	// Preparing and querying database using Gorm
	var roles []models.RolePut
	// select apps.id as appID, roles.id, roles.name, roles.description,roles.active from roles inner join apps on roles.app_id == apps.id where apps.uuid =="0191c74f-d039-71c6-a3be-66e2571a9cf1" ORDER BY roles.id;
	query_string := `select apps.id as appID, roles.id, roles.name, roles.description,roles.active from roles
						inner join apps on roles.app_id == apps.id
						where apps.uuid = ? ORDER BY roles.id;`

	if res := db.WithContext(tracer.Tracer).Raw(query_string, uuid).Scan(&roles); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "App Roles not found",
				Data:    nil,
			})
		}
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving App",
			Data:    nil,
		})
	}

	//  Finally returing response if All the above compeleted successfully
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one app.",
		Data:    &roles,
	})
}

// Add App to data
// @Summary Add a new App
// @Description Add App
// @Tags Apps
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param app body models.AppPost true "Add App"
// @Success 200 {object} common.ResponseHTTP{data=models.AppPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /app [post]
func PostApp(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Getting Database Connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// validator initialization
	validate := validator.New()

	//validating post data
	posted_app := new(models.AppPost)

	//first parse request data
	if err := contx.BodyParser(&posted_app); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_app); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  initiate -> app
	app := new(models.App)
	app.Name = posted_app.Name
	app.Description = posted_app.Description

	//  start transaction to database
	tx := db.WithContext(tracer.Tracer).Begin()

	// add  data using transaction if values are valid
	if err := tx.Create(&app).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "App Creation Failed",
			Data:    err,
		})
	}

	// close transaction
	tx.Commit()

	// return data if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "App created successfully.",
		Data:    app,
	})
}

// Patch App to data
// @Summary Patch App
// @Description Patch App
// @Tags Apps
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param app body models.AppPost true "Patch App"
// @Param app_id path int true "App ID"
// @Success 200 {object} common.ResponseHTTP{data=models.AppPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /app/{app_id} [patch]
func PatchApp(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Get database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  initialize data validator
	validate := validator.New()

	// validate path params
	id, err := strconv.Atoi(contx.Params("app_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate data struct
	patch_app := new(models.AppPatch)
	if err := contx.BodyParser(&patch_app); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validating
	if err := validate.Struct(patch_app); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// startng update transaction
	var app models.App
	app.ID = uint(id)
	tx := db.WithContext(tracer.Tracer).Begin()

	// Check if the record exists
	if err := db.WithContext(tracer.Tracer).First(&app, app.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If the record doesn't exist, return an error response
			tx.Rollback()
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "App not found",
				Data:    nil,
			})
		}
		// If there's an unexpected error, return an internal server error response
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Update the record
	if err := db.WithContext(tracer.Tracer).Model(&app).UpdateColumns(*patch_app).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Return  success response
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "App updated successfully.",
		Data:    app,
	})
}

// DeleteApps function removes a app by ID
// @Summary Remove App by ID
// @Description Remove app by ID
// @Tags Apps
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param app_id path int true "App ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /app/{app_id} [delete]
func DeleteApp(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// get deleted app attributes to return
	var app models.App

	// validate path params
	id, err := strconv.Atoi(contx.Params("app_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// perform delete operation if the object exists
	tx := db.WithContext(tracer.Tracer).Begin()

	// first getting app and checking if it exists
	if err := db.Where("id = ?", id).First(&app).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "App not found",
				Data:    nil,
			})
		}
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving app",
			Data:    nil,
		})
	}

	// Delete the app
	if err := db.Delete(&app).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error deleting app",
			Data:    nil,
		})
	}

	// Commit the transaction
	tx.Commit()

	// Return success respons
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "App deleted successfully.",
		Data:    app,
	})
}

// ################################################################
// Relationship Based Endpoints
// ################################################################

// Add App Role
// @Summary Add App to Role
// @Description Add App to Role
// @Tags Apps
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Param app_id query int true " App ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /roleapp/{role_id} [patch]
func AddRoleApps(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// connect
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	role_id, err := strconv.Atoi(contx.Params("role_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching Endpionts
	var role models.Role
	if res := db.WithContext(tracer.Tracer).Model(&models.Role{}).Where("id = ?", role_id).First(&role); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fetching role to be added
	app_id, _ := strconv.Atoi(contx.Query("app_id"))
	var app models.App
	if res := db.WithContext(tracer.Tracer).Model(&models.App{}).Where("id = ?", app_id).First(&app); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// startng update transaction

	tx := db.WithContext(tracer.Tracer).Begin()
	//  Adding one to many Relation
	if err := db.WithContext(tracer.Tracer).Model(&app).Association("Roles").Append(&role); err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error Adding Record",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Adding a Role to App.",
		Data:    app,
	})
}

// Delete App Role
// @Summary Delete App Role
// @Description Delete App Role
// @Tags Apps
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "App ID"
// @Param app_id query int true "Role ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /roleapp/{role_id} [delete]
func DeleteRoleApps(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	role_id, err := strconv.Atoi(contx.Params("role_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Getting Role
	var role models.Role
	if res := db.WithContext(tracer.Tracer).Model(&models.Role{}).Where("id = ?", role_id).First(&role); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fetching app to be added
	var app models.App
	app_id, _ := strconv.Atoi(contx.Query("app_id"))
	if res := db.WithContext(tracer.Tracer).Model(&models.App{}).Where("id = ?", app_id).First(&app); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// Removing Role From App
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&app).Association("Roles").Delete(&role); err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Deleteing a Role From App.",
		Data:    app,
	})
}

type AppsDropDown struct {
	ID   uint   `validate:"required" json:"id"`
	Name string `validate:"required" json:"name"`
}

// Get Feature Dropdown only active roles
// @Summary Get FeatureDropDown
// @Description Get FeatureDropDown
// @Tags Feature
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseHTTP{data=[]FeatureDropDown}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /appsdrop [get]
func GetDropApps(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	var apps_drop []AppsDropDown
	if res := db.WithContext(tracer.Tracer).Model(&models.App{}).Where("active = ?", true).Find(&apps_drop); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got response",
		Data:    &apps_drop,
	})
}
