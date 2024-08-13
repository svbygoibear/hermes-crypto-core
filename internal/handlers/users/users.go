package users

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	con "hermes-crypto-core/internal/constants"
	"hermes-crypto-core/internal/db"
	"hermes-crypto-core/internal/models"
)

// GetUsers handles GET requests to retrieve all users
func GetUsers(c *gin.Context) {
	users, err := db.DB.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser handles GET requests to retrieve a specific user (by id)
func GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := db.DB.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": con.USER_NOT_FOUND, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser handles POST requests to create a new user
func CreateUser(c *gin.Context) {
	var newUser models.User
	id := uuid.New() // Generate a new UUID for the user
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Checking to see if user already exists by email: %v", newUser.Email)
	// Check if user already exists by email
	user, err := db.DB.GetUserByEmail(newUser.Email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
	}
	log.Printf("We continue on at this point")
	// If user already exists, return the existing user
	if user != nil {
		c.JSON(http.StatusCreated, user)
		return
	}
	log.Printf("We have no user, thus we move")

	// If user does not exist, create a new user
	newUser.Id = id.String()
	newUser.Score = 0
	createdUser, err := db.DB.CreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdUser)
}

// UpdateUser handles PUT requests to update an existing user
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db.DB.UpdateUser(id, updatedUser, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE requests to remove a user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})
}
