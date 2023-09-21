package models

type Transaction struct{
	ID *int64 `gorm:"primaryKey"`
	Correlation_id int64 `json:"correlation_id"`
	CausationId int64 `json:"causation_id"`
	Message string `json:"message"`
	Status string `json:"status"`
}

//geteway = user request stored

//server a = proccecing user request -> event stored

//server b -> proccesing user request -> event stored

// gatewaye = save


// id cor_i