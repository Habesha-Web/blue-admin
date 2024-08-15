package controllers

import (
	"net/http"

	"blue-admin.com/common"
	"blue-admin.com/database"
	"blue-admin.com/models"
	"blue-admin.com/observe"
	"blue-admin.com/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

// Login Request for Endpoint
type LoginPost struct {
	GrantType string `json:"grant_type" validate:"required"`
	Email     string `json:"email" validate:"email,min=6,max=32"`
	Password  string `json:"password"`
	Token     string `json:"token"`
}

// Access token Response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// Login is a function to login by EMAIL and ID
// @Summary Auth
// @Description Login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body LoginPost true "Login"
// @Success 200 {object} common.ResponseHTTP{data=TokenResponse{}}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /login [post]
func PostLogin(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	// Getting Databse session
	db, _ := database.ReturnSession()
	validate := validator.New()

	//validating post data
	login_request_data := new(LoginPost)

	//first parse post data
	if err := contx.BodyParser(&login_request_data); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(login_request_data); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	switch login_request_data.GrantType {
	case "authorization_code":
		var user models.User
		res := db.WithContext(tracer.Tracer).Model(&models.User{}).Preload(clause.Associations).Where("email = ?", login_request_data.Email).First(&user)
		if res.Error != nil {
			return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
				Success: false,
				Message: res.Error.Error(),
				Data:    nil,
			})
		} else if utils.PasswordsMatch(user.Password, login_request_data.Password) {
			roles := make([]string, 0, 20)
			for _, value := range user.Roles {

				roles = append(roles, string(value.Name))
			}
			accessString, _ := utils.CreateJWTToken(user.Email, user.UUID, roles, 60)
			refreshString, _ := utils.CreateJWTToken(user.Email, user.UUID, roles, 65)

			data := TokenResponse{
				AccessToken:  accessString,
				RefreshToken: refreshString,
				TokenType:    "Bearer",
			}
			return contx.Status(http.StatusAccepted).JSON(common.ResponseHTTP{
				Success: true,
				Message: "Authorization Granted",
				Data:    data,
			})
		} else {
			return contx.Status(http.StatusUnauthorized).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Make sure You are Providing the Correct Credentials",
				Data:    "Authenthication Failed",
			})
		}
		// return "something"
	case "refresh_token":
		claims, err := utils.ParseJWTToken(login_request_data.Token)
		email := claims.Email
		uuid := claims.UUID
		roles := claims.Roles
		if err == nil {
			accessString, _ := utils.CreateJWTToken(email, uuid, roles, 60)
			refreshString, _ := utils.CreateJWTToken(email, uuid, roles, 65)
			data := TokenResponse{
				AccessToken:  accessString,
				RefreshToken: refreshString,
				TokenType:    "Bearer",
			}
			return contx.Status(http.StatusAccepted).JSON(common.ResponseHTTP{
				Success: true,
				Message: "Authorization Granted",
				Data:    data,
			})
		}

		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Request Type Unknown",
			Data:    "Currently Not Implemented",
		})
	case "token_decode":
		claims, err := utils.ParseJWTToken(login_request_data.Token)

		if err == nil {
			return contx.Status(http.StatusAccepted).JSON(common.ResponseHTTP{
				Success: true,
				Message: "Token decode sucessfull",
				Data:    claims,
			})
		}
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    "Unknown grant type",
		})
	default:
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Request Type Unknown",
			Data:    "Unknown grant type",
		})
	}

}

// CheckLogin is a function to checktoken Status
// @Summary Auth
// @Description CheckLogin
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body LoginPost true "Login"
// @Success 200 {object} common.ResponseHTTP{data=TokenResponse{}}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /checklogin [get]
func CheckLogin(contx *fiber.Ctx) error {

	// Getting header X-APP-TOKEN
	type Token struct {
		token string `reqHeader:"X-APP-TOKEN"`
	}
	var token Token
	err := contx.ReqHeaderParser(token)
	if err != nil {
		contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    "Error Getting Header Value",
		})
	}

	claims, err := utils.ParseJWTToken(token.token)
	//  Decoding the token
	if err != nil {
		return contx.Status(http.StatusAccepted).JSON(common.ResponseHTTP{
			Success: true,
			Message: "Token decode sucessfull",
			Data:    claims,
		})
	}

	// returning the value
	return contx.Status(http.StatusAccepted).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Token decode sucessfull",
		Data:    claims,
	})
}
