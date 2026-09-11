package entities

type Admin struct {
	BaseModel
	Email          string `gorm:"column:email"`
	Password       string `gorm:"column:password"`
	Name           string `gorm:"column:name"`
	Root           bool   `gorm:"column:root;default:false"`
	OrganizationID int32  `gorm:"column:organization_id"`
}

func (Admin) TableName() string {
	return "admins"
}
