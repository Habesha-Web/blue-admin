
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

// GetFeatureis a function to get a Features by ID
// @Summary Get Features
// @Description Get Features
// @Tags Features
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.FeatureGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /feature [get]
func GetFeatures(contx *fiber.Ctx) error {

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
	result, err := common.PaginationPureModel(db, models.Feature{}, []models.Feature{}, uint(Page), uint(Limit), tracer.Tracer)
	if err != nil {
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Failed to get all Feature.",
			Data:    "something",
		})
	}

	// returning result if all the above completed successfully
	return contx.Status(http.StatusOK).JSON(result)
}

// GetFeatureByID is a function to get a Features by ID
// @Summary Get Feature by ID
// @Description Get feature by ID
// @Tags Features
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param feature_id path int true "Feature ID"
// @Success 200 {object} common.ResponseHTTP{data=models.FeatureGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /feature/{feature_id} [get]
func GetFeatureByID(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  parsing Query Prameters
	id, err := strconv.Atoi(contx.Params("feature_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}


	// Preparing and querying database using Gorm
	var features_get models.FeatureGet
	var features models.Feature
	if res := db.WithContext(tracer.Tracer).Model(&models.Feature{}).Preload(clause.Associations).Where("id = ?", id).First(&features); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Feature not found",
				Data:    nil,
			})
		}
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving Feature",
			Data:    nil,
		})
	}

	// filtering response data according to filtered defined struct
	mapstructure.Decode(features, &features_get)

	//  Finally returing response if All the above compeleted successfully
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one feature.",
		Data:    &features_get,
	})
}

// Add Feature to data
// @Summary Add a new Feature
// @Description Add Feature
// @Tags Features
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param feature body models.FeaturePost true "Add Feature"
// @Success 200 {object} common.ResponseHTTP{data=models.FeaturePost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /feature [post]
func PostFeature(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	// Getting Database Connection
	db, _ := contx.Locals("db").(*gorm.DB)


	// validator initialization
	validate := validator.New()

	//validating post data
	posted_feature := new(models.FeaturePost)

	//first parse request data
	if err := contx.BodyParser(&posted_feature); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_feature); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  initiate -> feature
	feature := new(models.Feature)
	feature.Name = posted_feature.Name
	feature.Description = posted_feature.Description

	//  start transaction to database
	tx := db.WithContext(tracer.Tracer).Begin()

	// add  data using transaction if values are valid
	if err := tx.Create(&feature).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Feature Creation Failed",
			Data:    err,
		})
	}

	// close transaction
	tx.Commit()

	// return data if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Feature created successfully.",
		Data:    feature,
	})
}

// Patch Feature to data
// @Summary Patch Feature
// @Description Patch Feature
// @Tags Features
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param feature body models.FeaturePost true "Patch Feature"
// @Param feature_id path int true "Feature ID"
// @Success 200 {object} common.ResponseHTTP{data=models.FeaturePost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /feature/{feature_id} [patch]
func PatchFeature(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	// Get database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  initialize data validator
	validate := validator.New()

	// validate path params
	id, err := strconv.Atoi(contx.Params("feature_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate data struct
	patch_feature := new(models.FeaturePatch)
	if err := contx.BodyParser(&patch_feature); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validating
	if err := validate.Struct(patch_feature); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// startng update transaction
	var feature models.Feature
	feature.ID = uint(id)
	tx := db.WithContext(tracer.Tracer).Begin()

	// Check if the record exists
	if err := db.WithContext(tracer.Tracer).First(&feature, feature.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If the record doesn't exist, return an error response
			tx.Rollback()
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Feature not found",
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
	if err := db.WithContext(tracer.Tracer).Model(&feature).UpdateColumns(*patch_feature).Error; err != nil {
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
		Message: "Feature updated successfully.",
		Data:    feature,
	})
}

// DeleteFeatures function removes a feature by ID
// @Summary Remove Feature by ID
// @Description Remove feature by ID
// @Tags Features
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param feature_id path int true "Feature ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /feature/{feature_id} [delete]
func DeleteFeature(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	// Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// get deleted feature attributes to return
	var feature models.Feature

	// validate path params
	id, err := strconv.Atoi(contx.Params("feature_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}


	// perform delete operation if the object exists
	tx := db.WithContext(tracer.Tracer).Begin()

	// first getting feature and checking if it exists
	if err := db.Where("id = ?", id).First(&feature).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Feature not found",
				Data:    nil,
			})
		}
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving feature",
			Data:    nil,
		})
	}

	// Delete the feature
	if err := db.Delete(&feature).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error deleting feature",
			Data:    nil,
		})
	}

	// Commit the transaction
	tx.Commit()

	// Return success respons
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Feature deleted successfully.",
		Data:    feature,
	})
}

// ################################################################
// Relationship Based Endpoints
// ################################################################








// Add Feature Endpoint
// @Summary Add Feature to Endpoint
// @Description Add Feature to Endpoint
// @Tags Features
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param endpoint_id path int true "Endpoint ID"
// @Param feature_id query int true " Feature ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /endpointfeature/{endpoint_id} [patch]
func AddEndpointFeatures(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	// connect
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	endpoint_id, err := strconv.Atoi(contx.Params("endpoint_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching Endpionts
	var endpoint models.Endpoint
	if res := db.WithContext(tracer.Tracer).Model(&models.Endpoint{}).Where("id = ?", endpoint_id).First(&endpoint); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fetching endpoint to be added
	feature_id, _ := strconv.Atoi(contx.Query("feature_id"))
	var feature models.Feature
	if res := db.WithContext(tracer.Tracer).Model(&models.Feature{}).Where("id = ?", feature_id).First(&feature); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// startng update transaction

	tx := db.WithContext(tracer.Tracer).Begin()
	//  Adding one to many Relation
	if err := db.WithContext(tracer.Tracer).Model(&feature).Association("Endpoints").Append(&endpoint); err != nil {
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
		Message: "Success Adding a Endpoint to Feature.",
		Data:    feature,
	})
}

// Delete Feature Endpoint
// @Summary Delete Feature Endpoint
// @Description Delete Feature Endpoint
// @Tags Features
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param endpoint_id path int true "Feature ID"
// @Param feature_id query int true "Endpoint ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /endpointfeature/{endpoint_id} [delete]
func DeleteEndpointFeatures(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	//  database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	endpoint_id, err := strconv.Atoi(contx.Params("endpoint_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Getting Endpoint
	var endpoint models.Endpoint
	if res := db.WithContext(tracer.Tracer).Model(&models.Endpoint{}).Where("id = ?", endpoint_id).First(&endpoint); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fetching feature to be added
	var feature models.Feature
	feature_id, _ := strconv.Atoi(contx.Query("feature_id"))
	if res := db.WithContext(tracer.Tracer).Model(&models.Feature{}).Where("id = ?", feature_id).First(&feature); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// Removing Endpoint From Feature
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&feature).Association("Endpoints").Delete(&endpoint); err != nil {
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
		Message: "Success Deleteing a Endpoint From Feature.",
		Data:    feature,
	})
}






