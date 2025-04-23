package entities

import "time"

type Tier struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	SpaceLimit        int       `json:"space_limit" gorm:"not null"`
	DocumentLimit     int       `json:"document_limit" gorm:"not null"`
	QueryHistoryLimit int       `json:"query_history_limit" gorm:"not null"`
	QueryLimit        int       `json:"query_limit" gorm:"not null"`
	FileSizeLimitKb   int       `json:"file_size_limit_kb" gorm:"default:5120"`
	ApiCallLimit      int       `json:"api_call_limit" gorm:"default:100"`
	CostMonth         float64   `json:"cost_month" gorm:"type:decimal(10,2);not null"`
	Discount          float64   `json:"discount" gorm:"type:decimal(5,2);default:0"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Users             []User    `json:"users" gorm:"foreignKey:TierID"`
}

func (t Tier) GetIdType() string {
	return "uint"
}
