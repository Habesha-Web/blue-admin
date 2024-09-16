package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"blue-admin.com/controllers"
	"blue-admin.com/models"
	"github.com/stretchr/testify/assert" // add Testify package
)

type PostData struct {
	GrantType string `json:"grant_type"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Token     string `json:"token"`
}

// Define a structure for specifying input and output data
// of a single test case
var testLogin = []struct {
	name         string   //name of string
	description  string   // description of the test case
	route        string   // route path to test
	expectedCode int      // expected HTTP status code
	postData     PostData // expects post data to the uri
}{
	// First test case
	{
		name:         "valid credential login check",
		description:  "Valid Login",
		route:        "/api/v1/login",
		expectedCode: 202,
		postData: PostData{
			GrantType: "authorization_code",
			Email:     "superuser@mail.com",
			Password:  "default@123",
			Token:     "token1",
		},
	},
	// Second test case
	{
		name:         "invalid credential login check",
		description:  "Invalid Login",
		route:        "/api/v1/login",
		expectedCode: 401,
		postData: PostData{
			GrantType: "authorization_code",
			Email:     "superuser@mail.com",
			Password:  "default@12345",
			Token:     "token1",
		},
	},
	{
		name:         "invalid post data login check",
		description:  "get HTTP status 404, when token does not exist",
		route:        "/api/v1/login",
		expectedCode: 400,
		postData: PostData{
			GrantType: "password",
			Email:     "superuser@mail.com",
			Password:  "default@123",
			Token:     "token1",
		},
	},
}

func TestAppsLogin(t *testing.T) {
	// creating database for test
	models.InitDatabase()
	defer models.CleanDatabase()
	setupUserTestApp()

	//  firsit posting user
	t.Run("Posting user first", func(t *testing.T) {
		postUData, err := json.Marshal(userPost{
			Email:    "superuser@mail.com",
			Password: "default@123",
			Disabled: false,
		})
		if err != nil {
			t.Fatalf("Failed to marshal post data: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewReader(postUData))
		req.Header.Set("Content-Type", "application/json")

		resp_pu, err := TestApp.Test(req)
		if err != nil {
			fmt.Println(err)
		}

		assert.Equalf(t, 200, resp_pu.StatusCode, "Posting Before Login")
	})

	// Iterate through test Login cases
	for _, test := range testLogin {

		// testing Login Post routes
		t.Run(test.name, func(t *testing.T) {
			//  changing post data to json
			post_data, _ := json.Marshal(test.postData)

			// Create a new http request with the route from the test case
			req := httptest.NewRequest("POST", test.route, bytes.NewReader(post_data))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := TestApp.Test(req)

			// Verify, if the status code is as expected
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

			if resp.StatusCode == 202 {
				var response struct {
					Data controllers.TokenResponse `json:"data"`
				}
				body, _ := io.ReadAll(resp.Body)
				uerr := json.Unmarshal(body, &response)
				if uerr != nil {
					// fmt.Printf("Error marshaling response : %v", uerr)
				}

				t.Run("Checking Token Decode", func(t *testing.T) {
					post_data_token := PostData{
						GrantType: "token_decode",
						Email:     "superuser@mail.com",
						Password:  "default@123",
						Token:     response.Data.RefreshToken,
					}

					post_data_string, _ := json.Marshal(post_data_token)
					//  checking token decode options
					req := httptest.NewRequest("POST", test.route, bytes.NewReader(post_data_string))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("X-APP-TOKEN", response.Data.AccessToken)
					// checking refresh_token options
					resp, _ := TestApp.Test(req)
					assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
				})

				t.Run("Checking Refresh Token", func(t *testing.T) {
					post_data_token := PostData{
						GrantType: "refresh_token",
						Email:     "superuser@mail.com",
						Password:  "default@123",
						Token:     response.Data.RefreshToken,
					}

					post_data_string, _ := json.Marshal(post_data_token)
					//  checking token decode options
					req := httptest.NewRequest("POST", test.route, bytes.NewReader(post_data_string))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("X-APP-TOKEN", response.Data.AccessToken)
					// checking refresh_token options
					resp, _ := TestApp.Test(req)
					assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
				})

			}

		})

	}

	// checking loggin now
	// testing Login Post routes
	t.Run("Check Login Route checks", func(t *testing.T) {
		//  changing post data to json
		post_data, _ := json.Marshal(PostData{
			GrantType: "authorization_code",
			Email:     "superuser@mail.com",
			Password:  "default@123",
			Token:     "token1",
		})

		// Create a new http request with the route from the test case
		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(post_data))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := TestApp.Test(req)

		// Verify, if the status code is as expected
		assert.Equalf(t, 202, resp.StatusCode, "Loggin in First ")

		if resp.StatusCode == 202 {
			var response struct {
				Data controllers.TokenResponse `json:"data"`
			}
			body, _ := io.ReadAll(resp.Body)
			uerr := json.Unmarshal(body, &response)
			if uerr != nil {
				// fmt.Printf("Error marshaling response : %v", uerr)
			}
			check_token := response.Data.AccessToken
			t.Run("Checking Login Status -Working", func(t *testing.T) {
				//  checking token decode options
				req := httptest.NewRequest("GET", "/api/v1/checklogin", nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-APP-TOKEN", check_token)

				// checking refresh_token options
				resp_lg, _ := TestApp.Test(req)

				assert.Equalf(t, 202, resp_lg.StatusCode, "Check Login status working")
			})

			t.Run("Checking Login Status -No Header header", func(t *testing.T) {

				//  checking token decode options
				req := httptest.NewRequest("GET", "/api/v1/checklogin", nil)
				req.Header.Set("Content-Type", "application/json")

				// checking refresh_token options
				resp_lg, _ := TestApp.Test(req)
				assert.Equalf(t, 403, resp_lg.StatusCode, "Check Login status  No header")
			})

			t.Run("Checking Login Status -Error header", func(t *testing.T) {

				//  checking token decode options
				req := httptest.NewRequest("GET", "/api/v1/checklogin", nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-APP-TOKEN", "hi")

				// checking refresh_token options
				resp_lg, _ := TestApp.Test(req)
				assert.Equalf(t, 403, resp_lg.StatusCode, "Check Login status Error header")
			})
		}
	})
}
