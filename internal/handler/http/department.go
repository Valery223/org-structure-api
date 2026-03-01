package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Valery223/org-structure-api/internal/domain"
)

type DepartmentHandler struct {
	svc domain.DepartmentService
}

func NewDepartmentHandler(svc domain.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{svc: svc}
}

type createDeptRequest struct {
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}

type updateDeptRequest struct {
	Name     *string `json:"name"`
	ParentID *int    `json:"parent_id"`
}

func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createDeptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	dept, err := h.svc.Create(r.Context(), req.Name, req.ParentID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	depth := 1
	if d := r.URL.Query().Get("depth"); d != "" {
		if val, err := strconv.Atoi(d); err == nil && val >= 0 && val <= 5 {
			depth = val
		}
	}

	includeEmps := true
	if inc := r.URL.Query().Get("include_employees"); inc == "false" {
		includeEmps = false
	}

	dept, err := h.svc.GetTree(r.Context(), id, depth, includeEmps)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req updateDeptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	dept, err := h.svc.Update(r.Context(), id, req.Name, req.ParentID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		http.Error(w, "mode is required", http.StatusBadRequest)
		return
	}

	var reassignID *int
	if mode == "reassign" {
		ridStr := r.URL.Query().Get("reassign_to_department_id")
		rid, err := strconv.Atoi(ridStr)
		if err != nil {
			http.Error(w, "reassign_to_department_id is required and must be int", http.StatusBadRequest)
			return
		}
		reassignID = &rid
	}

	if err := h.svc.Delete(r.Context(), id, mode, reassignID); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrDuplicateName), errors.Is(err, domain.ErrSelfParenting), errors.Is(err, domain.ErrCycleDetected):
		http.Error(w, err.Error(), http.StatusConflict) // 409
	case errors.Is(err, domain.ErrInvalidValidation), errors.Is(err, domain.ErrDepartmentNotSet):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
