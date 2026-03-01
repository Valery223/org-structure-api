package service

import (
	"context"
	"strings"
	"time"

	"github.com/Valery223/org-structure-api/internal/domain"
)

type employeeService struct {
	empRepo  domain.EmployeeRepository
	deptRepo domain.DepartmentRepository
}

func NewEmployeeService(empRepo domain.EmployeeRepository, deptRepo domain.DepartmentRepository) domain.EmployeeService {
	return &employeeService{
		empRepo:  empRepo,
		deptRepo: deptRepo,
	}
}

func (s *employeeService) Create(ctx context.Context, deptID int, fullName, position string, hiredAtStr *string) (*domain.Employee, error) {
	fullName = strings.TrimSpace(fullName)
	position = strings.TrimSpace(position)

	if fullName == "" || position == "" {
		return nil, domain.ErrInvalidValidation
	}

	_, err := s.deptRepo.GetByID(ctx, deptID)
	if err != nil {
		return nil, err
	}

	var hiredAt *time.Time
	if hiredAtStr != nil {
		parsed, err := time.Parse("2006-01-02", *hiredAtStr)
		if err != nil {
			return nil, domain.ErrInvalidValidation
		}
		hiredAt = &parsed
	}

	emp := &domain.Employee{
		DepartmentID: deptID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
	}

	if err := s.empRepo.Create(ctx, emp); err != nil {
		return nil, err
	}

	return emp, nil
}
