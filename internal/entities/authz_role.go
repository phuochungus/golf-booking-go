package po

type AuthzRole struct {
	BaseModel
	Name           string `gorm:"column:name"`
	OrganizationID int32  `gorm:"column:organization_id"`
	FacilityID     *int32 `gorm:"column:facility_id"`
}

func (AuthzRole) TableName() string {
	return "authz_roles"
}
