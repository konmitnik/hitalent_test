package models

import "time"

type Department struct {
	Id        uint         `gorm:"primaryKey" json:"id"`
	Name      string       `gorm:"size:200;not null" json:"name"`
	ParentId  *uint        `gorm:"index" json:"parent_id"`
	CreatedAt time.Time    `gorm:"autoCreateTime;not null" json:"created_at"`
	Parent    *Department  `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Children  []Department `gorm:"foreignKey:ParentId" json:"children,omitempty"`
	Employees []Employee   `gorm:"foreignKey:DepartmentId" json:"employees,omitempty"`
}
