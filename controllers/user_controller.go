package controllers

import (
	"net/http"
	"strconv"

	"blue-admin.com/common"
	"blue-admin.com/models"
	"blue-admin.com/observe"
	"blue-admin.com/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetUseris a function to get a Users by ID
// @Summary Get Users
// @Description Get Users
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.UserGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /user [get]
func GetUsers(contx *fiber.Ctx) error {

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
	result, err := common.PaginationPureModel(db, models.User{}, []models.User{}, uint(Page), uint(Limit), tracer.Tracer)
	if err != nil {
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Failed to get all User.",
			Data:    "something",
		})
	}

	// returning result if all the above completed successfully
	return contx.Status(http.StatusOK).JSON(result)
}

// Get App Users a function to get app Users by ID
// @Summary Get App Users
// @Description Get App Users
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Param app_uuid query string true "app uuid"
// @Success 200 {object} common.ResponsePagination{data=[]models.UserGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /appusers [get]
func GetAppUsers(contx *fiber.Ctx) error {

	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  parsing Query Prameters
	Page, _ := strconv.Atoi(contx.Query("page"))
	Limit, _ := strconv.Atoi(contx.Query("size"))
	app_uuid := contx.Query("app_uuid")

	//  checking if query parameters  are correct
	if Page == 0 || Limit == 0 || app_uuid == "" {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Not Allowed, Bad request",
			Data:    nil,
		})
	}

	//  querying result with pagination using gorm function
	// result, err := common.PaginationPureModel(db, models.User{}, []models.User{}, uint(Page), uint(Limit), tracer.Tracer)
	query_string := `SELECT DISTINCT u.email, u.uuid, u.id, a.uuid
		FROM users u
		INNER JOIN user_roles ur ON u.id = ur.user_id
		INNER JOIN roles r ON ur.role_id = r.id
		INNER JOIN apps a ON r.app_id = a.id
		WHERE a.uuid = ? limit ? offset ?;
		`
	var users []models.UserGet
	if res := db.WithContext(tracer.Tracer).Raw(query_string, app_uuid, Limit, Page-1).Scan(&users); res.Error != nil {
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Failed to get all User.",
			Data:    "something",
		})
	}

	// returning result if all the above completed successfully
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "sucess get all app Users.",
		Data:    users,
	})
}

// GetUserByID is a function to get a Users by ID
// @Summary Get User by ID
// @Description Get user by ID
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} common.ResponseHTTP{data=models.UserGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /user/{user_id} [get]
func GetUserByID(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  parsing Query Prameters
	id, err := strconv.Atoi(contx.Params("user_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Preparing and querying database using Gorm
	var users_get models.UserGet
	var users models.User
	if res := db.WithContext(tracer.Tracer).Model(&models.User{}).Preload(clause.Associations).Where("id = ?", id).First(&users); res.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// filtering response data according to filtered defined struct
	mapstructure.Decode(users, &users_get)

	//  Finally returing response if All the above compeleted successfully
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one user.",
		Data:    &users_get,
	})
}

// GetAppUserByID is a function to get a Users by ID
// @Summary Get App User by ID
// @Description Get App user by ID
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param app_uuid query string true "app uuid"
// @Success 200 {object} common.ResponseHTTP{data=models.UserGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /appuser/{user_id} [get]
func GetAppUserByID(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  parsing Query Prameters
	id, err := strconv.Atoi(contx.Params("user_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	app_uuid := contx.Query("app_uuid")
	if app_uuid == "" {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "No provided app uuid",
			Data:    nil,
		})
	}

	// Preparing and querying database using Gorm
	var users_get models.UserGet
	query_string := `SELECT DISTINCT u.email, u.uuid, u.id, a.uuid
			FROM users u
			INNER JOIN user_roles ur ON u.id = ur.user_id
			INNER JOIN roles r ON ur.role_id = r.id
			INNER JOIN apps a ON r.app_id = a.id
			WHERE a.uuid = ? AND u.id = ?;`

	if res := db.WithContext(tracer.Tracer).Raw(query_string, app_uuid, id).Scan(&users_get); res.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Some thing happened",
			Data:    nil,
		})
	}

	//  Finally returing response if All the above compeleted successfully
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one user.",
		Data:    &users_get,
	})
}

// GetUserByUUID is a function to get a Users by UUID
// @Summary Get User by UUID
// @Description Get user by UUID
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param uuid query string true "User UUID"
// @Param app_uuid query string true "App UUID"
// @Success 200 {object} common.ResponseHTTP{data=models.UserNoRlnGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /useruuid [get]
func GetUserByUUID(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  parsing Query Prameters
	user_uuid := contx.Query("uuid")
	if user_uuid == "" {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "No User UUID provided",
			Data:    nil,
		})
	}

	//  parsing Query Prameters
	app_uuid := contx.Query("app_uuid")
	if app_uuid == "" {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "No App UUID provided",
			Data:    nil,
		})
	}

	// Preparing and querying database using Gorm
	var users_get models.UserNoRlnGet
	query_string := `SELECT DISTINCT u.email, u.uuid, u.id, a.uuid
			FROM users u
			INNER JOIN user_roles ur ON u.id = ur.user_id
			INNER JOIN roles r ON ur.role_id = r.id
			INNER JOIN apps a ON r.app_id = a.id
			WHERE a.uuid = ? AND u.uuid = ?;`

	if res := db.WithContext(tracer.Tracer).Raw(query_string, app_uuid, user_uuid).Scan(&users_get); res.Error != nil {
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	//  Finally returing response if All the above compeleted successfully
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one user.",
		Data:    &users_get,
	})
}

// Add User to data
// @Summary Add a new User
// @Description Add User
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user body models.UserPost true "Add User"
// @Success 200 {object} common.ResponseHTTP{data=models.UserPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /user [post]
func PostUser(contx *fiber.Ctx) error {
	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Getting Database Connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// validator initialization
	validate := validator.New()

	//validating post data
	posted_user := new(models.UserPost)

	//first parse request data
	if err := contx.BodyParser(&posted_user); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_user); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  initiate -> user
	user := new(models.User)
	user.Email = posted_user.Email
	user.Password = posted_user.Password
	user.Disabled = posted_user.Disabled

	//  start transaction to database
	tx := db.WithContext(tracer.Tracer).Begin()

	// add  data using transaction if values are valid
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "User Creation Failed",
			Data:    err,
		})
	}

	// close transaction
	tx.Commit()

	var user_get models.UserGet
	mapstructure.Decode(user, &user_get)

	// return data if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "User created successfully.",
		Data:    user_get,
	})
}

// Patch User to data
// @Summary Patch User
// @Description Patch User
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user body models.UserPost true "Patch User"
// @Param user_id path int true "User ID"
// @Success 200 {object} common.ResponseHTTP{data=models.UserPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /user/{user_id} [patch]
func PatchUser(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Get database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  initialize data validator
	validate := validator.New()

	// validate path params
	id, err := strconv.Atoi(contx.Params("user_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate data struct
	patch_user := new(models.UserPatch)
	if err := contx.BodyParser(&patch_user); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validating
	if err := validate.Struct(patch_user); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// startng update transaction
	var user models.User
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Where("id = ? ", id).First(&user).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := db.WithContext(tracer.Tracer).Model(&user).UpdateColumns(*patch_user).Update("disabled", patch_user.Disabled).Error; err != nil {
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
		Message: "User updated successfully.",
		Data:    user,
	})
}

// DeleteUsers function removes a user by ID
// @Summary Remove User by ID
// @Description Remove user by ID
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /user/{user_id} [delete]
func DeleteUser(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// get deleted user attributes to return
	var user models.User

	// validate path params
	id, err := strconv.Atoi(contx.Params("user_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// perform delete operation if the object exists
	tx := db.WithContext(tracer.Tracer).Begin()

	// first getting user and checking if it exists
	if err := db.WithContext(tracer.Tracer).Where("id = ?", id).First(&user).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Delete the user
	if id > 7 {
		if err := db.WithContext(tracer.Tracer).Delete(&user).Error; err != nil {
			tx.Rollback()
			return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Error deleting user",
				Data:    nil,
			})
		}
	}

	// Commit the transaction
	tx.Commit()

	// Return success respons
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "User deleted successfully.",
		Data:    user,
	})
}

// DeleteAppUsers function removes a user by ID ( specfic to provided app)
// @Summary Remove App User by ID
// @Description Remove App user by ID
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param app_uuid query string true "app uuid"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /appuser/{user_id} [delete]
func DeleteAppUser(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// get deleted user attributes to return
	var user models.User

	// validate path params
	id, err := strconv.Atoi(contx.Params("user_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  getting app uuid
	app_uuid := contx.Query("app_uuid")
	if app_uuid == "" {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "No provided app uuid",
			Data:    nil,
		})
	}

	// perform delete operation if the object exists
	tx := db.WithContext(tracer.Tracer).Begin()
	query_string := `SELECT DISTINCT u.email, u.uuid, u.id, a.uuid
			FROM users u
			INNER JOIN user_roles ur ON u.id = ur.user_id
			INNER JOIN roles r ON ur.role_id = r.id
			INNER JOIN apps a ON r.app_id = a.id
			WHERE a.uuid = ? AND u.id = ?;
	`
	// first getting user and checking if it exists
	if err := db.WithContext(tracer.Tracer).Raw(query_string, app_uuid, id).Scan(&user).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Delete the user
	if user.ID > 1 {
		if err := db.WithContext(tracer.Tracer).Delete(&user).Error; err != nil {
			tx.Rollback()
			return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Error deleting user",
				Data:    nil,
			})
		}
	}
	// Commit the transaction
	tx.Commit()

	// Return success respons
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "User deleted successfully.",
		Data:    user,
	})
}

// ################################################################
// Relationship Based Endpoints
// ################################################################

// Add Role to User
// @Summary Add User to Role
// @Description Add Role User
// @Tags RoleUsers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Param user_id path int true "User ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /roleuser/{role_id}/{user_id} [post]
func AddRoleUsers(contx *fiber.Ctx) error {

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
	user_id, err := strconv.Atoi(contx.Params("user_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching user to be added
	var user models.User
	// first getting user and checking if it exists
	if err := db.WithContext(tracer.Tracer).Where("id = ?", user_id).First(&user).Error; err != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  userending assocation
	var role models.Role
	role.ID = uint(role_id)
	if err := db.Find(&role); err.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error.Error(),
		})
	}

	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&role).Association("Users").Append(&user); err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Userending User Failed",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Creating a user Role.",
		Data:    user,
	})
}

// Add App Role to User
// @Summary Add User to Role
// @Description Add Role User
// @Tags RoleUsers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Param user_id path int true "User ID"
// @Param app_uuid query string true "app uuid"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /approleuser/{role_id}/{user_id} [post]
func AddAppsRoleUsers(contx *fiber.Ctx) error {

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
	user_id, err := strconv.Atoi(contx.Params("user_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  getting app uuid
	app_uuid := contx.Query("app_uuid")
	if app_uuid == "" {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "No provided app uuid",
			Data:    nil,
		})
	}

	// fetching user to be added
	var user models.User
	if err := db.WithContext(tracer.Tracer).Where("id = ?", user_id).First(&user).Error; err != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})

	}

	//  userending assocation
	var role models.Role
	query_string := `SELECT DISTINCT roles.id, roles.name, roles.description, roles.app_id, roles.active
			FROM roles
			INNER JOIN apps a ON roles.app_id = a.id
			WHERE a.uuid = ? AND roles.id = ?;`
	if res := db.WithContext(tracer.Tracer).Raw(query_string, app_uuid, role_id).Scan(&role); res.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    res.Error.Error(),
		})
	}

	if role.ID != 0 {

		tx := db.WithContext(tracer.Tracer).Begin()
		if err := db.WithContext(tracer.Tracer).Model(&role).Association("Users").Append(&user); err != nil {
			tx.Rollback()
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Appending User Failed",
				Data:    err.Error(),
			})
		}
		tx.Commit()

		// return value if transaction is sucessfull
		return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
			Success: true,
			Message: "Success Creating a user Role.",
			Data:    user,
		})
	}
	// return value if transaction is unsucessful
	return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
		Success: false,
		Message: "Either role or user does not Exist.",
		Data:    nil,
	})
}

// Delete User to Role
// @Summary Add User
// @Description Delete Role User
// @Tags Roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} common.ResponseHTTP{data=models.UserPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /roleuser/{role_id}/{user_id} [delete]
func DeleteRoleUsers(contx *fiber.Ctx) error {

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

	user_id, err := strconv.Atoi(contx.Params("user_id"))
	if err != nil || user_id == 0 {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching user to be deleted
	var user models.User
	if err := db.WithContext(tracer.Tracer).Where("id = ?", user_id).First(&user).Error; err != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})

	}

	// fettchng role
	var role models.Role
	if err := db.WithContext(tracer.Tracer).Where("id = ?", role_id).First(&role).Error; err != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error,
		})
	}

	// removing user
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&role).Association("Users").Delete(&user); err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNonAuthoritativeInfo).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Please Try Again Something Unexpected Huserened",
			Data:    err.Error(),
		})
	}

	tx.Commit()

	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Removing a user from role.",
		Data:    user,
	})
}

// Delete User to Role
// @Summary Add User
// @Description Delete Role User
// @Tags Roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Param user_id path int true "User ID"
// @Param app_uuid query string true "app uuid"
// @Success 200 {object} common.ResponseHTTP{data=models.UserPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /approleuser/{role_id}/{user_id} [delete]
func DeleteAppRoleUsers(contx *fiber.Ctx) error {

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

	user_id, err := strconv.Atoi(contx.Params("user_id"))
	if err != nil || user_id == 0 {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching user to be deleted
	var user models.User
	if err := db.WithContext(tracer.Tracer).Where("id = ?", user_id).First(&user).Error; err != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})

	}

	//  getting app uuid
	app_uuid := contx.Query("app_uuid")
	if app_uuid == "" {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "No provided app uuid",
			Data:    nil,
		})
	}

	// fettchng role
	//  userending assocation
	var role models.Role
	query_string := `SELECT DISTINCT roles.id, roles.name, roles.description, roles.app_id, roles.active
			FROM roles
			INNER JOIN apps a ON roles.app_id = a.id
			WHERE a.uuid = ? AND roles.id = ?;`
	if res := db.WithContext(tracer.Tracer).Raw(query_string, app_uuid, role_id).Scan(&role); res.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    res.Error.Error(),
		})
	}

	if role.ID != 0 {

		// removing user
		tx := db.WithContext(tracer.Tracer).Begin()
		if err := db.WithContext(tracer.Tracer).Model(&role).Association("Users").Delete(&user); err != nil {
			tx.Rollback()
			return contx.Status(http.StatusNonAuthoritativeInfo).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Please Try Again Something Unexpected Huserened",
				Data:    err.Error(),
			})
		}

		tx.Commit()

		// return value if transaction is sucessfull
		return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
			Success: true,
			Message: "Success Removing a user from role.",
			Data:    user,
		})
	}

	// return value if transaction is unsucessful
	return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
		Success: false,
		Message: "Either role or user does not Exist.",
		Data:    nil,
	})

}

// Activate/Deactivate User
// @Summary Activate/Deactivate User
// @Description Activate/Deactivate User
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param status query bool true "Disabled"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /user/{user_id} [put]
func ActivateDeactivateUser(contx *fiber.Ctx) error {
	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
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

	// Getting Query Parameter
	status := contx.QueryBool("status")

	// Fetching User
	var user models.User
	//Updating Didabled Status it it exists
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Where("id = ?", user_id).First(&user).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	db.WithContext(tracer.Tracer).Model(&user).Update("disabled", status)
	tx.Commit()

	var response_user models.UserGet
	if user.ID != 0 {
		mapstructure.Decode(user, &response_user)
		response_user.Disabled = status
		// return value if transaction is sucessfull
		return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
			Success: true,
			Message: "Success Updating a User.",
			Data:    response_user,
		})
	}

	//  Finally return if no record found
	return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
		Success: false,
		Message: "No Record Found",
		Data:    nil,
	})
}

type UserPassword struct {
	Email    string `validate:"required" json:"email" example:"someone@domain.com"`
	Password string `validate:"required" json:"password"`
}

// Update User Password Details
// @Summary Put User
// @Description Put User
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user body UserPassword true "Password User"
// @Param reset query bool true "Reset Password"
// @Success 200 {object} common.ResponseHTTP{data=models.UserGet}
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /user 	[put]
func ChangePassword(contx *fiber.Ctx) error {
	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	validate := validator.New()
	// get query parms
	reset_password := contx.QueryBool("reset")

	// first parsing
	patch_User := new(UserPassword)
	if err := contx.BodyParser(&patch_User); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	// then validating
	if err := validate.Struct(patch_User); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// startng update transaction
	var user_q models.User
	if err := db.WithContext(tracer.Tracer).Model(&user_q).Where("email =?", patch_User.Email).Find(&user_q).Error; err != nil {

		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err,
		})
	}

	var user models.UserGet

	if !reset_password {
		tx := db.WithContext(tracer.Tracer).Begin()
		patch_User.Password = utils.HashFunc(patch_User.Password)
		if err := db.WithContext(tracer.Tracer).Model(&user_q).UpdateColumns(*patch_User).Error; err != nil {
			tx.Rollback()
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Record not Found",
				Data:    err,
			})
		}
		tx.Commit()
	} else {
		tx := db.WithContext(tracer.Tracer).Begin()
		patch_User.Password = utils.HashFunc("default@123")
		if err := db.WithContext(tracer.Tracer).Model(&user_q).UpdateColumns(*patch_User).Error; err != nil {
			tx.Rollback()
			return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Record not Found",
				Data:    err,
			})
		}
		tx.Commit()
	}

	mapstructure.Decode(user_q, &user)
	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Updating a Password.",
		Data:    user,
	})
}
