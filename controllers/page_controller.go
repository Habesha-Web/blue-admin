package controllers

import (
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

// GetPageis a function to get a Pages by ID
// @Summary Get Pages
// @Description Get Pages
// @Tags Pages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.PageGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /page [get]
func GetPages(contx *fiber.Ctx) error {

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
	result, err := common.Pagination(db, models.Page{}, []models.Page{}, uint(Page), uint(Limit), tracer.Tracer)
	if err != nil {
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Failed to get all Page.",
			Data:    "something",
		})
	}

	// returning result if all the above completed successfully
	return contx.Status(http.StatusOK).JSON(result)
}

// GetPageByID is a function to get a Pages by ID
// @Summary Get Page by ID
// @Description Get page by ID
// @Tags Pages
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page_id path int true "Page ID"
// @Success 200 {object} common.ResponseHTTP{data=models.PageGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /page/{page_id} [get]
func GetPageByID(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  parsing Query Prameters
	id, err := strconv.Atoi(contx.Params("page_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Preparing and querying database using Gorm
	var pages_get models.PageGet
	var pages models.Page
	if res := db.WithContext(tracer.Tracer).Model(&models.Page{}).Preload(clause.Associations).Where("id = ?", id).First(&pages); res.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// filtering response data according to filtered defined struct
	mapstructure.Decode(pages, &pages_get)

	//  Finally returing response if All the above compeleted successfully
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one page.",
		Data:    &pages_get,
	})
}

// Add Page to data
// @Summary Add a new Page
// @Description Add Page
// @Tags Pages
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page body models.PagePost true "Add Page"
// @Success 200 {object} common.ResponseHTTP{data=models.PagePost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /page [post]
func PostPage(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Getting Database Connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// validator initialization
	validate := validator.New()

	//validating post data
	posted_page := new(models.PagePost)

	//first parse request data
	if err := contx.BodyParser(&posted_page); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_page); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  initiate -> page
	page := new(models.Page)
	page.Name = posted_page.Name
	page.Description = posted_page.Description

	//  start transaction to database
	tx := db.WithContext(tracer.Tracer).Begin()

	// add  data using transaction if values are valid
	if err := tx.Create(&page).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Page Creation Failed",
			Data:    err,
		})
	}

	// close transaction
	tx.Commit()

	// return data if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Page created successfully.",
		Data:    page,
	})
}

// Patch Page to data
// @Summary Patch Page
// @Description Patch Page
// @Tags Pages
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page body models.PagePost true "Patch Page"
// @Param page_id path int true "Page ID"
// @Success 200 {object} common.ResponseHTTP{data=models.PagePatch}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /page/{page_id} [patch]
func PatchPage(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Get database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  initialize data validator
	validate := validator.New()

	// validate path params
	id, err := strconv.Atoi(contx.Params("page_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate data struct
	patch_page := new(models.PagePatch)
	if err := contx.BodyParser(&patch_page); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validating
	if err := validate.Struct(patch_page); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// startng update transaction
	var page models.Page
	tx := db.WithContext(tracer.Tracer).Begin()
	// Check if the record exists
	if err := db.WithContext(tracer.Tracer).Where("id = ?", id).First(&page).Error; err != nil {
		// If the record doesn't exist, return an error response
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Update the record
	if err := db.WithContext(tracer.Tracer).Model(&page).UpdateColumns(*patch_page).Update("active", patch_page.Active).Error; err != nil {
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
		Message: "Page updated successfully.",
		Data:    page,
	})
}

// DeletePages function removes a page by ID
// @Summary Remove Page by ID
// @Description Remove page by ID
// @Tags Pages
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page_id path int true "Page ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /page/{page_id} [delete]
func DeletePage(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// get deleted page attributes to return
	var page models.Page

	// validate path params
	id, err := strconv.Atoi(contx.Params("page_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// perform delete operation if the object exists
	tx := db.WithContext(tracer.Tracer).Begin()

	// first getting page and checking if it exists
	if err := db.Where("id = ?", id).First(&page).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})

	}

	// Delete the page
	if id > 9 {
		if err := db.Delete(&page).Error; err != nil {
			tx.Rollback()
			return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Error deleting page",
				Data:    nil,
			})
		}
	}
	// Commit the transaction
	tx.Commit()

	// Return success respons
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Page deleted successfully.",
		Data:    page,
	})
}

// ################################################################
// Relationship Based Endpoints
// ################################################################

// Add Role to Page
// @Summary Add Page to Role
// @Description Add Role Page
// @Tags RolePages
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Param page_id path int true "Page ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /rolepage/{role_id}/{page_id} [post]
func AddRolePages(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// database connection
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

	// validate path params
	page_id, err := strconv.Atoi(contx.Params("page_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching page to be added
	var page models.Page
	if res := db.WithContext(tracer.Tracer).Where("id = ?", page_id).First(&page); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	//  pageending assocation
	var role models.Role
	if err := db.WithContext(tracer.Tracer).Where("ID = ? ", role_id).First(&role); err.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error.Error(),
		})
	}

	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&role).Association("Pages").Append(&page); err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Pageending Page Failed",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Creating a page Role.",
		Data:    page,
	})
}

// Delete Page to Role
// @Summary Add Page
// @Description Delete Role Page
// @Tags Roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Param page_id path int true "Page ID"
// @Success 200 {object} common.ResponseHTTP{data=models.PagePost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /rolepage/{role_id}/{page_id} [delete]
func DeleteRolePages(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//Connect to Database
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	role_id, err := strconv.Atoi(contx.Params("role_id"))
	if err != nil || role_id == 0 {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	page_id, err := strconv.Atoi(contx.Params("page_id"))
	if err != nil || page_id == 0 {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	// fetching page to be deleted
	var page models.Page
	if res := db.WithContext(tracer.Tracer).Where("id = ?", page_id).First(&page); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fettchng role
	var role models.Role
	role.ID = uint(role_id)
	if err := db.WithContext(tracer.Tracer).Where("id = ?", role_id).First(&role); err.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error.Error(),
		})
	}

	// removing page
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&role).Association("Pages").Delete(&page); err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNonAuthoritativeInfo).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Please Try Again Something Unexpected Hpageened",
			Data:    err.Error(),
		})
	}

	tx.Commit()

	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Removing a page from role.",
		Data:    page,
	})
}
