package po

type Organization struct {
	BaseModel
	Name  string `gorm:"column:name"`
	Phone string `gorm:"column:phone"`
	Email string `gorm:"column:email"`
}

func (Organization) TableName() string {
	return "organizations"
}
