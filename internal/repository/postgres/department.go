package postgres

import (
	"context"
	"errors"

	"github.com/Valery223/org-structure-api/internal/domain"
	"gorm.io/gorm"
)

type departmentRepo struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) domain.DepartmentRepository {
	return &departmentRepo{db: db}
}

func (r *departmentRepo) Create(ctx context.Context, dept *domain.Department) error {
	return r.db.WithContext(ctx).Create(dept).Error
}

func (r *departmentRepo) GetByID(ctx context.Context, id int) (*domain.Department, error) {
	var dept domain.Department
	err := r.db.WithContext(ctx).First(&dept, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &dept, nil
}

func (r *departmentRepo) Update(ctx context.Context, dept *domain.Department) error {
	return r.db.WithContext(ctx).Model(dept).Select("*").Updates(dept).Error
}

func (r *departmentRepo) DeleteCascade(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&domain.Department{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *departmentRepo) DeleteReassign(ctx context.Context, id int, reassignToID int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&domain.Department{}).Where("id = ?", reassignToID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errors.New("reassign target department does not exist")
		}

		if err := tx.Model(&domain.Employee{}).
			Where("department_id = ?", id).
			Update("department_id", reassignToID).Error; err != nil {
			return err
		}

		if err := tx.Delete(&domain.Department{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *departmentRepo) ExistsByNameAndParent(ctx context.Context, name string, parentID *int) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&domain.Department{}).Where("name = ?", name)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *departmentRepo) GetParentID(ctx context.Context, id int) (*int, error) {
	var dept domain.Department
	err := r.db.WithContext(ctx).Select("parent_id").First(&dept, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return dept.ParentID, nil
}

func (r *departmentRepo) GetChildrenByParentID(ctx context.Context, parentID int) ([]domain.Department, error) {
	var children []domain.Department
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&children).Error
	return children, err
}
