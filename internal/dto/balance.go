package dto

type Balance struct {
	Balance  float32 `json:"current"`
	Withdraw float32 `json:"withdrawn"`
}
