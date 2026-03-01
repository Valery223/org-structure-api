package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Valery223/org-structure-api/internal/domain"
)

type EmployeeHandler struct {
	svc domain.EmployeeService
}

func NewEmployeeHandler(svc domain.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{svc: svc}
}

type createEmpRequest struct {
	FullName string  `json:"full_name"`
	Position string  `json:"position"`
	HiredAt  *string `json:"hired_at"`
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	deptIDStr := r.PathValue("id")
	deptID, err := strconv.Atoi(deptIDStr)
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	var req createEmpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	emp, err := h.svc.Create(r.Context(), deptID, req.FullName, req.Position, req.HiredAt)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}
