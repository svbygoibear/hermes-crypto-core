package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	con "hermes-crypto-core/internal/constants"
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

func (m *MockDB) UpdateUser(id string, user models.User, updateScore bool) (*models.User, error) {
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

// Users Tests
func TestGetUsers(t *testing.T) {
	r, mockDB := setupTestRouter()
	r.GET("/users", GetUsers)

	mockUsers := []models.User{
		{Id: "1", Name: "Test User", Email: "test@test.com", Votes: []models.Vote{
			{VoteDirection: "up", CoinValue: 0.5, CoinValueAtVote: 0.5, CoinValueCurrency: con.COIN_CURRENCY_USD, VoteCoin: con.COIN_TYPE_BTC, VoteDateTime: models.TimestampTime{time.Time{}}}}}}
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

func TestGetUser(t *testing.T) {
	r, mockDB := setupTestRouter()
	r.GET("/users/:id", GetUser)

	mockUser := &models.User{Id: "1", Name: "Test User"}
	mockDB.On("GetUserByID", "1").Return(mockUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, *mockUser, response)
}

func TestCreateUser(t *testing.T) {
	r, mockDB := setupTestRouter()
	r.POST("/users", CreateUser)

	newUser := models.User{Name: "New User", Email: "test@test.com"}
	mockDB.On("GetUserByEmail", newUser.Email).Return(&newUser, nil)
	mockDB.On("CreateUser", newUser).Return(&newUser, nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
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
	mockDB.On("UpdateUser", "1", updatedUser, false).Return(&updatedUser, nil)

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

// Votes Tests
func TestGetUserVotes(t *testing.T) {
	r, mockDB := setupTestRouter()
	r.GET("/users/:id/votes", GetUserVotesById)

	mockUser := &models.User{Id: "1", Name: "Test User", Votes: []models.Vote{{VoteDirection: "up"}, {VoteDirection: "down"}}}
	mockDB.On("GetUserByID", "1").Return(mockUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/votes", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response []models.Vote
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, *&mockUser.Votes, response)
}

func TestGetUserLastVoteResult(t *testing.T) {
	r, mockDB := setupTestRouter()
	r.GET("/users/:id/votes/result", GetLastUserVoteResult)

	voteDateTime1, _ := time.Parse(time.RFC3339, "2023-10-12T07:20:50.52Z")
	voteDateTime2, _ := time.Parse(time.RFC3339, "2024-01-01T07:20:50.52Z")
	mockUser := &models.User{Id: "18890123000123", Name: "Test User", Email: "test@gmail.com", Votes: []models.Vote{
		{VoteDirection: "up", CoinValue: 59760, CoinValueAtVote: 45234, CoinValueCurrency: con.COIN_CURRENCY_USD, VoteCoin: con.COIN_TYPE_BTC, VoteDateTime: models.TimestampTime{Time: voteDateTime1}},
		{VoteDirection: "down", CoinValue: 59760, CoinValueAtVote: 45234, CoinValueCurrency: con.COIN_CURRENCY_USD, VoteCoin: con.COIN_TYPE_BTC, VoteDateTime: models.TimestampTime{Time: voteDateTime2}}},
	}
	mockDB.On("GetUserByID", "18890123000123").Return(mockUser, nil)
	mockDB.On("UpdateUser", "18890123000123", *mockUser, true).Return(mockUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/18890123000123/votes/result", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response models.Vote
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.NotEqual(t, *&mockUser.Votes[1], response)
}

func TestGetUserLastVoteResultNoValue(t *testing.T) {
	r, mockDB := setupTestRouter()
	r.GET("/users/:id/votes/result", GetLastUserVoteResult)

	voteDateTime1, _ := time.Parse(time.RFC3339, "2023-10-12T07:20:50.52Z")
	voteDateTime2, _ := time.Parse(time.RFC3339, "2024-01-01T19:30:50.52Z")
	mockUser := &models.User{Id: "18890123000123", Name: "Test User", Email: "test@gmail.com", Votes: []models.Vote{
		{VoteDirection: "up", CoinValue: 59760, CoinValueAtVote: 45234, CoinValueCurrency: con.COIN_CURRENCY_USD, VoteCoin: con.COIN_TYPE_BTC, VoteDateTime: models.TimestampTime{Time: voteDateTime1}},
		{VoteDirection: "down", CoinValue: 0, CoinValueAtVote: 45234, CoinValueCurrency: con.COIN_CURRENCY_USD, VoteCoin: con.COIN_TYPE_BTC, VoteDateTime: models.TimestampTime{Time: voteDateTime2}}},
	}
	mockDB.On("GetUserByID", "18890123000123").Return(mockUser, nil)
	mockDB.On("UpdateUser", "18890123000123", *mockUser, true).Return(mockUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/18890123000123/votes/result", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response models.Vote
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.NotEqual(t, *&mockUser.Votes[1], response)
}
