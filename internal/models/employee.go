package models

import "time"

type Employee struct {
	Id           uint        `gorm:"primaryKey" json:"id"`
	DepartmentId uint        `gorm:"not null;index" json:"department_id"`
	FullName     string      `gorm:"size:200;not null" json:"full_name"`
	Position     string      `gorm:"size:200;not null" json:"position"`
	HiredAt      *time.Time  `gorm:"type:date" json:"hired_at"`
	CreatedAt    time.Time   `gorm:"autoCreateTime;not null" json:"created_at"`
	Department   *Department `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}
