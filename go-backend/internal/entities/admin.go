package entities

type Admin struct {
	BaseModel
	Email          string `gorm:"column:email"`
	Password       string `gorm:"column:password"`
	Name           string `gorm:"column:name"`
	OrganizationID int32  `gorm:"column:organization_id"`
}

func (Admin) TableName() string {
	return "admins"
}
