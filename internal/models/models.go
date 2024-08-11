package models

type HealthCheck struct {
	Status string `json:"status"`
	Api    string `json:"api"`
}

type Vote struct {
	VoteDirection     string  `json:"vote_direction"`
	VoteDateTime      string  `json:"vote_date_time"`
	VoteCoin          string  `json:"vote_coin"`
	CoinValue         float64 `json:"coin_value"`
	CoinValueCurrency string  `json:"coin_value_currency"`
}

type User struct {
	ID    string `json:"id"` // Partition key
	Name  string `json:"name"`
	Email string `json:"email"`
	Votes []Vote `json:"votes"`
}
