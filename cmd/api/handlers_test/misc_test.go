package handlers_test

import (
	"github.com/go-resty/resty/v2"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	"github.com/rhodeon/go-backend-template/test_utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
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
		Get(test_utils.JoinUrlPath(app.Config.Server.BaseUrl, "ping"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, "OK", pingResponseBody.Status)
}
