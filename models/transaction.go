package models

type Transaction struct{
	ID *int64 `gorm:"primaryKey"`
	Correlation_id string `json:"correlation_id"`
	Is_root bool `json:"is_root"`
	UserId string `json:"user_id"`
	RootId string `json:"root_id"`
	ReqType string `json:"request_type"`
	Payload string `json:"payload"`
}