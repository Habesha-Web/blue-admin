package controllers

import (
	"net/http"

	"blue-admin.com/common"
	"blue-admin.com/messages"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Send Email to list of users using rabbit
// @Summary Send Email to
// @Description Sending Email
// @Tags Utilities
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param User body messages.EmailMessage true "messages"
// @Success 200 {object} common.ResponseHTTP{data=messages.EmailMessage}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /email [post]
func SendEmail(contx *fiber.Ctx) error {
	validate := validator.New()

	//validating post data
	posted_message := new(messages.EmailMessage)

	//first parse post data
	if err := contx.BodyParser(&posted_message); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_message); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//send to rabbit app module qeue using channel
	// Attempt to publish a message to the queue.
	if err := messages.PublishEmailQueue(*posted_message, "email"); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	// close connection and channel of the rabbitmq server

	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Sent Emails.",
		Data:    posted_message,
	})

}
