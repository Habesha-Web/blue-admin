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

// GetRoleis a function to get a Roles by ID
// @Summary Get Roles
// @Description Get Roles
// @Tags Roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.RoleGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /role [get]
func GetRoles(contx *fiber.Ctx) error {

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
	result, err := common.PaginationPureModel(db, models.Role{}, []models.Role{}, uint(Page), uint(Limit), tracer.Tracer)
	if err != nil {
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Failed to get all Role.",
			Data:    "something",
		})
	}

	// returning result if all the above completed successfully
	return contx.Status(http.StatusOK).JSON(result)
}

// GetRoleByID is a function to get a Roles by ID
// @Summary Get Role by ID
// @Description Get role by ID
// @Tags Roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Success 200 {object} common.ResponseHTTP{data=models.RoleGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /role/{role_id} [get]
func GetRoleByID(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  parsing Query Prameters
	id, err := strconv.Atoi(contx.Params("role_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Preparing and querying database using Gorm
	var roles_get models.RoleGet
	var roles models.Role
	if res := db.WithContext(tracer.Tracer).Model(&models.Role{}).Preload(clause.Associations).Where("id = ?", id).First(&roles); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Role not found",
				Data:    nil,
			})
		}
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving Role",
			Data:    nil,
		})
	}

	// filtering response data according to filtered defined struct
	mapstructure.Decode(roles, &roles_get)

	//  Finally returing response if All the above compeleted successfully
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one role.",
		Data:    &roles_get,
	})
}

type RoleDropDown struct {
	ID   uint   `validate:"required" json:"id"`
	Name string `validate:"required" json:"name"`
}

// Get Roles Dropdown only active roles
// @Summary Get RoleDropDown
// @Description Get RoleDropDown
// @Tags Role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseHTTP{data=[]RoleDropDown}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /droproles [get]
func GetDropDownRoles(contx *fiber.Ctx) error {

	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	var roles_drop []RoleDropDown
	if res := db.WithContext(tracer.Tracer).Model(&models.Role{}).Where("active = ?", true).Find(&roles_drop); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one role.",
		Data:    &roles_drop,
	})
}

// Add Role to data
// @Summary Add a new Role
// @Description Add Role
// @Tags Roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role body models.RolePost true "Add Role"
// @Success 200 {object} common.ResponseHTTP{data=models.RolePost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /role [post]
func PostRole(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Getting Database Connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// validator initialization
	validate := validator.New()

	//validating post data
	posted_role := new(models.RolePost)

	//first parse request data
	if err := contx.BodyParser(&posted_role); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_role); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  initiate -> role
	role := new(models.Role)
	role.Name = posted_role.Name
	role.Description = posted_role.Description

	//  start transaction to database
	tx := db.WithContext(tracer.Tracer).Begin()

	// add  data using transaction if values are valid
	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Role Creation Failed",
			Data:    err,
		})
	}

	// close transaction
	tx.Commit()

	// return data if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Role created successfully.",
		Data:    role,
	})
}

// Patch Role to data
// @Summary Patch Role
// @Description Patch Role
// @Tags Roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role body models.RolePost true "Patch Role"
// @Param role_id path int true "Role ID"
// @Success 200 {object} common.ResponseHTTP{data=models.RolePost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /role/{role_id} [patch]
func PatchRole(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Get database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  initialize data validator
	validate := validator.New()

	// validate path params
	id, err := strconv.Atoi(contx.Params("role_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate data struct
	patch_role := new(models.RolePatch)
	if err := contx.BodyParser(&patch_role); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validating
	if err := validate.Struct(patch_role); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// startng update transaction
	var role models.Role
	role.ID = uint(id)
	tx := db.WithContext(tracer.Tracer).Begin()

	// Check if the record exists
	if err := db.WithContext(tracer.Tracer).First(&role, role.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If the record doesn't exist, return an error response
			tx.Rollback()
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Role not found",
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
	if err := db.WithContext(tracer.Tracer).Model(&role).UpdateColumns(*patch_role).Error; err != nil {
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
		Message: "Role updated successfully.",
		Data:    role,
	})
}

// DeleteRoles function removes a role by ID
// @Summary Remove Role by ID
// @Description Remove role by ID
// @Tags Roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /role/{role_id} [delete]
func DeleteRole(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// get deleted role attributes to return
	var role models.Role

	// validate path params
	id, err := strconv.Atoi(contx.Params("role_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// perform delete operation if the object exists
	tx := db.WithContext(tracer.Tracer).Begin()

	// first getting role and checking if it exists
	if err := db.WithContext(tracer.Tracer).Where("id = ?", id).First(&role).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Role not found",
				Data:    nil,
			})
		}
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving role",
			Data:    nil,
		})
	}

	// Delete the role
	if err := db.WithContext(tracer.Tracer).Delete(&role).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error deleting role",
			Data:    nil,
		})
	}

	// Commit the transaction
	tx.Commit()

	// Return success respons
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Role deleted successfully.",
		Data:    role,
	})
}

// ################################################################
// Relationship Based Endpoints
// ################################################################

// Add User to Role
// @Summary Add Role to User
// @Description Add User Role
// @Tags UserRoles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param role_id path int true "Role ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /userrole/{user_id}/{role_id} [post]
func AddUserRoles(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	user_id, err := strconv.Atoi(contx.Params("user_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate path params
	role_id, err := strconv.Atoi(contx.Params("role_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching role to be added
	var role models.Role
	if res := db.WithContext(tracer.Tracer).Where("id = ?", role_id).First(&role); res.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	//  roleending assocation
	var user models.User
	if err := db.WithContext(tracer.Tracer).Where("id = ?", user_id).First(&user); err.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error.Error(),
		})
	}

	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&user).Association("Roles").Append(&role); err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Adding Role Failed",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Creating a role User.",
		Data:    role,
	})
}

// Delete Role to User
// @Summary Add Role
// @Description Delete User Role
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param role_id path int true "Role ID"
// @Success 200 {object} common.ResponseHTTP{data=models.RolePost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /userrole/{user_id}/{role_id} [delete]
func DeleteUserRoles(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//Connect to Database
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	user_id, err := strconv.Atoi(contx.Params("user_id"))
	if err != nil || user_id == 0 {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	role_id, err := strconv.Atoi(contx.Params("role_id"))
	if err != nil || role_id == 0 {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	// fetching role to be deleted
	var role models.Role
	if res := db.WithContext(tracer.Tracer).Where("id = ?", role_id).First(&role); res.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fettchng user
	var user models.User
	if err := db.WithContext(tracer.Tracer).Where("id = ?", user_id).First(&user); err.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error.Error(),
			Data:    nil,
		})
	}

	// removing role
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&user).Association("Roles").Delete(&role); err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNonAuthoritativeInfo).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Please Try Again Something Unexpected Hroleened",
			Data:    err.Error(),
		})
	}

	tx.Commit()

	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Removing a role from user.",
		Data:    role,
	})
}

// Add Role Feature
// @Summary Add Role to Feature
// @Description Add Role to Feature
// @Tags Roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param feature_id path int true "Feature ID"
// @Param role_id query int true " Role ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /featurerole/{feature_id} [patch]
func AddFeatureRoles(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// connect
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	feature_id, err := strconv.Atoi(contx.Params("feature_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching Endpionts
	var feature models.Feature
	if res := db.WithContext(tracer.Tracer).Model(&models.Feature{}).Where("id = ?", feature_id).First(&feature); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fetching feature to be added
	role_id, _ := strconv.Atoi(contx.Query("role_id"))
	var role models.Role
	if res := db.WithContext(tracer.Tracer).Model(&models.Role{}).Where("id = ?", role_id).First(&role); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// startng update transaction

	tx := db.WithContext(tracer.Tracer).Begin()
	//  Adding one to many Relation
	if err := db.WithContext(tracer.Tracer).Model(&role).Association("Features").Append(&feature); err != nil {
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
		Message: "Success Adding a Feature to Role.",
		Data:    role,
	})
}

// Delete Role Feature
// @Summary Delete Role Feature
// @Description Delete Role Feature
// @Tags Roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param feature_id path int true "Role ID"
// @Param role_id query int true "Feature ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /featurerole/{feature_id} [delete]
func DeleteFeatureRoles(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	feature_id, err := strconv.Atoi(contx.Params("feature_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Getting Feature
	var feature models.Feature
	if res := db.WithContext(tracer.Tracer).Model(&models.Feature{}).Where("id = ?", feature_id).First(&feature); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fetching role to be added
	var role models.Role
	role_id, _ := strconv.Atoi(contx.Query("role_id"))
	if res := db.WithContext(tracer.Tracer).Model(&models.Role{}).Where("id = ?", role_id).First(&role); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// Removing Feature From Role
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&role).Association("Features").Delete(&feature); err != nil {
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
		Message: "Success Deleteing a Feature From Role.",
		Data:    role,
	})
}

// Activate/Deactivate Role to data
// @Summary Activate/Deactivate
// @Description Activate/Deactivate
// @Tags Role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Param active query bool true "Active"
// @Success 200 {object} common.ResponseHTTP{data=models.RolePost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /roles/{role_id} [put]
func ActivateDeactivateRoles(contx *fiber.Ctx) error {
	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	id, err := strconv.Atoi(contx.Params("role_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	//  Get qurery Parm
	active := contx.QueryBool("active")
	// startng update transaction
	var role models.Role
	role.ID = uint(id)
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Where("id = ? ", id).Model(&role).Update("active", active).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err,
		})
	}
	tx.Commit()

	if role.ID != 0 {
		role.Active = active
		// return value if transaction is sucessfull
		return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
			Success: true,
			Message: "Success Updating a role.",
			Data:    role,
		})
	}

	return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
		Success: false,
		Message: "No Record Found",
		Data:    nil,
	})
}

type EndpiontsRoles struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetRole EndPoints By ID is a function to get a Roles by ID
// @Summary Get EndPoints Role by ID
// @Description Get role EndPoints by ID
// @Tags Role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id query int true "Role ID"
// @Success 200 {object} common.ResponseHTTP{data=[]models.EndpointGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /role_endpoints [get]
func GetRoleEndpointsID(contx *fiber.Ctx) error {
	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	role_id := contx.QueryInt("role_id")
	var endpoints []EndpiontsRoles
	var roles models.Role
	if res := db.WithContext(tracer.Tracer).Model(&models.Role{}).Preload(clause.Associations).Preload("Features.Endpoints").Preload(clause.Associations).Where("id = ?", role_id).First(&roles); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    "nil",
		})
	}

	for x := range roles.Features {
		if len(roles.Features[x].Endpoints) > 0 {
			for i := range roles.Features[x].Endpoints {
				resp_endpoint := EndpiontsRoles{ID: roles.Features[x].Endpoints[i].ID, Name: roles.Features[x].Endpoints[i].Name}
				endpoints = append(endpoints, resp_endpoint)
			}
		}
	}

	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one role.",
		Data:    &endpoints,
	})
}
