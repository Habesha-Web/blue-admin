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
var testsAppsPostID = []struct {
	name         string         //name of string
	description  string         // description of the test case
	route        string         // route path to test
	app_id       string         //path param
	post_data    models.AppPost // patch_data
	expectedCode int            // expected HTTP status code
}{
	// First test case
	{
		name:        "post App - 1",
		description: "post Single App",
		route:       "/" + group_path + "/app",
		post_data: models.AppPost{
			Name:        "New one Posted 3",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	// Second test case
	{
		name:        "post App - 2",
		description: "post Single ",
		route:       "/" + group_path + "/app",
		post_data: models.AppPost{
			Name:        "New one Posted 3",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	// Second Third case
	{
		name:        "get App By ID check - 3",
		description: "get HTTP status 404, when App Does not exist",
		route:       "/" + group_path + "/app",
		post_data: models.AppPost{
			Name:        "Name one",
			Description: "Description of Name one",
		},
		expectedCode: 500,
	},
}

func TestPostAppsByID(t *testing.T) {

	ReturnTestApp()

	// Iterate through test single test cases
	for _, test := range testsAppsPostID {
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
				t.Run("Checking the Delete Request Path for Apps", func(t *testing.T) {

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
var testsAppsPatchID = []struct {
	name         string          //name of string
	description  string          // description of the test case
	route        string          // route path to test
	patch_data   models.AppPatch // patch_data
	expectedCode int             // expected HTTP status code
}{
	// First test case
	{
		name:        "patch Apps By ID check - 1",
		description: "patch Single App by ID",
		route:       "/" + group_path + "/app/1",
		patch_data: models.AppPatch{
			Name:        "Name one eight",
			Description: "Description of Name one for test one",
		},
		expectedCode: 200,
	},

	// Second test case
	{
		name:        "get App By ID check - 2",
		description: "get HTTP status 404, when App Does not exist",
		route:       "/" + group_path + "/app/1000",
		patch_data: models.AppPatch{
			Name:        "Name one eight",
			Description: "Description of Name one for test 3",
		},
		expectedCode: 404,
	},
	// Second test case
	{
		name:        "get App By ID check - 4",
		description: "get HTTP status 404, when App Does not exist",
		route:       "/" + group_path + "/app/@@",
		patch_data: models.AppPatch{
			Name:        "Name one eight",
			Description: "Description of Name one for test 2",
		},
		expectedCode: 400,
	},
}

func TestPatchAppsByID(t *testing.T) {

	ReturnTestApp()

	// Iterate through test single test cases
	for _, test := range testsAppsPatchID {
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
var testsAppsGet = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get Apps working - 1",
		description:  "get HTTP status 200",
		route:        "/" + group_path + "/app?page=1&size=10",
		expectedCode: 200,
	},
	// First test case
	{
		name:         "get Apps working - 2",
		description:  "get HTTP status 200",
		route:        "/" + group_path + "/app?page=0&size=-5",
		expectedCode: 400,
	},
	// Second test case
	{
		name:         "get Apps Working - 3",
		description:  "get HTTP status 404, when App Does not exist",
		route:        "/" + group_path + "/app?page=1&size=0",
		expectedCode: 400,
	},
}

func TestGetApps(t *testing.T) {
	ReturnTestApp()

	// Iterate through test single test cases
	for _, test := range testsAppsGet {
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

var testsAppsGetByID = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get Apps By ID check - 1",
		description:  "get Single App by ID",
		route:        "/" + group_path + "/app/1",
		expectedCode: 200,
	},

	// First test case
	{
		name:         "get Apps By ID check - 2",
		description:  "get Single App by ID",
		route:        "/" + group_path + "/app/-1",
		expectedCode: 404,
	},
	// Second test case
	{
		name:         "get App By ID check - 3",
		description:  "get HTTP status 404, when App Does not exist",
		route:        "/" + group_path + "/app/1000",
		expectedCode: 404,
	},
}

func TestGetAppsByID(t *testing.T) {

	ReturnTestApp()

	// Iterate through test single test cases
	for _, test := range testsAppsGetByID {
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