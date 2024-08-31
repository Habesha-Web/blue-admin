package controllers

import (
	"net/http"
	"strconv"

	"blue-admin.com/common"
	"blue-admin.com/models"
	"blue-admin.com/observe"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetJWTSaltis a function to get a JWTSalts by ID
// @Summary Get JWTSalts
// @Description Get JWTSalts
// @Tags JWTSalts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.JWTSalt}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /jwtsalt [get]
func GetJWTSalts(contx *fiber.Ctx) error {

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
	result, err := common.PaginationPureModel(db, models.JWTSalt{}, []models.JWTSalt{}, uint(Page), uint(Limit), tracer.Tracer)
	if err != nil {
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Failed to get all JWTSalt.",
			Data:    "something",
		})
	}

	// returning result if all the above completed successfully
	return contx.Status(http.StatusOK).JSON(result)
}
