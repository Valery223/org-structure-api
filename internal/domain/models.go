package domain

import "time"

type Department struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	ParentID  *int      `json:"parent_id"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	Employees []Employee   `json:"employees,omitempty" gorm:"foreignKey:DepartmentID"`
	Children  []Department `json:"children,omitempty" gorm:"-"`
}

type Employee struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	DepartmentID int        `json:"department_id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at" gorm:"type:date"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
}
