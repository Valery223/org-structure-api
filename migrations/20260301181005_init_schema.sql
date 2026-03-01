-- +goose Up
-- +goose StatementBegin
CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL CHECK (name <> ''),
    parent_id INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_department_parent FOREIGN KEY (parent_id) 
        REFERENCES departments(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_departments_name_parent ON departments (name, coalesce(parent_id, 0));

CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    department_id INT NOT NULL,
    full_name VARCHAR(200) NOT NULL CHECK (full_name <> ''),
    position VARCHAR(200) NOT NULL CHECK (position <> ''),
    hired_at DATE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_employee_department FOREIGN KEY (department_id) 
        REFERENCES departments(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;
-- +goose StatementEnd
