package service

import (
	"context"
	"testing"

	"github.com/Valery223/org-structure-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MOCKS

type MockDepartmentRepo struct {
	mock.Mock
}

func (m *MockDepartmentRepo) Create(ctx context.Context, dept *domain.Department) error {
	args := m.Called(ctx, dept)
	return args.Error(0)
}

func (m *MockDepartmentRepo) GetByID(ctx context.Context, id int) (*domain.Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Department), args.Error(1)
}

func (m *MockDepartmentRepo) Update(ctx context.Context, dept *domain.Department) error {
	args := m.Called(ctx, dept)
	return args.Error(0)
}

func (m *MockDepartmentRepo) DeleteCascade(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockDepartmentRepo) DeleteReassign(ctx context.Context, id, targetID int) error {
	args := m.Called(ctx, id, targetID)
	return args.Error(0)
}

func (m *MockDepartmentRepo) ExistsByNameAndParent(ctx context.Context, name string, parentID *int) (bool, error) {
	args := m.Called(ctx, name, parentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDepartmentRepo) GetChildrenByParentID(ctx context.Context, parentID int) ([]domain.Department, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]domain.Department), args.Error(1)
}

func (m *MockDepartmentRepo) GetParentID(ctx context.Context, id int) (*int, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

type MockEmployeeRepo struct {
	mock.Mock
}

func (m *MockEmployeeRepo) Create(ctx context.Context, emp *domain.Employee) error {
	return nil
}
func (m *MockEmployeeRepo) GetByDepartmentID(ctx context.Context, deptID int) ([]domain.Employee, error) {
	return nil, nil
}

// TESTS

func TestDepartmentService_Update_CycleDetection(t *testing.T) {
	mockDeptRepo := new(MockDepartmentRepo)
	mockEmpRepo := new(MockEmployeeRepo)
	service := NewDepartmentService(mockDeptRepo, mockEmpRepo)
	ctx := context.Background()

	// Сценарий: Попытка переместить Департамент А (ID=1) внутрь Департамента Б (ID=2),
	// при том, что Департамент Б уже является потомком А.
	// Структура: ID 1 -> ID 2. Мы пытаемся сделать ID 1 дочерним для ID 2.

	dept1 := &domain.Department{ID: 1, Name: "Dept 1", ParentID: nil}

	newParentID := 2

	mockDeptRepo.On("GetByID", ctx, 1).Return(dept1, nil)

	parentOf2 := 1
	mockDeptRepo.On("GetParentID", ctx, 2).Return(&parentOf2, nil)

	_, err := service.Update(ctx, 1, nil, &newParentID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCycleDetected, err)

	mockDeptRepo.AssertExpectations(t)
}

func TestDepartmentService_Create_Success(t *testing.T) {
	mockDeptRepo := new(MockDepartmentRepo)
	mockEmpRepo := new(MockEmployeeRepo)
	service := NewDepartmentService(mockDeptRepo, mockEmpRepo)
	ctx := context.Background()

	name := "New Dept"

	mockDeptRepo.On("ExistsByNameAndParent", ctx, name, (*int)(nil)).Return(false, nil)
	mockDeptRepo.On("Create", ctx, mock.AnythingOfType("*domain.Department")).Return(nil)

	res, err := service.Create(ctx, name, nil)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, name, res.Name)
}
