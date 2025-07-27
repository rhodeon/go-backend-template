package handlers_test

import (
	"context"
	"net/http"
	"testing"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/models/requests"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	"github.com/rhodeon/go-backend-template/testutils"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_success(t *testing.T) {
	app, err := spawnServer()
	if err != nil {
		t.Fatal(err)
	}

	requestBody := requests.CreateUserRequestBody{
		Username: "johndoe",
		Email:    "johndoe@example.com",
	}

	responseBody := responses.UserResponseBody{}

	resp, err := resty.New().
		R().
		SetResult(&responseBody).
		SetBody(requestBody).
		Post(testutils.JoinUrlPath(app.Config.Server.BaseUrl + "/users"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, "johndoe", responseBody.Username)
	assert.Equal(t, "johndoe@example.com", responseBody.Email)

	// Confirm that the persisted details match the expected values.
	type user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	var persistedUser user

	// Get most recently created user.
	err = app.DbPool.QueryRow(context.Background(), `
SELECT username, email 
FROM users 
ORDER BY created_at DESC 
LIMIT 1`).
		Scan(&persistedUser.Username, &persistedUser.Email)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "johndoe", persistedUser.Username)
	assert.Equal(t, "johndoe@example.com", persistedUser.Email)
}

func TestCreateUser_failure(t *testing.T) {
	app, err := spawnServer()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name           string
		req            requests.CreateUserRequestBody
		respStatusCode int
		respMessage    string
	}{
		{
			name: "duplicate_email",
			req: requests.CreateUserRequestBody{
				Username: "johndoe",
				Email:    "janedoe@example.com",
			},
			respStatusCode: http.StatusConflict,
			respMessage:    `user with email "janedoe@example.com" already exists`,
		},
		{
			name: "duplicate_username",
			req: requests.CreateUserRequestBody{
				Username: "janedoe",
				Email:    "johndoe@example.com",
			},
			respStatusCode: http.StatusConflict,
			respMessage:    `user with username "janedoe" already exists`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			responseError := apierrors.ApiError{}

			resp, err := resty.New().
				R().
				SetBody(tc.req).
				SetError(&responseError).
				Post(testutils.JoinUrlPath(app.Config.Server.BaseUrl + "/users"))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusConflict, resp.StatusCode())
			assert.Equal(t, tc.respMessage, responseError.Error())
		})
	}
}

func TestGetUser_success(t *testing.T) {
	app, err := spawnServer()
	if err != nil {
		t.Fatal(err)
	}

	responseBody := responses.UserResponseBody{}

	resp, err := resty.New().
		R().
		SetResult(&responseBody).
		Get(testutils.JoinUrlPath(app.Config.Server.BaseUrl, "users", "1"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, "janedoe@example.com", responseBody.Email)
	assert.Equal(t, "janedoe", responseBody.Username)
}

func TestGetUser_failure(t *testing.T) {
	app, err := spawnServer()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name           string
		reqUserId      string
		respStatusCode int
		respMessage    string
	}{
		{
			name:           "duplicate_email",
			reqUserId:      "16",
			respStatusCode: http.StatusNotFound,
			respMessage:    `user not found`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			responseError := apierrors.ApiError{}

			resp, err := resty.New().
				R().
				SetError(&responseError).
				Get(testutils.JoinUrlPath(app.Config.Server.BaseUrl, "users", tc.reqUserId))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.respStatusCode, resp.StatusCode())
			assert.Equal(t, tc.respMessage, responseError.Error())
		})
	}
}
