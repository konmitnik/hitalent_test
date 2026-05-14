package repository

import (
	"github.com/konmitnik/hitalent_test/internal/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateEmployee(emp *models.Employee) error {
	return gorm.G[models.Employee](r.conn).Create(r.ctx, emp)
}

func (r *Repository) GetEmployee(id uint) (*models.Employee, error) {
	emp, err := gorm.G[models.Employee](r.conn).Where("id = ?", id).First(r.ctx)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}
