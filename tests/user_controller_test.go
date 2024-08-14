package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"blue-admin.com/models"
	"github.com/stretchr/testify/assert"
)

// go test -coverprofile=coverage.out ./...
// go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

// ##########################################################################
var testsUsersPostID = []struct {
	name         string          //name of string
	description  string          // description of the test case
	route        string          // route path to test
	user_id      string          //path param
	post_data    models.UserPost // patch_data
	expectedCode int             // expected HTTP status code
}{
	// First test case
	{
		name:        "post User - 1",
		description: "post Single User",
		route:       "/" + group_path + "/user",
		post_data:   models.UserPost{
			// Name:        "New one Posted 3",
			// Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	// Second test case
	{
		name:        "post User - 2",
		description: "post Single ",
		route:       "/" + group_path + "/user",
		post_data:   models.UserPost{
			// Name:        "New one Posted 3",
			// Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	// Second Third case
	{
		name:        "get User By ID check - 3",
		description: "get HTTP status 404, when User Does not exist",
		route:       "/" + group_path + "/user",
		post_data:   models.UserPost{
			// Name:        "Name one",
			// Description: "Description of Name one",
		},
		expectedCode: 500,
	},
}

func TestPostUsersByID(t *testing.T) {

	ReturnTestApp()

	// Iterate through test single test cases
	for _, test := range testsUsersPostID {
		t.Run(test.name, func(t *testing.T) {
			//  changing post data to json
			post_data, _ := json.Marshal(test.post_data)

			req := httptest.NewRequest(http.MethodPost, test.route, bytes.NewReader(post_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

			resp, _ := TestApp.Test(req)

			var responseMap map[string]interface{}
			body, _ := io.ReadAll(resp.Body)
			uerr := json.Unmarshal(body, &responseMap)
			if uerr != nil {
				// fmt.Printf("Error marshaling response : %v", uerr)
				fmt.Println()
			}

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
			//  running delete test if post is success
			if resp.StatusCode == 200 {
				t.Run("Checking the Delete Request Path for Users", func(t *testing.T) {

					test_route := fmt.Sprintf("%v/%v", test.route, responseMap["data"].(map[string]interface{})["id"])

					req_delete := httptest.NewRequest(http.MethodDelete, test_route, bytes.NewReader(post_data))

					// Add specfic headers if needed as below
					req_delete.Header.Set("Content-Type", "application/json")

					resp, _ := TestApp.Test(req_delete)

					assert.Equalf(t, 200, resp.StatusCode, test.description+"deleteing")
				})
			} else {
				t.Run("Checking the Delete Request Path for  that does not exit", func(t *testing.T) {

					test_route_1 := fmt.Sprintf("%v/:%v", test.route, 1000000)

					req_delete := httptest.NewRequest(http.MethodDelete, test_route_1, bytes.NewReader(post_data))

					// Add specfic headers if needed as below
					req_delete.Header.Set("Content-Type", "application/json")

					resp, _ := TestApp.Test(req_delete)
					assert.Equalf(t, 500, resp.StatusCode, test.description+"deleteing")
				})

				t.Run("Checking the Delete Request Path that is not valid", func(t *testing.T) {

					test_route_2 := fmt.Sprintf("%v/%v", test.route, "$$$")

					req_delete := httptest.NewRequest(http.MethodDelete, test_route_2, bytes.NewReader(post_data))

					// Add specfic headers if needed as below
					req_delete.Header.Set("Content-Type", "application/json")
					resp, _ := TestApp.Test(req_delete)

					assert.Equalf(t, 500, resp.StatusCode, test.description+"deleteing")
				})
			}
		})
	}

}

// ##########################################################################
var testsUsersPatchID = []struct {
	name         string           //name of string
	description  string           // description of the test case
	route        string           // route path to test
	patch_data   models.UserPatch // patch_data
	expectedCode int              // expected HTTP status code
}{
	// First test case
	{
		name:        "patch Users By ID check - 1",
		description: "patch Single User by ID",
		route:       "/" + group_path + "/user/1",
		patch_data:  models.UserPatch{
			// Name:        "Name one eight",
			// Description: "Description of Name one for test one",
		},
		expectedCode: 200,
	},

	// Second test case
	{
		name:        "get User By ID check - 2",
		description: "get HTTP status 404, when User Does not exist",
		route:       "/" + group_path + "/user/1000",
		patch_data:  models.UserPatch{
			// Name:        "Name one eight",
			// Description: "Description of Name one for test 3",
		},
		expectedCode: 404,
	},
	// Second test case
	{
		name:        "get User By ID check - 4",
		description: "get HTTP status 404, when User Does not exist",
		route:       "/" + group_path + "/user/@@",
		patch_data:  models.UserPatch{
			// Name:        "Name one eight",
			// Description: "Description of Name one for test 2",
		},
		expectedCode: 400,
	},
}

func TestPatchUsersByID(t *testing.T) {

	ReturnTestApp()

	// Iterate through test single test cases
	for _, test := range testsUsersPatchID {
		t.Run(test.name, func(t *testing.T) {

			//  changing post data to json
			patch_data, _ := json.Marshal(test.patch_data)

			req := httptest.NewRequest(http.MethodPatch, test.route, bytes.NewReader(patch_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder

			resp, _ := TestApp.Test(req)

			// for debuging you can uncomment
			// fmt.Println("########")
			// fmt.Println(resp.StatusCode)
			// body, _ := io.ReadAll(resp.Result().Body)
			// fmt.Println(string(body))
			// fmt.Println("########")

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

}

// ##########################################################################
// Define a structure for specifying input and output data
// of a single test case
var testsUsersGet = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get Users working - 1",
		description:  "get HTTP status 200",
		route:        "/" + group_path + "/user?page=1&size=10",
		expectedCode: 200,
	},
	// First test case
	{
		name:         "get Users working - 2",
		description:  "get HTTP status 200",
		route:        "/" + group_path + "/user?page=0&size=-5",
		expectedCode: 400,
	},
	// Second test case
	{
		name:         "get Users Working - 3",
		description:  "get HTTP status 404, when User Does not exist",
		route:        "/" + group_path + "/user?page=1&size=0",
		expectedCode: 400,
	},
}

func TestGetUsers(t *testing.T) {
	ReturnTestApp()

	// Iterate through test single test cases
	for _, test := range testsUsersGet {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			// Add specfic headers if needed as below
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder

			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

}

// ##############################################################

var testsUsersGetByID = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get Users By ID check - 1",
		description:  "get Single User by ID",
		route:        "/" + group_path + "/user/1",
		expectedCode: 200,
	},

	// First test case
	{
		name:         "get Users By ID check - 2",
		description:  "get Single User by ID",
		route:        "/" + group_path + "/user/-1",
		expectedCode: 404,
	},
	// Second test case
	{
		name:         "get User By ID check - 3",
		description:  "get HTTP status 404, when User Does not exist",
		route:        "/" + group_path + "/user/1000",
		expectedCode: 404,
	},
}

func TestGetUsersByID(t *testing.T) {

	ReturnTestApp()

	// Iterate through test single test cases
	for _, test := range testsUsersGetByID {
		t.Run(test.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			// Add specfic headers if needed as below
			// req.Header.Set("X-APP-TOKEN", "hi")

			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

}
