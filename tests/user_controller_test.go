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

type userPost struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Disabled bool   `json:"disabled"`
}

var testsUsersPost = []struct {
	name         string
	description  string
	route        string
	expectedCode int
	postData     userPost
}{
	{
		name:         "post users check new post",
		description:  "HTTP status 200",
		route:        groupPath + "/user",
		expectedCode: 200,
		postData: userPost{
			Email:    "testaddone@mail.com",
			Password: "default@123",
			Disabled: true,
		},
	},
	{
		name:         "post Error post data",
		description:  "HTTP status 500 error data",
		route:        groupPath + "/user",
		expectedCode: 500,
		postData: userPost{
			Email:    "testaddone@mail.com",
			Password: "",
			Disabled: true,
		},
	},
	{
		name:         "post users check two",
		description:  "HTTP status 200 new",
		route:        groupPath + "/user",
		expectedCode: 200,
		postData: userPost{
			Email:    "testaddtwo@mail.com",
			Password: "default@123",
			Disabled: true,
		},
	},
	{
		name:         "post users check three",
		description:  "HTTP status 200 new",
		route:        groupPath + "/user",
		expectedCode: 200,
		postData: userPost{
			Email:    "testaddthree@mail.com",
			Password: "default@123",
			Disabled: true,
		},
	},
	{
		name:         "post users unique check three",
		description:  "HTTP status 400 unique",
		route:        groupPath + "/user",
		expectedCode: 500,
		postData: userPost{
			Email:    "testaddone@mail.com",
			Password: "default@123",
			Disabled: true,
		},
	},
}

var testsUsersPatchID = []struct {
	name         string
	description  string
	route        string
	patchData    models.UserPatch
	expectedCode int
}{
	{
		name:        "patch Users By ID check - 1",
		description: "Patch single user by ID",
		route:       groupPath + "/user/1",
		patchData: models.UserPatch{
			Email:    "testaddtwoupdate@mail.com",
			Disabled: false,
		},
		expectedCode: 200,
	},
	{
		name:        "patch User By ID check - 2",
		description: "HTTP status 404 when user does not exist",
		route:       groupPath + "/user/1000",
		patchData: models.UserPatch{
			Email:    "testaddtwo@mail.com",
			Disabled: false,
		},
		expectedCode: 404,
	},
	{
		name:        "patch User By ID check - 3",
		description: "HTTP status 400 when ID is invalid",
		route:       groupPath + "/user/@@",
		patchData: models.UserPatch{
			Email:    "testaddtwo@mail.com",
			Disabled: false,
		},
		expectedCode: 400,
	},
}

var testsUsersGet = []struct {
	name         string
	description  string
	route        string
	expectedCode int
}{
	{
		name:         "get Users working - 1",
		description:  "HTTP status 200",
		route:        groupPath + "/user?page=1&size=10",
		expectedCode: 200,
	},
	{
		name:         "get Users working - 2",
		description:  "HTTP status 400 when size is negative",
		route:        groupPath + "/user?page=0&size=-5",
		expectedCode: 400,
	},
	{
		name:         "get Users working - 3",
		description:  "HTTP status 400 when page size is zero",
		route:        groupPath + "/user?page=1&size=0",
		expectedCode: 400,
	},
}

var testsAppUsersGet = []struct {
	name         string
	description  string
	route        string
	expectedCode int
}{
	{
		name:         "get App Users working - 1",
		description:  "HTTP status 200",
		route:        groupPath + "/appusers?page=1&size=10&app_uuid=\"something\"",
		expectedCode: 200,
	},
	{
		name:         "get App Users working - 2",
		description:  "HTTP status 400 when size is negative",
		route:        groupPath + "/appusers?page=0&size=-5",
		expectedCode: 400,
	},
	{
		name:         "get App Users working - 3",
		description:  "HTTP status 400 when page size is zero",
		route:        groupPath + "/appusers?page=1&size=0",
		expectedCode: 400,
	},
}

var testsUsersGetByID = []struct {
	name         string
	description  string
	route        string
	expectedCode int
}{
	{
		name:         "get User By ID check - 1",
		description:  "HTTP status 200 for valid user ID",
		route:        groupPath + "/user/1",
		expectedCode: 200,
	},
	{
		name:         "get User By ID check - 2",
		description:  "HTTP status 404 for negative user ID",
		route:        groupPath + "/user/-1",
		expectedCode: 404,
	},
	{
		name:         "get User By ID check - 3",
		description:  "HTTP status 404 for non-existent user ID",
		route:        groupPath + "/user/1000",
		expectedCode: 404,
	},
}

var testsAppUsersGetByID = []struct {
	name         string
	description  string
	route        string
	expectedCode int
}{
	{
		name:         "get App User By ID check - 1",
		description:  "HTTP status 200 for valid user ID",
		route:        groupPath + "/appuser/1?app_uuid=\"someuuid\"",
		expectedCode: 200,
	},
	{
		name:         "get App User By ID check - 2",
		description:  "HTTP status 404 for negative user ID",
		route:        groupPath + "/appuser/-a5?app_uuid=\"someuuid\"",
		expectedCode: 400,
	},
	{
		name:         "get App User By ID check - 3",
		description:  "HTTP status 404 for non-existent user ID",
		route:        groupPath + "/appuser",
		expectedCode: 404,
	},
}

var testsAppUsersGetByUUID = []struct {
	name         string
	description  string
	route        string
	expectedCode int
}{
	{
		name:         "get App User By UUID check - 1",
		description:  "HTTP status 200 for valid user ID",
		route:        groupPath + "/useruuid?uuid=\"someuuid\"&app_uuid=\"someuuid\"",
		expectedCode: 200,
	},
	{
		name:         "get App User By UUID check - 2",
		description:  "HTTP status 404 for negative user ID",
		route:        groupPath + "/useruuid?&app_uuid=\"someuuid\"",
		expectedCode: 400,
	},
	{
		name:         "get App User By UUID check - 3",
		description:  "HTTP status 404 for non-existent user ID",
		route:        groupPath + "/useruuid?uuid=\"someuuid\"",
		expectedCode: 400,
	},
}

func TestUserOperations(t *testing.T) {
	// creating database for test
	models.InitDatabase()
	defer models.CleanDatabase()
	setupUserTestApp()

	//  Testing user Post operations
	for _, test := range testsUsersPost {
		t.Run(test.name, func(t *testing.T) {
			postData, err := json.Marshal(test.postData)
			if err != nil {
				t.Fatalf("Failed to marshal post data: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, test.route, bytes.NewReader(postData))
			req.Header.Set("Content-Type", "application/json")

			resp, err := TestApp.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}

			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

	// Testng user Patch operations
	for _, test := range testsUsersPatchID {
		t.Run(test.name, func(t *testing.T) {
			patchData, err := json.Marshal(test.patchData)
			if err != nil {
				t.Fatalf("Failed to marshal patch data: %v", err)
			}

			req := httptest.NewRequest(http.MethodPatch, test.route, bytes.NewReader(patchData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := TestApp.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}

			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		})
	}

	// Testng user Get operations
	for _, test := range testsUsersGet {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			resp, err := TestApp.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}

			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		})
	}

	// Testng App users Get operations
	for _, test := range testsAppUsersGet {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			resp, err := TestApp.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}

			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		})
	}

	// Testng user Get By ID operations
	for _, test := range testsUsersGetByID {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			resp, err := TestApp.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}

			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		})
	}

	// Testing App user Get By ID operations
	for _, test := range testsAppUsersGetByID {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			resp, err := TestApp.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}

			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		})
	}

	// Testing App user Get By UUID operations
	for _, test := range testsAppUsersGetByUUID {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			resp, err := TestApp.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}

			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		})
	}

	// Testing Delete TestUserOperations
	t.Run("Checking the Delete Request Path for Users", func(t *testing.T) {
		testRoute := fmt.Sprintf("/api/v1/user/%v", 3)
		reqDelete := httptest.NewRequest(http.MethodDelete, testRoute, nil)
		reqDelete.Header.Set("Content-Type", "application/json")

		resp, err := TestApp.Test(reqDelete)
		if err != nil {
			t.Fatalf("Delete request failed: %v", err)
		}

		assert.Equalf(t, http.StatusOK, resp.StatusCode, "user  delete")
	})

	t.Run("Checking the Delete Request Path for Non-existent User", func(t *testing.T) {
		testRoute := fmt.Sprintf("/api/v1/user/%v", 1000000)
		reqDelete := httptest.NewRequest(http.MethodDelete, testRoute, nil)
		reqDelete.Header.Set("Content-Type", "application/json")
		resp, err := TestApp.Test(reqDelete)
		if err != nil {
			t.Fatalf("Delete request failed: %v", err)
		}

		assert.Equalf(t, http.StatusNotFound, resp.StatusCode, " delete non-existent")
	})

	t.Run("Checking the Delete Request Path with Invalid ID", func(t *testing.T) {
		testRoute := fmt.Sprintf("/api/v1/user/%v", "$$$")
		reqDelete := httptest.NewRequest(http.MethodDelete, testRoute, nil)
		reqDelete.Header.Set("Content-Type", "application/json")

		resp, err := TestApp.Test(reqDelete)
		if err != nil {
			t.Fatalf("Delete request failed: %v", err)
		}

		assert.Equalf(t, http.StatusBadRequest, resp.StatusCode, " delete invalid")
	})

}
