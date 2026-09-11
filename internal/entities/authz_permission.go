package entities

type AuthzPermission struct {
	BaseModel
	ObjectID int32 `gorm:"column:object_id"`
	ActionID int32 `gorm:"column:action_id"`
}

func (AuthzPermission) TableName() string {
	return "authz_permissions"
}
