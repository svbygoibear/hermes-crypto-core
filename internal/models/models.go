package models

type HealthCheck struct {
	Status string `json:"status" example:"ok"`
	Api    string `json:"api" example:"users"`
}

type Vote struct {
	VoteDirection     string        `json:"vote_direction" example:"up" enums:"up,down"`
	VoteDateTime      TimestampTime `json:"vote_date_time" swaggertype:"primitive,string" example:"2019-10-12T07:20:50.52Z"`
	VoteCoin          string        `json:"vote_coin" example:"bitcoin"`
	CoinValue         float64       `json:"coin_value" example:"1000"`
	CoinValueCurrency string        `json:"coin_value_currency" example:"usd"`
}

type User struct {
	Id    string `json:"id" example:"78712300234"` // Partition key
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"test@test.com"` // Sort key
	Votes []Vote `json:"votes"`
}

type ReturnError struct {
	ErrorMessage string `json:"error_message" example:"Failed to retrieve data from CoinGecko API"`
}

// Error implements error.
func (r ReturnError) Error() string {
	return r.ErrorMessage
}
