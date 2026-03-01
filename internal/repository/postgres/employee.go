package postgres

import (
	"context"

	"github.com/Valery223/org-structure-api/internal/domain"
	"gorm.io/gorm"
)

type employeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) domain.EmployeeRepository {
	return &employeeRepo{db: db}
}

func (r *employeeRepo) Create(ctx context.Context, emp *domain.Employee) error {
	return r.db.WithContext(ctx).Create(emp).Error
}

func (r *employeeRepo) GetByDepartmentID(ctx context.Context, deptID int) ([]domain.Employee, error) {
	var employees []domain.Employee
	err := r.db.WithContext(ctx).Where("department_id = ?", deptID).Find(&employees).Error
	return employees, err
}
