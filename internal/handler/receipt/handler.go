package receipt

import (
	"net/http"
	"strconv"

	"student_service_app/backend/internal/handler"
	"student_service_app/backend/internal/mapper"
	"student_service_app/backend/internal/response"
	servicecommon "student_service_app/backend/internal/service"
	servicepkg "student_service_app/backend/internal/service/receipt"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service servicepkg.Service
}

func NewHandler(service servicepkg.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	result := make([]any, 0, len(items))
	for _, item := range items {
		result = append(result, mapper.ReceiptToDTO(item))
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) GetByKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if id, err := strconv.ParseInt(key, 10, 64); err == nil {
		item, err := h.service.GetByID(r.Context(), id)
		if err != nil {
			response.Error(w, err)
			return
		}
		response.JSON(w, http.StatusOK, mapper.ReceiptToDTO(*item))
		return
	}

	item, err := h.service.GetByReceiptNo(r.Context(), key)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, mapper.ReceiptToDTO(*item))
}

func ParseID(r *http.Request, param string) (int64, error) {
	return handler.ParseID(r, param)
}
