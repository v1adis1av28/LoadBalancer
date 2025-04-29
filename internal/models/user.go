package models

type User struct {
	ClientID   string `json:"client_id"`
	Capacity   int    `json:"capacity"`
	RefillRate int    `json:"rate_per_sec"`
}
