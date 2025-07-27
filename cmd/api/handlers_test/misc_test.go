package handlers_test

import (
	"net/http"
	"testing"

	"github.com/rhodeon/go-backend-template/testutils"

	"github.com/go-resty/resty/v2"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	app, err := spawnServer()
	if err != nil {
		t.Fatal(err)
	}

	pingResponseBody := responses.PingResponseBody{}

	resp, err := resty.New().
		R().
		SetResult(&pingResponseBody).
		Get(testutils.JoinUrlPath(app.Config.Server.BaseUrl, "ping"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, "OK", pingResponseBody.Status)
}
