package service

import (
	"context"
	"strings"

	"github.com/Valery223/org-structure-api/internal/domain"
)

type departmentService struct {
	deptRepo domain.DepartmentRepository
	empRepo  domain.EmployeeRepository
}

func NewDepartmentService(deptRepo domain.DepartmentRepository, empRepo domain.EmployeeRepository) domain.DepartmentService {
	return &departmentService{
		deptRepo: deptRepo,
		empRepo:  empRepo,
	}
}

func (s *departmentService) Create(ctx context.Context, name string, parentID *int) (*domain.Department, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrInvalidValidation
	}

	exists, err := s.deptRepo.ExistsByNameAndParent(ctx, name, parentID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrDuplicateName
	}

	dept := &domain.Department{
		Name:     name,
		ParentID: parentID,
	}

	if err := s.deptRepo.Create(ctx, dept); err != nil {
		return nil, err
	}

	return dept, nil
}

func (s *departmentService) GetTree(ctx context.Context, id, depth int, includeEmployees bool) (*domain.Department, error) {
	root, err := s.deptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.buildTree(ctx, root, depth, includeEmployees); err != nil {
		return nil, err
	}

	return root, nil
}

// buildTree recursively populates the Children slice for the given department node.
// It stops when the specified depth is reached.
func (s *departmentService) buildTree(ctx context.Context, current *domain.Department, depth int, includeEmployees bool) error {
	if includeEmployees {
		emps, err := s.empRepo.GetByDepartmentID(ctx, current.ID)
		if err != nil {
			return err
		}
		current.Employees = emps
	}

	if depth <= 0 {
		return nil
	}

	children, err := s.deptRepo.GetChildrenByParentID(ctx, current.ID)
	if err != nil {
		return err
	}

	current.Children = children
	for i := range current.Children {
		if err := s.buildTree(ctx, &current.Children[i], depth-1, includeEmployees); err != nil {
			return err
		}
	}
	return nil
}

func (s *departmentService) Update(ctx context.Context, id int, name *string, parentID *int) (*domain.Department, error) {
	dept, err := s.deptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		cleanedName := strings.TrimSpace(*name)
		if cleanedName == "" {
			return nil, domain.ErrInvalidValidation
		}
		pidToCheck := dept.ParentID
		if parentID != nil {
			pidToCheck = parentID
		}

		exists, err := s.deptRepo.ExistsByNameAndParent(ctx, cleanedName, pidToCheck)
		if err != nil {
			return nil, err
		}
		if exists && cleanedName != dept.Name {
			return nil, domain.ErrDuplicateName
		}
		dept.Name = cleanedName
	}

	if parentID != nil {
		newParentID := *parentID
		if newParentID == id {
			return nil, domain.ErrSelfParenting
		}
		if err := s.validateCycle(ctx, id, newParentID); err != nil {
			return nil, err
		}
		dept.ParentID = &newParentID
	}

	if err := s.deptRepo.Update(ctx, dept); err != nil {
		return nil, err
	}

	return dept, nil
}

// validateCycle checks if moving a department creates a circular reference in the tree.
// It traverses up the tree from the newParentID. If we encounter the movingDeptID
// during the traversal, it means the new parent is currently a descendant of the moving department,
// which would close a loop.
func (s *departmentService) validateCycle(ctx context.Context, movingDeptID, newParentID int) error {
	currentID := newParentID

	for {
		if currentID == 0 {
			break
		}
		if currentID == movingDeptID {
			return domain.ErrCycleDetected
		}

		parentID, err := s.deptRepo.GetParentID(ctx, currentID)
		if err != nil {
			return err
		}
		if parentID == nil {
			break
		}
		currentID = *parentID
	}
	return nil
}

func (s *departmentService) Delete(ctx context.Context, id int, mode string, reassignToID *int) error {
	_, err := s.deptRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if mode == "cascade" {
		return s.deptRepo.DeleteCascade(ctx, id)
	} else if mode == "reassign" {
		if reassignToID == nil {
			return domain.ErrDepartmentNotSet
		}
		return s.deptRepo.DeleteReassign(ctx, id, *reassignToID)
	}

	return domain.ErrInvalidValidation
}
