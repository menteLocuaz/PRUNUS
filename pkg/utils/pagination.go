package utils

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prunus/pkg/dto"
)

// ParsePaginationParams extrae los parámetros de paginación de una solicitud HTTP
func ParsePaginationParams(r *http.Request) dto.PaginationParams {
	params := dto.PaginationParams{
		Limit: dto.DefaultLimit,
	}

	query := r.URL.Query()

	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			if limit > dto.MaxLimit {
				limit = dto.MaxLimit
			}
			params.Limit = limit
		}
	}

	params.LastID = query.Get("last_id")

	if lastDateStr := query.Get("last_date"); lastDateStr != "" {
		if t, err := time.Parse(time.RFC3339, lastDateStr); err == nil {
			params.LastDate = &t
		}
	}

	return params
}
