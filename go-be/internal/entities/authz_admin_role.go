package entities

type AuthzAdminRole struct {
	BaseModel
	AdminID int32 `gorm:"column:admin_id"`
	RoleID  int32 `gorm:"column:role_id"`
}

func (AuthzAdminRole) TableName() string {
	return "authz_admin_roles"
}
