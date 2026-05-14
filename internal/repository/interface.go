package repository

import "github.com/konmitnik/hitalent_test/internal/models"

type RepositoryInterface interface {
	CreateDepartment(dep *models.Department) error
	GetDepartmentById(id uint) (*models.Department, error)
	SaveDepartment(dep *models.Department) error
	DeleteDepartmentById(id uint) error
	GetDepartmentByName(name string, parentId *uint) (*models.Department, error)
	LoadDepartmentChildren(dep *models.Department, depth int)
	LoadDepartmentEmployees(dep *models.Department)
	FitForParent(newParentId uint, depId uint) (bool, error)
	ReassignEmployeesAndDelete(depId uint, targetId uint) error

	CreateEmployee(emp *models.Employee) error
	GetEmployee(id uint) (*models.Employee, error)
}
