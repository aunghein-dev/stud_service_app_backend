package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"student_service_app/backend/internal/errs"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Base struct {
	validate *validator.Validate
}

func NewBase(validate *validator.Validate) Base {
	return Base{validate: validate}
}

func (b Base) Decode(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errs.BadRequest("invalid request body")
	}
	if err := b.validate.Struct(v); err != nil {
		return errs.BadRequest(err.Error())
	}
	return nil
}

func ParseID(r *http.Request, param string) (int64, error) {
	idStr := chi.URLParam(r, param)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return 0, errs.BadRequest(fmt.Sprintf("invalid %s", param))
	}
	return id, nil
}
