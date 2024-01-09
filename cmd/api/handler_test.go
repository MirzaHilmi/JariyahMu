package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/debug"
	"testing"

	"github.com/lmittmann/tint"
	"github.com/stretchr/testify/assert"
)

func executeRequest(req *http.Request, app *application) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.routes().ServeHTTP(rr, req)

	return rr
}

func TestStatus(t *testing.T) {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	app, err := bootstrap(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		t.Fatal()
	}

	req, err := http.NewRequest(http.MethodGet, "/api/v1/status", nil)
	assert.Nil(t, err)

	res := executeRequest(req, app)

	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	var statusObj struct {
		Status string
	}

	err = json.Unmarshal(body, &statusObj)
	assert.Nil(t, err)

	assert.Equal(t, "OK", statusObj.Status)
}
