package models

// HealthCheck is a struct that represents the health check response
type HealthCheck struct {
	Status string `json:"status" example:"ok"`
	Api    string `json:"api" example:"users"`
}

// Represents an individual vote
type Vote struct {
	VoteDirection     string        `json:"vote_direction" example:"up" enums:"up,down"`
	VoteDateTime      TimestampTime `json:"vote_date_time" swaggertype:"primitive,string" example:"2019-10-12T07:20:50.52Z"`
	VoteCoin          string        `json:"vote_coin" example:"bitcoin"`
	CoinValue         float64       `json:"coin_value" example:"58950.000000"`
	CoinValueAtVote   float64       `json:"coin_value_at_vote" example:"58940.000000"`
	CoinValueCurrency string        `json:"coin_value_currency" example:"USD"`
}

// User is a struct that represents a user with all of their votes
type User struct {
	Id    string `json:"id" example:"78712300234"` // Partition key
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"test@test.com"` // Sort key
	Score int    `json:"score" example:"0"`
	Votes []Vote `json:"votes"`
}

// CoinResult is a struct that represents the result of a coin query
type CoinResult struct {
	Coin              string        `json:"vote_coin" example:"bitcoin"`
	CoinValue         float64       `json:"coin_value" example:"58950.000000"`
	CoinValueCurrency string        `json:"coin_value_currency" example:"USD"`
	QueryTime         TimestampTime `json:"query_time" example:"2021-10-12T07:20:50.52Z"`
}

type ReturnError struct {
	ErrorMessage string `json:"error_message" example:"Failed to retrieve data from CoinGecko API"`
}

// Error implements error.
func (r ReturnError) Error() string {
	return r.ErrorMessage
}
