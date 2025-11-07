package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/handlers/auth"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	dbusers "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/users"
	"github.com/rhodeon/go-backend-template/utils/testutils"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser_success(t *testing.T) {
	t.Parallel()

	app, err := spawnServer()
	if err != nil {
		t.Fatal(err)
	}

	reqBody := auth.RegisterRequestBody{
		Email:     "johndoe@example.com",
		Username:  "johndoe",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	}

	responseBody := responses.Envelope[responses.User]{}

	resp, err := resty.New().
		R().
		SetResult(&responseBody).
		SetBody(reqBody).
		Post(testutils.JoinUrlPath(app.Config.Server.BaseUrl, "/auth/register"))
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, http.StatusOK, resp.StatusCode()) {
		fmt.Println(string(resp.Body()))
	}

	assert.Equal(t, "johndoe", responseBody.Data.Username)
	assert.Equal(t, "johndoe@example.com", responseBody.Data.Email)

	// Confirm that the persisted details match the expected values.
	var persistedUser dbusers.User

	// Get the most recently created user.
	err = app.Db.Pool().QueryRow(context.Background(), `
SELECT username, email
FROM public.users
ORDER BY created_at DESC
LIMIT 1`).
		Scan(&persistedUser.Username, &persistedUser.Email)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "johndoe", persistedUser.Username)
	assert.Equal(t, "johndoe@example.com", persistedUser.Email)
}

func TestRegisterUser_failure(t *testing.T) {
	t.Parallel()

	app, err := spawnServer()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name           string
		reqBody        auth.RegisterRequestBody
		respStatusCode int
		respMessage    string
	}{
		{
			name: "duplicate_email",
			reqBody: auth.RegisterRequestBody{
				Username:  "johndoe",
				Email:     "janedoe@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Password:  "123456",
			},
			respStatusCode: http.StatusConflict,
			respMessage:    "Email already taken",
		},
		{
			name: "duplicate_username",
			reqBody: auth.RegisterRequestBody{
				Username:  "janedoe",
				Email:     "johndoe@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Password:  "123456",
			},
			respStatusCode: http.StatusConflict,
			respMessage:    "Username already taken",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			responseError := apierrors.ApiError{}

			resp, err := resty.New().
				R().
				SetBody(tc.reqBody).
				SetError(&responseError).
				Post(testutils.JoinUrlPath(app.Config.Server.BaseUrl, "/auth/register"))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusConflict, resp.StatusCode())
			assert.Equal(t, tc.respMessage, responseError.Error())
		})
	}
}
