package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hermes-crypto-core/internal/db"
	"hermes-crypto-core/internal/models"
)

// GetUsers handles GET requests to retrieve all users
func GetUsers(c *gin.Context) {
	users, err := db.DB.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserAndVotes handles GET requests to retrieve a specific user (by id) & their votes
// it will also check the latest vote result if 60 seconds have passed since the last vote
// and update the result if necessary
func GetUserAndVotes(c *gin.Context) {
	id := c.Param("id")
	user, err := db.DB.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUserAndVotes handles POST requests to create a new user including votes
// if there are any to create
func CreateUserAndVotes(c *gin.Context) {
	var newUser models.User
	id := uuid.New() // Generate a new UUID for the user
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser.Id = id.String()
	createdUser, err := db.DB.CreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
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

	user, err := db.DB.UpdateUser(id, updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE requests to remove a user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})
}
