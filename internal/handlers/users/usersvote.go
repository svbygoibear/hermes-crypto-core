package users

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"hermes-crypto-core/internal/coin"
	con "hermes-crypto-core/internal/constants"
	"hermes-crypto-core/internal/db"
	"hermes-crypto-core/internal/models"
)

// GetUserVotes handles GET requests to retrieve the specified (by id) user's votes
func GetUserVotesById(c *gin.Context) {
	log.Default().Println("Getting user votes")
	id := c.Param("id")
	user, err := db.DB.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": con.USER_NOT_FOUND})
		return
	}

	c.JSON(http.StatusOK, user.Votes)
}

// GetLastUserVoteResult handles GET requests to retrieve the specified (by id) user's last vote result
// If there is a result, we update this here as well with the users score
func GetLastUserVoteResult(c *gin.Context) {
	id := c.Param("id")
	user, err := db.DB.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": con.USER_NOT_FOUND, "message": err.Error()})
		return
	}

	log.Default().Println("Getting last user vote result")

	latestVote := GetLatestVote(*user)
	// Check if there is a latest vote for us to check results against

	log.Default().Println("Latest vote=", latestVote)
	if latestVote != nil {
		isYoungerThanTimeout := time.Since(latestVote.VoteDateTime.Time) < 60*time.Second
		log.Default().Println("Is recent=", isYoungerThanTimeout)
		if isYoungerThanTimeout {
			c.JSON(http.StatusOK, latestVote)
			return
		}

		// We have a coin value - it is all fine
		if latestVote.CoinValue != 0 {
			c.JSON(http.StatusOK, latestVote)
			return
		}

		// If the vote is not recent, we need to check the exchange rate and update the vote
		currentExchangeRate, err := coin.BinanceGetCurrentExchangeRate()
		if err != nil {
			c.JSON(http.StatusFailedDependency, gin.H{"error": "Could not determine current exchange rate", "message": err.Error()})
			return
		}
		latestVote.CoinValue = *currentExchangeRate

		log.Printf("Current exchange=$%f", latestVote.CoinValue)

		// Find and update the vote in the user's votes array
		for i, vote := range user.Votes {
			if vote.VoteDateTime == latestVote.VoteDateTime {
				user.Votes[i] = *latestVote
				break
			}
		}

		// Update the user score
		currentUserScore := user.Score
		// If the value at vote is higher than the current value, then the value dropped
		didCoinValueDrop := latestVote.CoinValueAtVote > latestVote.CoinValue
		// Now determine how the user score should be updated given their last vote direction
		if didCoinValueDrop {
			// Thus if the user voted down, they get a point
			if latestVote.VoteDirection == "down" {
				currentUserScore += 1
			} else {
				currentUserScore -= 1
			}
		} else {
			if latestVote.VoteDirection == "up" {
				currentUserScore += 1
			} else {
				currentUserScore -= 1
			}
		}
		user.Score = currentUserScore

		// Update the user with the new vote
		updatedUser, err := db.DB.UpdateUser(id, *user, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": con.USER_VOTE_UPDATE_FAILED, "message": err.Error()})
			return
		}
		log.Printf("Updated user vote for %v", updatedUser)

		// Return the updated vote
		c.JSON(http.StatusOK, latestVote)
		return
	}

	// Otherwise there is nothing to return
	c.JSON(http.StatusOK, nil)
}

// CreateUserVote handles POST requests to create a new user vote. We also run validation to see
// if there is an ongoing vote already before going ahead and creating a new vote.
func CreateUserVote(c *gin.Context) {
	id := c.Param("id")
	var newVote models.Vote
	if err := c.ShouldBindJSON(&newVote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	user, err := db.DB.GetUserByID(id)
	// If user does not exist, return an error since we can't add a vote to a non-existent user
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": con.USER_NOT_FOUND, "message": err.Error()})
		return
	}

	// If user already exists, check if there is an ongoing vote
	// Get the latest vote
	latestVote := GetLatestVote(*user)
	// If there is a vote, and it is recent, return an error noting there is an ongoing vote already
	if latestVote != nil {
		isRecent := time.Since(latestVote.VoteDateTime.Time) < 60*time.Second
		// if there is an unresolved vote, return an error
		if isRecent {
			c.JSON(http.StatusConflict, gin.H{"error": "User already has an ongoing vote"})
			return
		}
		// if there is an unchecked vote, return an error
		if latestVote.CoinValue == 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Users last vote has not been resolved"})
			return
		}
	}

	currentExchangeRate, err := coin.BinanceGetCurrentExchangeRate()
	if err != nil {
		c.JSON(http.StatusFailedDependency, gin.H{"error": "Could not determine current exchange rate"})
		return
	}
	// Set the current exchange rate as the value of the coin at the time of the vote
	newVote.CoinValueAtVote = *currentExchangeRate
	// Add default values for the vote
	newVote.VoteDateTime = models.TimestampTime{Time: time.Now()}
	newVote.CoinValue = 0
	newVote.VoteCoin = con.COIN_TYPE_BTC
	newVote.CoinValueCurrency = con.COIN_CURRENCY_USD

	// If there is no ongoing vote, create a new vote
	user.Votes = append(user.Votes, newVote)

	// Update the user with the extra votes
	updatedUser, err := db.DB.UpdateUser(id, *user, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": con.USER_VOTE_UPDATE_FAILED, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, updatedUser.Votes)
}

func GetLatestVote(user models.User) *models.Vote {
	if len(user.Votes) == 0 {
		return nil
	}

	newestVote := &user.Votes[0]
	newestTime := newestVote.VoteDateTime.Time

	for i := 1; i < len(user.Votes); i++ {
		currentVote := &user.Votes[i]
		currentTime := currentVote.VoteDateTime.Time

		if currentTime.After(newestTime) {
			newestVote = currentVote
			newestTime = currentTime
		}
	}

	return newestVote
}
