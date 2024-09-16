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
var testsFeaturesPostID = []struct {
	name         string             //name of string
	description  string             // description of the test case
	route        string             // route path to test
	feature_id   string             //path param
	post_data    models.FeaturePost // patch_data
	expectedCode int                // expected HTTP status code
}{
	// First test case
	{
		name:        "post Feature - 1",
		description: "post Single Feature",
		route:       groupPath + "/feature",
		post_data: models.FeaturePost{
			Name:        "New one",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	// Second test case
	{
		name:        "post Feature - 2",
		description: "post Single 2 ",
		route:       groupPath + "/feature",
		post_data: models.FeaturePost{
			Name:        "New two",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	// Second Third case
	{
		name:        "Post Feature - 3",
		description: "Post feature 3",
		route:       groupPath + "/feature",
		post_data: models.FeaturePost{
			Name:        "Name three",
			Description: "Description of Name one",
		},
		expectedCode: 200,
	},
	// Second Third case
	{
		name:        "Post Four - 4",
		description: "Post Four 4",
		route:       groupPath + "/feature",
		post_data: models.FeaturePost{
			Name:        "Name four",
			Description: "Description of Name one",
		},
		expectedCode: 200,
	},
	{
		name:        "Post Four - 5",
		description: "Post Four Unique check 5",
		route:       groupPath + "/feature",
		post_data: models.FeaturePost{
			Name:        "Name four",
			Description: "Description of Name one",
		},
		expectedCode: 500,
	},
}

// ##########################################################################
var testsFeaturesPatchID = []struct {
	name         string              //name of string
	description  string              // description of the test case
	route        string              // route path to test
	patch_data   models.FeaturePatch // patch_data
	expectedCode int                 // expected HTTP status code
}{
	// First test case
	{
		name:        "patch Features By ID check - 1",
		description: "patch Single Feature by ID",
		route:       groupPath + "/feature/1",
		patch_data: models.FeaturePatch{
			Name:        "Name one eight",
			Description: "Description of Name one for test one",
		},
		expectedCode: 200,
	},

	// Second test case
	{
		name:        "get Feature By ID check - 2",
		description: "get HTTP status 404, when Feature Does not exist",
		route:       groupPath + "/feature/1000",
		patch_data: models.FeaturePatch{
			Name:        "Name one eight",
			Description: "Description of Name one for test 3",
		},
		expectedCode: 404,
	},
	// Second test case
	{
		name:        "get Feature By ID check - 4",
		description: "get HTTP status 404, when Feature Does not exist",
		route:       groupPath + "/feature/@@",
		patch_data: models.FeaturePatch{
			Name:        "Name one eight",
			Description: "Description of Name one for test 2",
		},
		expectedCode: 400,
	},
}

// ##########################################################################
// Define a structure for specifying input and output data
// of a single test case
var testsFeaturesGet = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get Features working - 1",
		description:  "get HTTP status 200",
		route:        groupPath + "/feature?page=1&size=10",
		expectedCode: 200,
	},
	// First test case
	{
		name:         "get Features working - 2",
		description:  "get HTTP status 200",
		route:        groupPath + "/feature?page=0&size=-5",
		expectedCode: 400,
	},
	// Second test case
	{
		name:         "get Features Working - 3",
		description:  "get HTTP status 404, when Feature Does not exist",
		route:        groupPath + "/feature?page=1&size=0",
		expectedCode: 400,
	},
}

var testsFeaturesGetByID = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get Features By ID check - 1",
		description:  "get Single Feature by ID",
		route:        groupPath + "/feature/1",
		expectedCode: 200,
	},

	// First test case
	{
		name:         "get Features By ID check - 2",
		description:  "get Single Feature by ID",
		route:        groupPath + "/feature/-1",
		expectedCode: 404,
	},
	// Second test case
	{
		name:         "get Feature By ID check - 3",
		description:  "get HTTP status 404, when Feature Does not exist",
		route:        groupPath + "/feature/1000",
		expectedCode: 404,
	},
}

func TestFeaturesOperations(t *testing.T) {
	// creating database for test
	models.InitDatabase()
	defer models.CleanDatabase()
	setupUserTestApp()

	// test feature Post Operations
	for _, test := range testsFeaturesPostID {
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
			// if resp.StatusCode == 200 {
			// 	t.Run("Checking the Delete Request Path for Features", func(t *testing.T) {

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

	//  test feature Patch Operations
	for _, test := range testsFeaturesPatchID {
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

	// test features get operations cases
	for _, test := range testsFeaturesGet {
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

	// test single feature Get cases
	for _, test := range testsFeaturesGetByID {
		t.Run(test.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			// Add specfic headers if needed as below
			// req.Header.Set("X-APP-TOKEN", "hi")

			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

	// test delete features endpoint
	// test delete TestEndpointOperations
	t.Run("Checking the Delete Request Path for Features", func(t *testing.T) {
		test_route := fmt.Sprintf("/api/v1/feature/%v", 3)
		req_delete := httptest.NewRequest(http.MethodDelete, test_route, nil)

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")
		resp, _ := TestApp.Test(req_delete)

		assert.Equalf(t, 200, resp.StatusCode, "deleteing Feature")
	})

	t.Run("Checking the Delete Request Path for  that does not exit", func(t *testing.T) {

		test_route_1 := fmt.Sprintf("/api/v1/feature/%v", 1000000)

		req_delete := httptest.NewRequest(http.MethodDelete, test_route_1, nil)

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")

		resp, _ := TestApp.Test(req_delete)
		assert.Equalf(t, 404, resp.StatusCode, "deleteing non existant Endpoint")
	})

	t.Run("Checking the Delete Request Path that is not valid", func(t *testing.T) {

		test_route_2 := fmt.Sprintf("/api/v1/feature/%v", "$$$")

		req_delete := httptest.NewRequest(http.MethodDelete, test_route_2, nil)

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")
		resp, _ := TestApp.Test(req_delete)

		assert.Equalf(t, 400, resp.StatusCode, "deleteing error path")
	})
}
