package request

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/utils/validator"
)

var (
	ErrInvalidJSON = errors.New("JSON inválido")
	ErrInvalidID   = errors.New("ID inválido")
)

// DecodeJSON decodifica el cuerpo de la solicitud en el destino proporcionado
func DecodeJSON(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return ErrInvalidJSON
	}
	return nil
}

// DecodeAndValidate decodifica y valida el cuerpo de la solicitud
func DecodeAndValidate(r *http.Request, dst interface{}) error {
	if err := DecodeJSON(r, dst); err != nil {
		return err
	}
	return validator.Validate.Struct(dst)
}

// GetID recupera un UUID de los parámetros de la URL
func GetID(r *http.Request, param string) (uuid.UUID, error) {
	idStr := chi.URLParam(r, param)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, ErrInvalidID
	}
	return id, nil
}
