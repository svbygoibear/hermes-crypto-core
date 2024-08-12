package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"hermes-crypto-core/internal/db"
	"hermes-crypto-core/internal/models"
)

// MockDB is a mock of the db package
type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockDB) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) CreateUser(user models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) UpdateUser(id string, user models.User) (*models.User, error) {
	args := m.Called(id, user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter() (*gin.Engine, *MockDB) {
	r := gin.Default()
	mockDB := new(MockDB)
	db.DB = mockDB
	return r, mockDB
}

func TestGetUsers(t *testing.T) {
	r, mockDB := setupTestRouter()
	r.GET("/users", GetUsers)

	mockUsers := []models.User{{Id: "1", Name: "Test User"}}
	mockDB.On("GetAllUsers").Return(mockUsers, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response []models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, mockUsers, response)
}

func TestGetUserAndVotes(t *testing.T) {
	r, mockDB := setupTestRouter()
	r.GET("/users/:id/vote", GetUserAndVotes)

	mockUser := &models.User{Id: "1", Name: "Test User"}
	mockDB.On("GetUserByID", "1").Return(mockUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/vote", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, *mockUser, response)
}

func TestCreateUserAndVotes(t *testing.T) {
	r, mockDB := setupTestRouter()
	r.POST("/users/vote", CreateUserAndVotes)

	newUser := models.User{Name: "New User"}
	mockDB.On("CreateUser", newUser).Return(&newUser, nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/users/vote", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, newUser, response)
}

func TestUpdateUser(t *testing.T) {
	r, mockDB := setupTestRouter()
	r.PUT("/users/:id", UpdateUser)

	updatedUser := models.User{Id: "1", Name: "Updated User"}
	mockDB.On("UpdateUser", "1", updatedUser).Return(&updatedUser, nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, updatedUser, response)
}

func TestDeleteUser(t *testing.T) {
	r, mockDB := setupTestRouter()
	r.DELETE("/users/:id", DeleteUser)

	mockDB.On("DeleteUser", "1").Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "User successfully deleted", response["message"])
}
