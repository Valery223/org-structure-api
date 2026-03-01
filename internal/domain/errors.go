package domain

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrCycleDetected     = errors.New("cycle detected in department tree")
	ErrSelfParenting     = errors.New("department cannot be parent of itself")
	ErrDuplicateName     = errors.New("department name must be unique within the same parent")
	ErrInvalidValidation = errors.New("validation failed")
	ErrDepartmentNotSet  = errors.New("reassign department id must be provided for reassign mode")
)
