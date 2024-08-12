package users

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hermes-crypto-core/internal/db"
	"hermes-crypto-core/internal/models"
)

// GetUserVotes handles GET requests to retrieve the specified (by id) user's votes
func GetUserVotes(c *gin.Context) {
	id := c.Param("id")
	user, err := db.DB.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		return
	}
	// TODO: marshal user votes
	c.JSON(http.StatusOK, user)
}

// GetLastUserVoteResult handles GET requests to retrieve the specified (by id) user's last vote result
// If there is a result, we update this here as well
func GetLastUserVoteResult(c *gin.Context) {
	id := c.Param("id")
	user, err := db.DB.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		return
	}
	// TODO: marshal user votes
	c.JSON(http.StatusOK, user)
}

// CreateUserVote handles POST requests to create a new user vote. We also run validation to see
// if there is an ongoing vote already before going ahead and creating a new vote.
func CreateUserVote(c *gin.Context) {
	var newUser models.User
	id := uuid.New() // Generate a new UUID for the user
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists by email
	user, err := db.DB.GetUserByEmail(newUser.Email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
	}

	// If user already exists, return the existing user
	if user != nil {
		c.JSON(http.StatusCreated, user)
		return
	}

	// If user does not exist, create a new user
	newUser.Id = id.String()
	createdUser, err := db.DB.CreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, createdUser)
}
