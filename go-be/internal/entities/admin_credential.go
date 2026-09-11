package entities

type AdminCredential struct {
	AdminID        int32  `gorm:"column:admin_id;primaryKey"`
	HashedPassword string `gorm:"column:hashed_password"`
}

func (AdminCredential) TableName() string {
	return "admin_credentials"
}
