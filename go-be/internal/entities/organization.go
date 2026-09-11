package entities

type Organization struct {
	BaseModel
	Name string `gorm:"column:name"`
}

func (Organization) TableName() string {
	return "organizations"
}
