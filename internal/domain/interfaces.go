package domain

import "context"

// DepartmentRepository defines the interface for department data access operations.
type DepartmentRepository interface {
	// Create inserts a new department into the database.
	Create(ctx context.Context, dept *Department) error

	// GetByID retrieves a department by its unique identifier.
	GetByID(ctx context.Context, id int) (*Department, error)

	// Update modifies an existing department in the database.
	Update(ctx context.Context, dept *Department) error

	// DeleteCascade removes a department and all its sub-departments recursively.
	DeleteCascade(ctx context.Context, id int) error

	// DeleteReassign removes a department and reassigns its employees to another department.
	DeleteReassign(ctx context.Context, id int, reassignToID int) error

	// ExistsByNameAndParent checks if a department with the given name exists under the specified parent.
	ExistsByNameAndParent(ctx context.Context, name string, parentID *int) (bool, error)

	// GetParentID returns the parent ID of a department, or nil if it's a root department.
	GetParentID(ctx context.Context, id int) (*int, error)

	// GetChildrenByParentID retrieves all direct children (sub-departments) of a department.
	GetChildrenByParentID(ctx context.Context, parentID int) ([]Department, error)
}

// EmployeeRepository defines the interface for employee data access operations.
type EmployeeRepository interface {
	// Create inserts a new employee into the database.
	Create(ctx context.Context, emp *Employee) error

	// GetByDepartmentID retrieves all employees belonging to a specific department.
	GetByDepartmentID(ctx context.Context, deptID int) ([]Employee, error)
}

// DepartmentService defines the business logic interface for department operations.
type DepartmentService interface {
	// Create creates a new department with the given name and optional parent.
	Create(ctx context.Context, name string, parentID *int) (*Department, error)

	// GetTree retrieves a department with its nested children up to the specified depth.
	GetTree(ctx context.Context, id, depth int, includeEmployees bool) (*Department, error)

	// Update modifies department properties such as name or parent, performing cycle detection.
	Update(ctx context.Context, id int, name *string, parentID *int) (*Department, error)

	// Delete removes a department using the specified deletion mode (cascade or reassign).
	Delete(ctx context.Context, id int, mode string, reassignToID *int) error
}

// EmployeeService defines the business logic interface for employee operations.
type EmployeeService interface {
	// Create creates a new employee in the specified department.
	Create(ctx context.Context, deptID int, fullName, position string, hiredAt *string) (*Employee, error)
}
