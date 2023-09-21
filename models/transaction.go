package models

type Transaction struct{
	ID *int64 `gorm:"primaryKey"`
	IsRoot bool `json:"is_root"`
	Correlation_id string `json:"correlation_id"`
	CausationId string `json:"causation_id"`
	UserIp string `json:"user_ip"`
	UserId string `json:"user_id"`
	ReqType string `json:"request_type"`
	Is_Query string `json:"is_query"`
	Is_Command string `json:"is_Command"`
	Status string `json:"status"`
	Payload string `json:"payload"`
}