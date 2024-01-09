package main

import (
	"net/http"
	"time"

	"github.com/MirzaHilmi/JariyahMu/internal/password"
	"github.com/MirzaHilmi/JariyahMu/internal/request"
	"github.com/MirzaHilmi/JariyahMu/internal/response"
	"github.com/MirzaHilmi/JariyahMu/internal/validator"

	"github.com/pascaldekloe/jwt"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) signUserUp(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		FullName             string              `json:"fullName"`
		Email                string              `json:"email"`
		Password             string              `json:"password"`
		PasswordConfirmation string              `json:"passwordConfirmation"`
		Validator            validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	_, found, err := app.db.GetUserByEmail(payload.Email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	payload.Validator.CheckField(payload.Email != "", "Email", "Email is required")
	payload.Validator.CheckField(validator.Matches(payload.Email, validator.RgxEmail), "Email", "Must be a valid email address")
	payload.Validator.CheckField(!found, "Email", "Email is already in use")

	payload.Validator.CheckField(payload.Password != "", "Password", "Password is required")
	payload.Validator.CheckField(len(payload.Password) >= 8, "Password", "Password is too short")
	payload.Validator.CheckField(len(payload.Password) <= 72, "Password", "Password is too long")
	payload.Validator.CheckField(validator.NotIn(payload.Password, password.CommonPasswords...), "Password", "Password is too common")

	if payload.Validator.HasErrors() {
		app.failedValidation(w, r, payload.Validator)
		return
	}

	hashedPassword, err := password.Hash(payload.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	_, err = app.db.InsertUser(payload.Email, hashedPassword)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) logUserIn(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email     string              `json:"Email"`
		Password  string              `json:"Password"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, found, err := app.db.GetUserByEmail(payload.Email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	payload.Validator.CheckField(payload.Email != "", "Email", "Email is required")
	payload.Validator.CheckField(found, "Email", "Email address could not be found")

	if found {
		passwordMatches, err := password.Matches(payload.Password, user.HashedPassword)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		payload.Validator.CheckField(payload.Password != "", "Password", "Password is required")
		payload.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
	}

	if payload.Validator.HasErrors() {
		app.failedValidation(w, r, payload.Validator)
		return
	}

	var claims jwt.Claims
	claims.Subject = user.ID

	expiry := time.Now().Add(24 * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)

	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.baseURL}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]string{
		"AuthenticationToken":       string(jwtBytes),
		"AuthenticationTokenExpiry": expiry.Format(time.RFC3339),
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected handler"))
}
