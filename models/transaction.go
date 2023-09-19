package models

type Transaction struct{
	ID *int64 `gorm:"primaryKey"`
	Correlation_id string `json:"correlation_id"`
	ReqType string `json:"request_type"`
	Payload string `json:"payload"`
}