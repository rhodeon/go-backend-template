package handlers_test

import (
	_ "net/http"
	_ "testing"

	_ "github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	_ "github.com/rhodeon/go-backend-template/utils/testutils"

	_ "github.com/go-resty/resty/v2"
	_ "github.com/stretchr/testify/assert"
)

//
// func TestPing(t *testing.T) {
// 	app, err := spawnServer()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	responseBody := responses.Envelope[string]{}
//
// 	resp, err := resty.New().
// 		R().
// 		SetResult(&responseBody).
// 		Get(testutils.JoinUrlPath(app.Config.Server.BaseUrl, "ping"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	assert.Equal(t, http.StatusOK, resp.StatusCode())
// 	assert.Equal(t, "success", responseBody.Data)
// }
