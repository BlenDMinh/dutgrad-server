package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type UserMFA struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	UserID      uint        `json:"user_id"`
	User        User        `json:"-" gorm:"foreignKey:UserID"`
	Secret      string      `json:"secret"`
	BackupCodes BackupCodes `json:"backup_codes" gorm:"type:jsonb"`
	Verified    bool        `json:"verified"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (u UserMFA) GetIdType() string {
	return "uint"
}

type BackupCodes []string

func (bc BackupCodes) Value() (driver.Value, error) {
	if bc == nil {
		return "[]", nil
	}
	return json.Marshal(bc)
}

func (bc *BackupCodes) Scan(value interface{}) error {
	if value == nil {
		*bc = BackupCodes{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal backup codes")
	}

	return json.Unmarshal(bytes, bc)
}
