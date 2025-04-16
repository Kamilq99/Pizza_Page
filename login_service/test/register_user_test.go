package handlers_test

import (
	"bytes"
	"encoding/json"
	"login_service/handlers"
	"login_service/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", handlers.RegisterUser)

	user := models.UserRegister{
		Username: "newuser",
		Password: "newpassword123",
		Email:    "newuser@example.com",
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Registration successful")
}

func TestRegisterUser_BadRequest_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", handlers.RegisterUser)

	invalidJSON := `{"username":"newuser", "password":"pass123"` // Brak nawiasu
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "error")
}

func TestRegisterUser_BadRequest_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", handlers.RegisterUser)

	// Brakuje pola email
	body := `{"username":"user123","password":"pass123"}`
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "error")
}

func TestRegisterUser_BadRequest_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", handlers.RegisterUser)

	req, _ := http.NewRequest(http.MethodPost, "/register", nil)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "error")
}
