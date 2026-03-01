package app

import (
	"net/http"

	handler "github.com/Valery223/org-structure-api/internal/handler/http"
)

func SetupRouter(deptH *handler.DepartmentHandler, empH *handler.EmployeeHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// Department routes
	mux.HandleFunc("POST /departments/", deptH.Create)
	mux.HandleFunc("GET /departments/{id}", deptH.Get)
	mux.HandleFunc("PATCH /departments/{id}", deptH.Update)
	mux.HandleFunc("DELETE /departments/{id}", deptH.Delete)

	// Employee routes
	mux.HandleFunc("POST /departments/{id}/employees/", empH.Create)

	return mux
}
