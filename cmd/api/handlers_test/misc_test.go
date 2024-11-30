package handlers_test

import (
	"github.com/go-resty/resty/v2"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
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
		Get(app.Config.Server.BaseUrl + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode() != http.StatusOK {
		t.Errorf("want status code %d; got %d", http.StatusOK, resp.StatusCode())
	}

	if pingResponseBody.Status != "OK" {
		t.Errorf("want status code %s; got %s", "OK", pingResponseBody.Status)
	}
}
