package entities

type Organization struct {
	BaseModel
	Name        string `gorm:"column:name"`
	RootAdminID *int32 `gorm:"column:root_admin_id"`
}

func (Organization) TableName() string {
	return "organizations"
}
