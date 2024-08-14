
package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"blue-admin.com/common"
	"blue-admin.com/models"
	"blue-admin.com/observe"
)

// GetEndpointis a function to get a Endpoints by ID
// @Summary Get Endpoints
// @Description Get Endpoints
// @Tags Endpoints
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.EndpointGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /endpoint [get]
func GetEndpoints(contx *fiber.Ctx) error {

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
	result, err := common.PaginationPureModel(db, models.Endpoint{}, []models.Endpoint{}, uint(Page), uint(Limit), tracer.Tracer)
	if err != nil {
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Failed to get all Endpoint.",
			Data:    "something",
		})
	}

	// returning result if all the above completed successfully
	return contx.Status(http.StatusOK).JSON(result)
}

// GetEndpointByID is a function to get a Endpoints by ID
// @Summary Get Endpoint by ID
// @Description Get endpoint by ID
// @Tags Endpoints
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param endpoint_id path int true "Endpoint ID"
// @Success 200 {object} common.ResponseHTTP{data=models.EndpointGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /endpoint/{endpoint_id} [get]
func GetEndpointByID(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  parsing Query Prameters
	id, err := strconv.Atoi(contx.Params("endpoint_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}


	// Preparing and querying database using Gorm
	var endpoints_get models.EndpointGet
	var endpoints models.Endpoint
	if res := db.WithContext(tracer.Tracer).Model(&models.Endpoint{}).Preload(clause.Associations).Where("id = ?", id).First(&endpoints); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Endpoint not found",
				Data:    nil,
			})
		}
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving Endpoint",
			Data:    nil,
		})
	}

	// filtering response data according to filtered defined struct
	mapstructure.Decode(endpoints, &endpoints_get)

	//  Finally returing response if All the above compeleted successfully
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one endpoint.",
		Data:    &endpoints_get,
	})
}

// Add Endpoint to data
// @Summary Add a new Endpoint
// @Description Add Endpoint
// @Tags Endpoints
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param endpoint body models.EndpointPost true "Add Endpoint"
// @Success 200 {object} common.ResponseHTTP{data=models.EndpointPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /endpoint [post]
func PostEndpoint(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	// Getting Database Connection
	db, _ := contx.Locals("db").(*gorm.DB)


	// validator initialization
	validate := validator.New()

	//validating post data
	posted_endpoint := new(models.EndpointPost)

	//first parse request data
	if err := contx.BodyParser(&posted_endpoint); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_endpoint); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  initiate -> endpoint
	endpoint := new(models.Endpoint)
	endpoint.Name = posted_endpoint.Name
	endpoint.Description = posted_endpoint.Description

	//  start transaction to database
	tx := db.WithContext(tracer.Tracer).Begin()

	// add  data using transaction if values are valid
	if err := tx.Create(&endpoint).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Endpoint Creation Failed",
			Data:    err,
		})
	}

	// close transaction
	tx.Commit()

	// return data if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Endpoint created successfully.",
		Data:    endpoint,
	})
}

// Patch Endpoint to data
// @Summary Patch Endpoint
// @Description Patch Endpoint
// @Tags Endpoints
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param endpoint body models.EndpointPost true "Patch Endpoint"
// @Param endpoint_id path int true "Endpoint ID"
// @Success 200 {object} common.ResponseHTTP{data=models.EndpointPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /endpoint/{endpoint_id} [patch]
func PatchEndpoint(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	// Get database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  initialize data validator
	validate := validator.New()

	// validate path params
	id, err := strconv.Atoi(contx.Params("endpoint_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate data struct
	patch_endpoint := new(models.EndpointPatch)
	if err := contx.BodyParser(&patch_endpoint); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validating
	if err := validate.Struct(patch_endpoint); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// startng update transaction
	var endpoint models.Endpoint
	endpoint.ID = uint(id)
	tx := db.WithContext(tracer.Tracer).Begin()

	// Check if the record exists
	if err := db.WithContext(tracer.Tracer).First(&endpoint, endpoint.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If the record doesn't exist, return an error response
			tx.Rollback()
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Endpoint not found",
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
	if err := db.WithContext(tracer.Tracer).Model(&endpoint).UpdateColumns(*patch_endpoint).Error; err != nil {
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
		Message: "Endpoint updated successfully.",
		Data:    endpoint,
	})
}

// DeleteEndpoints function removes a endpoint by ID
// @Summary Remove Endpoint by ID
// @Description Remove endpoint by ID
// @Tags Endpoints
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param endpoint_id path int true "Endpoint ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /endpoint/{endpoint_id} [delete]
func DeleteEndpoint(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	// Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// get deleted endpoint attributes to return
	var endpoint models.Endpoint

	// validate path params
	id, err := strconv.Atoi(contx.Params("endpoint_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}


	// perform delete operation if the object exists
	tx := db.WithContext(tracer.Tracer).Begin()

	// first getting endpoint and checking if it exists
	if err := db.Where("id = ?", id).First(&endpoint).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Endpoint not found",
				Data:    nil,
			})
		}
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving endpoint",
			Data:    nil,
		})
	}

	// Delete the endpoint
	if err := db.Delete(&endpoint).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error deleting endpoint",
			Data:    nil,
		})
	}

	// Commit the transaction
	tx.Commit()

	// Return success respons
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Endpoint deleted successfully.",
		Data:    endpoint,
	})
}

// ################################################################
// Relationship Based Endpoints
// ################################################################






