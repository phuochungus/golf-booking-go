package entities

type AuthzRolePermission struct {
	BaseModel
	RoleID       int32 `gorm:"column:role_id"`
	PermissionID int32 `gorm:"column:permission_id"`
}

func (AuthzRolePermission) TableName() string {
	return "authz_role_permissions"
}
