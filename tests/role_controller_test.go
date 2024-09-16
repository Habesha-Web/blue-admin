package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"blue-admin.com/models"
	"github.com/stretchr/testify/assert"
)

// go test -coverprofile=coverage.out ./...
// go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

// ##########################################################################
var testsRolesPostID = []struct {
	name         string          //name of string
	description  string          // description of the test case
	route        string          // route path to test
	role_id      string          //path param
	post_data    models.RolePost // patch_data
	expectedCode int             // expected HTTP status code
}{
	// First test case
	{
		name:        "post Role - 1",
		description: "post Single Role 1",
		route:       groupPath + "/role",
		post_data: models.RolePost{
			Name:        "Name one",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	{
		name:        "post Role - 2",
		description: "post Single Role 2",
		route:       groupPath + "/role",
		post_data: models.RolePost{
			Name:        "Name two",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	{
		name:        "post Role - 3",
		description: "post Single Role 3",
		route:       groupPath + "/role",
		post_data: models.RolePost{
			Name:        "Name three",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	// Second test case
	{
		name:        "post Role - 4",
		description: "post Single 4",
		route:       groupPath + "/role",
		post_data: models.RolePost{
			Name:        "New four",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	{
		name:        "invalid post data Role - 4.1",
		description: "invalid post data Single 4",
		route:       groupPath + "/role",
		post_data: models.RolePost{
			Name:        "",
			Description: "",
		},
		expectedCode: 400,
	},
	{
		name:        "unique Role By ID check - 5",
		description: " unique when Role Does exist 5",
		route:       groupPath + "/role",
		post_data: models.RolePost{
			Name:        "Name one",
			Description: "Description of Name one",
		},
		expectedCode: 500,
	},
}

// ##########################################################################
// Define a structure for specifying input and output data
// of a single test case
var testsRolesGet = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get Roles working - 1",
		description:  "get HTTP status 200",
		route:        groupPath + "/role?page=1&size=10",
		expectedCode: 200,
	},
	// First test case
	{
		name:         "get Roles working - 2",
		description:  "get HTTP status 200",
		route:        groupPath + "/role?page=0&size=-5",
		expectedCode: 400,
	},
	// Second test case
	{
		name:         "get Roles Working - 3",
		description:  "get HTTP status 404, when Role Does not exist",
		route:        groupPath + "/role?page=1&size=0",
		expectedCode: 400,
	},
}

// ##########################################################################
var testsRolesPatchID = []struct {
	name         string           //name of string
	description  string           // description of the test case
	route        string           // route path to test
	patch_data   models.RolePatch // patch_data
	expectedCode int              // expected HTTP status code
}{
	// First test case
	{
		name:        "patch Roles By ID check - 1",
		description: "patch Single Role by ID",
		route:       groupPath + "/role/1",
		patch_data: models.RolePatch{
			Name:        "Name one eight",
			Description: "Description of Name one for test one",
		},
		expectedCode: 200,
	},

	// Second test case
	{
		name:        "get Role By ID check - 2",
		description: "get HTTP status 404, when Role Does not exist",
		route:       groupPath + "/role/1000",
		patch_data: models.RolePatch{
			Name:        "Name one eight",
			Description: "Description of Name one for test 3",
		},
		expectedCode: 404,
	},
	// Second test case
	{
		name:        "get Role By ID check - 4",
		description: "get HTTP status 404, when Role Does not exist",
		route:       groupPath + "/role/@@",
		patch_data: models.RolePatch{
			Name:        "Name one eight",
			Description: "Description of Name one for test 2",
		},
		expectedCode: 400,
	},
}

// ##############################################################
var testsRolesGetByID = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get Roles By ID check - 1",
		description:  "get Single Role by ID- 1",
		route:        groupPath + "/role/1",
		expectedCode: 200,
	},

	// First test case
	{
		name:         "get Roles By ID check - 2",
		description:  "get Single Role by ID -2",
		route:        groupPath + "/role/-1",
		expectedCode: 404,
	},
	// Second test case
	{
		name:         "get Role By ID check - 3",
		description:  "get Role By ID check - 3 when Role Does not exist",
		route:        groupPath + "/role/1000",
		expectedCode: 404,
	},
}

func TestRolesOperations(t *testing.T) {
	// creating database for test
	models.InitDatabase()
	defer models.CleanDatabase()
	setupUserTestApp()

	// role post operations
	for _, test := range testsRolesPostID {
		t.Run(test.name, func(t *testing.T) {
			//  changing post data to json
			post_data, _ := json.Marshal(test.post_data)

			req := httptest.NewRequest(http.MethodPost, test.route, bytes.NewReader(post_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
			//  running delete test if post is success
			// if resp.StatusCode == 200 {
			// 	t.Run("Checking the Delete Request Path for Roles", func(t *testing.T) {

			// 		test_route := fmt.Sprintf("%v/%v", test.route, responseMap["data"].(map[string]interface{})["id"])

			// 		req_delete := httptest.NewRequest(http.MethodDelete, test_route, bytes.NewReader(post_data))

			// 		// Add specfic headers if needed as below
			// 		req_delete.Header.Set("Content-Type", "application/json")

			// 		resp, _ := TestApp.Test(req_delete)

			// 		assert.Equalf(t, 200, resp.StatusCode, test.description+"deleteing")
			// 	})
			// } else {
			// 	t.Run("Checking the Delete Request Path for  that does not exit", func(t *testing.T) {

			// 		test_route_1 := fmt.Sprintf("%v/:%v", test.route, 1000000)

			// 		req_delete := httptest.NewRequest(http.MethodDelete, test_route_1, bytes.NewReader(post_data))

			// 		// Add specfic headers if needed as below
			// 		req_delete.Header.Set("Content-Type", "application/json")

			// 		resp, _ := TestApp.Test(req_delete)
			// 		assert.Equalf(t, 500, resp.StatusCode, test.description+"deleteing")
			// 	})

			// 	t.Run("Checking the Delete Request Path that is not valid", func(t *testing.T) {

			// 		test_route_2 := fmt.Sprintf("%v/%v", test.route, "$$$")

			// 		req_delete := httptest.NewRequest(http.MethodDelete, test_route_2, bytes.NewReader(post_data))

			// 		// Add specfic headers if needed as below
			// 		req_delete.Header.Set("Content-Type", "application/json")
			// 		resp, _ := TestApp.Test(req_delete)

			// 		assert.Equalf(t, 500, resp.StatusCode, test.description+"deleteing")
			// 	})
			// }
		})
	}

	//test role patch operations
	for _, test := range testsRolesPatchID {
		t.Run(test.name, func(t *testing.T) {

			//  changing post data to json
			patch_data, _ := json.Marshal(test.patch_data)

			req := httptest.NewRequest(http.MethodPatch, test.route, bytes.NewReader(patch_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder

			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

	//get role operations
	for _, test := range testsRolesGet {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

	// get roles by id
	for _, test := range testsRolesGetByID {
		t.Run(test.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, test.route, nil)
			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

	// test role delete TestRolesOperations
	t.Run("Checking the Delete Request Path for Role", func(t *testing.T) {
		test_route := fmt.Sprintf("/api/v1/role/%v", 3)
		req_delete := httptest.NewRequest(http.MethodDelete, test_route, nil)

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")
		resp, _ := TestApp.Test(req_delete)

		assert.Equalf(t, 200, resp.StatusCode, "deleteing role")
	})

	t.Run("Checking the Delete Request Path for  that does not exit", func(t *testing.T) {

		test_route_1 := fmt.Sprintf("/api/v1/role/%v", 1000000)

		req_delete := httptest.NewRequest(http.MethodDelete, test_route_1, nil)

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")

		resp, _ := TestApp.Test(req_delete)
		assert.Equalf(t, 404, resp.StatusCode, "deleteing non existant Endpoint")
	})

	t.Run("Checking the Delete Request Path that is not valid", func(t *testing.T) {

		test_route_2 := fmt.Sprintf("/api/v1/role/%v", "$$$")

		req_delete := httptest.NewRequest(http.MethodDelete, test_route_2, nil)

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")
		resp, _ := TestApp.Test(req_delete)

		assert.Equalf(t, 400, resp.StatusCode, "deleteing error path")
	})
}
