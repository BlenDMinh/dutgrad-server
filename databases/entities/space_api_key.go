package entities

type SpaceAPIKey struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SpaceID     uint   `json:"space_id"`
	Space       *Space `json:"space" gorm:"foreignKey:SpaceID"`
}

func (s SpaceAPIKey) GetIdType() string {
	return "uint"
}
