package repository

import (
	"strings"

	"github.com/konmitnik/hitalent_test/internal/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateDepartment(dep *models.Department) error {
	return gorm.G[models.Department](r.conn).Create(r.ctx, dep)
}

func (r *Repository) GetDepartmentById(id uint) (*models.Department, error) {
	dep, err := gorm.G[models.Department](r.conn).Where("id = ?", id).First(r.ctx)
	if err != nil {
		return nil, err
	}
	return &dep, nil
}

func (r *Repository) SaveDepartment(dep *models.Department) error {
	return r.conn.Save(dep).Error
}

func (r *Repository) DeleteDepartmentById(id uint) error {
	_, err := gorm.G[models.Department](r.conn).Where("id = ?", id).Delete(r.ctx)
	return err
}

func (r *Repository) GetDepartmentByName(name string, parentId *uint) (*models.Department, error) {
	dep, err := gorm.G[models.Department](r.conn).Where(
		"name = ? AND parent_id IS NOT DISTINCT FROM ?",
		name,
		parentId,
	).First(r.ctx)
	if err != nil {
		return nil, err
	}
	return &dep, nil
}

func (r *Repository) LoadDepartmentChildren(dep *models.Department, depth int) {
	preloadChain := "Children" + strings.Repeat(".Children", depth-1)
	tempDep, err := gorm.G[models.Department](r.conn).
		Where("id = ?", dep.Id).
		Preload(preloadChain, nil).
		First(r.ctx)
	if err == nil {
		dep.Children = tempDep.Children
	}
}

func (r *Repository) LoadDepartmentEmployees(dep *models.Department) {
	tempDep, err := gorm.G[models.Department](r.conn).
		Where("id = ?", dep.Id).
		Preload("Employees", func(db gorm.PreloadBuilder) error {
			db.Order("created_at ASC")
			return nil
		}).
		First(r.ctx)
	if err == nil {
		dep.Employees = tempDep.Employees
	}
}

func (r *Repository) FitForParent(newParentId uint, depId uint) (bool, error) {
	if newParentId == depId {
		return false, nil
	}
	current := newParentId
	for i := 0; i < 100; i++ {
		d, err := gorm.G[models.Department](r.conn).Where("id = ?", current).First(r.ctx)
		if err != nil {
			return false, err
		}
		if d.ParentId == nil {
			return true, nil
		}
		if *d.ParentId == depId {
			return false, nil
		}
		current = *d.ParentId
	}
	return false, nil
}

func (r *Repository) ReassignEmployeesAndDelete(depId uint, targetId uint) error {
	return r.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Employee{}).
			Where("department_id = ?", depId).
			Update("department_id", targetId).Error; err != nil {
			return err
		}
		_, err := gorm.G[models.Department](tx).Where("id = ?", depId).Delete(r.ctx)
		return err
	})
}
