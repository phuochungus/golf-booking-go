package entities

type Admin struct {
	BaseModel
	Email           string `gorm:"column:email"`
	Name            string `gorm:"column:name"`
	OrganizationID  *int32 `gorm:"column:organization_id"`
	AdminCredential *AdminCredential
}

func (Admin) TableName() string {
	return "admins"
}
