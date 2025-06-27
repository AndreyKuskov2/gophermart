package models

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}
