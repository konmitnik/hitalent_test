package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/konmitnik/hitalent_test/internal/helpers"
	"github.com/konmitnik/hitalent_test/internal/models"
	"github.com/konmitnik/hitalent_test/internal/repository"
)

type Handler struct {
	rep *repository.Repository
	log *slog.Logger
}

func NewHandler(rep *repository.Repository) *Handler {
	return &Handler{rep: rep, log: slog.Default()}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	rawPath := strings.TrimPrefix(r.URL.Path, "/departments")
	rawPath = strings.Trim(rawPath, "/")

	var parts []string
	if rawPath != "" {
		parts = strings.Split(rawPath, "/")
	}

	h.log.Info("request", "method", r.Method, "path", r.URL.Path)

	switch r.Method {
	case http.MethodPost:
		if len(parts) == 0 {
			h.CreateDepartment(w, r)
		} else if len(parts) == 2 && parts[1] == "employees" {
			h.CreateEmployee(w, r, parts[0])
		} else {
			http.NotFound(w, r)
		}
	case http.MethodGet:
		if len(parts) == 1 {
			if dep := h.departmentOrNotFound(w, parts[0]); dep != nil {
				h.GetDepartment(w, r, dep)
			}
		} else {
			http.NotFound(w, r)
		}
	case http.MethodPatch:
		if len(parts) == 1 {
			if dep := h.departmentOrNotFound(w, parts[0]); dep != nil {
				h.PatchDepartment(w, r, dep)
			}
		} else {
			http.NotFound(w, r)
		}
	case http.MethodDelete:
		if len(parts) == 1 {
			if dep := h.departmentOrNotFound(w, parts[0]); dep != nil {
				h.DeleteDepartment(w, r, dep)
			}
		} else {
			http.NotFound(w, r)
		}
	default:
		helpers.ResponseJSON(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) departmentOrNotFound(w http.ResponseWriter, idStr string) *models.Department {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		helpers.ResponseJSON(w, map[string]string{"error": "invalid department id"}, http.StatusBadRequest)
		return nil
	}
	dep, err := h.rep.GetDepartmentById(uint(id))
	if err != nil {
		helpers.ResponseJSON(w, map[string]string{"error": "department not found"}, http.StatusNotFound)
		return nil
	}
	return dep
}

// POST /departments/
func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		ParentId *uint  `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.ResponseJSON(w, map[string]string{"error": "invalid request body"}, http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || len(req.Name) > 200 {
		helpers.ResponseJSON(w, map[string]string{"error": "name must be between 1 and 200 characters"}, http.StatusBadRequest)
		return
	}

	if req.ParentId != nil {
		if _, err := h.rep.GetDepartmentById(*req.ParentId); err != nil {
			helpers.ResponseJSON(w, map[string]string{"error": "parent department not found"}, http.StatusNotFound)
			return
		}
	}

	if dup, _ := h.rep.GetDepartmentByName(req.Name, req.ParentId); dup != nil {
		helpers.ResponseJSON(w, map[string]string{"error": "department with this name already exists under the same parent"}, http.StatusConflict)
		return
	}

	dep := &models.Department{Name: req.Name, ParentId: req.ParentId}
	if err := h.rep.CreateDepartment(dep); err != nil {
		h.log.Error("create department", "err", err)
		helpers.ResponseJSON(w, map[string]string{"error": "internal server error"}, http.StatusInternalServerError)
		return
	}

	helpers.ResponseJSON(w, dep, http.StatusCreated)
}

// POST /departments/{id}/employees/
func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request, depIdStr string) {
	depId, err := strconv.ParseUint(depIdStr, 10, 64)
	if err != nil || depId == 0 {
		helpers.ResponseJSON(w, map[string]string{"error": "invalid department id"}, http.StatusBadRequest)
		return
	}
	if _, err := h.rep.GetDepartmentById(uint(depId)); err != nil {
		helpers.ResponseJSON(w, map[string]string{"error": "department not found"}, http.StatusNotFound)
		return
	}

	var req struct {
		FullName string     `json:"full_name"`
		Position string     `json:"position"`
		HiredAt  *time.Time `json:"hired_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.ResponseJSON(w, map[string]string{"error": "invalid request body"}, http.StatusBadRequest)
		return
	}

	req.FullName = strings.TrimSpace(req.FullName)
	req.Position = strings.TrimSpace(req.Position)

	if req.FullName == "" || len(req.FullName) > 200 {
		helpers.ResponseJSON(w, map[string]string{"error": "full_name must be between 1 and 200 characters"}, http.StatusBadRequest)
		return
	}
	if req.Position == "" || len(req.Position) > 200 {
		helpers.ResponseJSON(w, map[string]string{"error": "position must be between 1 and 200 characters"}, http.StatusBadRequest)
		return
	}

	emp := &models.Employee{
		DepartmentId: uint(depId),
		FullName:     req.FullName,
		Position:     req.Position,
		HiredAt:      req.HiredAt,
	}
	if err := h.rep.CreateEmployee(emp); err != nil {
		h.log.Error("create employee", "err", err)
		helpers.ResponseJSON(w, map[string]string{"error": "internal server error"}, http.StatusInternalServerError)
		return
	}

	helpers.ResponseJSON(w, emp, http.StatusCreated)
}

// GET /departments/{id}?depth=1&include_employees=true
func (h *Handler) GetDepartment(w http.ResponseWriter, r *http.Request, dep *models.Department) {
	depth, _ := strconv.Atoi(r.URL.Query().Get("depth"))
	if depth <= 0 {
		depth = 1
	}
	if depth > 5 {
		depth = 5
	}

	includeParam := r.URL.Query().Get("include_employees")
	includeEmployees := includeParam == "" || includeParam == "true"

	h.rep.LoadDepartmentChildren(dep, depth)
	if includeEmployees {
		h.rep.LoadDepartmentEmployees(dep)
	}

	if dep.Children == nil {
		dep.Children = []models.Department{}
	}
	if includeEmployees && dep.Employees == nil {
		dep.Employees = []models.Employee{}
	}

	helpers.ResponseJSON(w, dep, http.StatusOK)
}

// PATCH /departments/{id}
func (h *Handler) PatchDepartment(w http.ResponseWriter, r *http.Request, dep *models.Department) {
	var rawReq map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&rawReq); err != nil {
		helpers.ResponseJSON(w, map[string]string{"error": "invalid request body"}, http.StatusBadRequest)
		return
	}

	nameChanged, parentChanged := false, false

	if nameRaw, ok := rawReq["name"]; ok {
		var name string
		if err := json.Unmarshal(nameRaw, &name); err != nil {
			helpers.ResponseJSON(w, map[string]string{"error": "invalid name"}, http.StatusBadRequest)
			return
		}
		name = strings.TrimSpace(name)
		if name == "" || len(name) > 200 {
			helpers.ResponseJSON(w, map[string]string{"error": "name must be between 1 and 200 characters"}, http.StatusBadRequest)
			return
		}
		dep.Name = name
		nameChanged = true
	}

	if parentRaw, ok := rawReq["parent_id"]; ok {
		var parentId *uint
		if err := json.Unmarshal(parentRaw, &parentId); err != nil {
			helpers.ResponseJSON(w, map[string]string{"error": "invalid parent_id"}, http.StatusBadRequest)
			return
		}
		if parentId != nil {
			if *parentId == dep.Id {
				helpers.ResponseJSON(w, map[string]string{"error": "department cannot be its own parent"}, http.StatusConflict)
				return
			}
			if _, err := h.rep.GetDepartmentById(*parentId); err != nil {
				helpers.ResponseJSON(w, map[string]string{"error": "parent department not found"}, http.StatusNotFound)
				return
			}
			safe, err := h.rep.FitForParent(*parentId, dep.Id)
			if err != nil {
				h.log.Error("fit for parent", "err", err)
				helpers.ResponseJSON(w, map[string]string{"error": "internal server error"}, http.StatusInternalServerError)
				return
			}
			if !safe {
				helpers.ResponseJSON(w, map[string]string{"error": "circular reference: cannot move department into its own subtree"}, http.StatusConflict)
				return
			}
		}
		dep.ParentId = parentId
		parentChanged = true
	}

	if nameChanged || parentChanged {
		if dup, _ := h.rep.GetDepartmentByName(dep.Name, dep.ParentId); dup != nil && dup.Id != dep.Id {
			helpers.ResponseJSON(w, map[string]string{"error": "department name already exists under this parent"}, http.StatusConflict)
			return
		}
	}

	if err := h.rep.SaveDepartment(dep); err != nil {
		h.log.Error("patch department", "err", err)
		helpers.ResponseJSON(w, map[string]string{"error": "internal server error"}, http.StatusInternalServerError)
		return
	}

	helpers.ResponseJSON(w, dep, http.StatusOK)
}

// DELETE /departments/{id}?mode={cascade|reassign}&reassign_to_department_id={id}
func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request, dep *models.Department) {
	mode := r.URL.Query().Get("mode")
	if mode != "cascade" && mode != "reassign" {
		helpers.ResponseJSON(w, map[string]string{"error": "mode must be 'cascade' or 'reassign'"}, http.StatusBadRequest)
		return
	}

	if mode == "reassign" {
		reassignIdStr := r.URL.Query().Get("reassign_to_department_id")
		reassignId, err := strconv.ParseUint(reassignIdStr, 10, 64)
		if err != nil || reassignId == 0 {
			helpers.ResponseJSON(w, map[string]string{"error": "reassign_to_department_id is required for mode=reassign"}, http.StatusBadRequest)
			return
		}
		if _, err := h.rep.GetDepartmentById(uint(reassignId)); err != nil {
			helpers.ResponseJSON(w, map[string]string{"error": "reassign target department not found"}, http.StatusNotFound)
			return
		}
		if err := h.rep.ReassignEmployeesAndDelete(dep.Id, uint(reassignId)); err != nil {
			h.log.Error("reassign and delete", "err", err)
			helpers.ResponseJSON(w, map[string]string{"error": "internal server error"}, http.StatusInternalServerError)
			return
		}
	} else {
		if err := h.rep.DeleteDepartmentById(dep.Id); err != nil {
			h.log.Error("cascade delete", "err", err)
			helpers.ResponseJSON(w, map[string]string{"error": "internal server error"}, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
