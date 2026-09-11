package po

type AuthzAction struct {
	BaseModel
	Name                string `gorm:"column:name"`
	Description         string `gorm:"column:description"`
	IsOrganizationLevel bool   `gorm:"column:is_organization_level;default:false"`
	IsFacilityLevel     bool   `gorm:"column:is_facility_level;default:false"`
}

func (AuthzAction) TableName() string {
	return "authz_actions"
}
