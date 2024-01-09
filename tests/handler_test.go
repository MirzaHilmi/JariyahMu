package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Routes().ServeHTTP(rr, req)

	return rr
}

func TestStatus(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/v1/status", nil)
	assert.Nil(t, err)

	res := executeRequest(req)

	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	var statusObj struct {
		Status string
	}

	err = json.Unmarshal(body, &statusObj)
	assert.Nil(t, err)

	assert.Equal(t, "OK", statusObj.Status)
}

// func TestSignUp(t *testing.T) {
// 	newUser := struct {
// 		FullName             string `json:"fullName"`
// 		Email                string `json:"email"`
// 		Password             string `json:"password"`
// 		PasswordConfirmation string `json:"passwordConfirmation"`
// 	}{
// 		FullName:             "Mirza",
// 		Email:                "mirza@gmail.com",
// 		Password:             "12345678",
// 		PasswordConfirmation: "12345678",
// 	}

// 	body, err := json.Marshal(newUser)
// 	assert.Nil(t, err)

// 	req := http.NewRequest(http.MethodPost, "/api/v1/auth/signup", body)
// 	res := executeRequest(req)

// }
