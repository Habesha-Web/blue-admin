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
var testsPagesPostID = []struct {
	name         string          //name of string
	description  string          // description of the test case
	route        string          // route path to test
	page_id      string          //path param
	post_data    models.PagePost // patch_data
	expectedCode int             // expected HTTP status code
}{
	// First test case
	{
		name:        "post Page - 1",
		description: "post Single Page 1",
		route:       groupPath + "/page",
		post_data: models.PagePost{
			Name:        "New one",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	// Second test case
	{
		name:        "post Page - 2",
		description: "post Single 2",
		route:       groupPath + "/page",
		post_data: models.PagePost{
			Name:        "New two",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	// Second Third case
	{
		name:        "post page  - 3",
		description: "Posting page -3",
		route:       groupPath + "/page",
		post_data: models.PagePost{
			Name:        "Name three",
			Description: "Description of Name one",
		},
		expectedCode: 200,
	},
	{
		name:        "post page  - 4",
		description: "Posting page -4",
		route:       groupPath + "/page",
		post_data: models.PagePost{
			Name:        "Name four",
			Description: "Description of Name one",
		},
		expectedCode: 200,
	},
	{
		name:        "post page  - 4",
		description: "Posting page that Exists -5",
		route:       groupPath + "/page",
		post_data: models.PagePost{
			Name:        "Name four",
			Description: "Description of Name one",
		},
		expectedCode: 500,
	},
}

// ##########################################################################
var testsPagesPatchID = []struct {
	name         string           //name of string
	description  string           // description of the test case
	route        string           // route path to test
	patch_data   models.PagePatch // patch_data
	expectedCode int              // expected HTTP status code
}{
	// First test case
	{
		name:        "patch Pages By ID check - 1",
		description: "patch Single Page by ID",
		route:       groupPath + "/page/1",
		patch_data: models.PagePatch{
			Name:        "Name one eight",
			Description: "Description of Name one for test one",
		},
		expectedCode: 200,
	},

	// Second test case
	{
		name:        "get Page By ID check - 2",
		description: "get HTTP status 404, when Page Does not exist",
		route:       groupPath + "/page/1000",
		patch_data: models.PagePatch{
			Name:        "Name one eight",
			Description: "Description of Name one for test 3",
		},
		expectedCode: 404,
	},
	// Second test case
	{
		name:        "get Page By ID check - 4",
		description: "get HTTP status 404, when Page Does not exist",
		route:       groupPath + "/page/@@",
		patch_data: models.PagePatch{
			Name:        "Name one eight",
			Description: "Description of Name one for test 2",
		},
		expectedCode: 400,
	},
}

// ##########################################################################
// Define a structure for specifying input and output data
// of a single test case
var testsPagesGet = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get Pages working - 1",
		description:  "get HTTP status 200",
		route:        groupPath + "/page?page=1&size=10",
		expectedCode: 200,
	},
	// First test case
	{
		name:         "get Pages working - 2",
		description:  "get HTTP status 200",
		route:        groupPath + "/page?page=0&size=-5",
		expectedCode: 400,
	},
	// Second test case
	{
		name:         "get Pages Working - 3",
		description:  "get HTTP status 404, when Page Does not exist",
		route:        groupPath + "/page?page=1&size=0",
		expectedCode: 400,
	},
}

// ##############################################################
var testsPagesGetByID = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get Pages By ID check - 1",
		description:  "get Single Page by ID",
		route:        groupPath + "/page/1",
		expectedCode: 200,
	},

	// First test case
	{
		name:         "get Pages By ID check - 2",
		description:  "get Single Page by ID",
		route:        groupPath + "/page/-1",
		expectedCode: 404,
	},
	// Second test case
	{
		name:         "get Page By ID check - 3",
		description:  "get HTTP status 404, when Page Does not exist",
		route:        groupPath + "/page/1000",
		expectedCode: 404,
	},
}

func TestPagesOperations(t *testing.T) {
	// creating database for test
	models.InitDatabase()
	defer models.CleanDatabase()
	setupUserTestApp()

	// Test Page Post operations
	for _, test := range testsPagesPostID {
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
			// //  running delete test if post is success
			// if resp.StatusCode == 200 {
			// 	t.Run("Checking the Delete Request Path for Pages", func(t *testing.T) {

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

	// test Page Patch Operations cases
	for _, test := range testsPagesPatchID {
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

	// test Pages Get Operatyions cases
	for _, test := range testsPagesGet {
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

	// test Get Pages by ID
	for _, test := range testsPagesGetByID {
		t.Run(test.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			// Add specfic headers if needed as below
			// req.Header.Set("X-APP-TOKEN", "hi")

			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

	// test Delete TestPagesOperations
	// test delete TestEndpointOperations
	t.Run("Checking the Delete Request Path for Page", func(t *testing.T) {
		test_route := fmt.Sprintf("/api/v1/page/%v", 3)
		req_delete := httptest.NewRequest(http.MethodDelete, test_route, nil)

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")
		resp, _ := TestApp.Test(req_delete)

		assert.Equalf(t, 200, resp.StatusCode, "deleteing Page")
	})

	t.Run("Checking the Delete Request Path for  that does not exit", func(t *testing.T) {

		test_route_1 := fmt.Sprintf("/api/v1/page/%v", 1000000)

		req_delete := httptest.NewRequest(http.MethodDelete, test_route_1, nil)

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")

		resp, _ := TestApp.Test(req_delete)
		assert.Equalf(t, 404, resp.StatusCode, "deleteing non existant Page")
	})

	t.Run("Checking the Delete Request Path that is not valid", func(t *testing.T) {

		test_route_2 := fmt.Sprintf("/api/v1/page/%v", "$$$")

		req_delete := httptest.NewRequest(http.MethodDelete, test_route_2, nil)

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")
		resp, _ := TestApp.Test(req_delete)

		assert.Equalf(t, 400, resp.StatusCode, "deleteing error path")
	})
}
