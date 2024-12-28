package model

type Order struct {
	UserID  int    `json:"user_id"`
	Product string `json:"product"`
	Count   int    `json:"count"`
	Address string `json:"address"`
}

type Status string

const (
	OrderStatusPending  Status = "pending"
	OrderStatusDone     Status = "done"
	OrderStatusCanceled Status = "canceled"
)
