package entities

type AuthzAdminPermission struct {
	BaseModel
	AdminID      int32 `gorm:"column:admin_id"`
	PermissionID int32 `gorm:"column:permission_id"`
}

func (AuthzAdminPermission) TableName() string {
	return "authz_admin_permissions"
}
