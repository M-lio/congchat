package dto

type TransferRequest struct {
	SenderID    uint    `json:"sender_id"`
	ReceiverID  uint    `json:"receiver_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description,omitempty"`
	Status      string  // "Pending", "Completed", "Refunded"
}
