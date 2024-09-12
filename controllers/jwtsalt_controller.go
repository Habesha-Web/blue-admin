package controllers

import (
	"errors"
	"net/http"

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
// @Success 200 {object} common.ResponsePagination{data=models.JWTSalt}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /jwtsalt [get]
func GetJWTSalts(contx *fiber.Ctx) error {

	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	//  querying result with pagination using gorm function
	var salts models.JWTSalt
	if res := db.WithContext(tracer.Tracer).Model(&models.JWTSalt{}).Where("id = ?", 1).First(&salts); res.Error != nil {
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

	// returning result if all the above completed successfully
	//  Finally returing response if All the above compeleted successfully
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one salts.",
		Data:    &salts,
	})
}
