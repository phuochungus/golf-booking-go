package po

type AuthzObject struct {
	BaseModel
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

func (AuthzObject) TableName() string {
	return "authz_objects"
}
