package po

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains the columns shared by all persisted entities.
type BaseModel struct {
	ID        int32          `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}
