package entities

type Facility struct {
	BaseModel
	Name           string   `gorm:"column:name"`
	OrganizationID int32    `gorm:"column:organization_id"`
	Rating         *float32 `gorm:"column:rating;default:0"`
}

func (Facility) TableName() string {
	return "facilities"
}
