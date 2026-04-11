package models

import (
	"fmt"
	"strings"
	"time"
)

// JSONDate es un tipo personalizado para manejar fechas en JSON con múltiples formatos
type JSONDate time.Time

// Formatos soportados
var dateFormats = []string{
	"2006-01-02T15:04:05Z07:00", // RFC3339
	"2006-01-02",                // YYYY-MM-DD
}

// UnmarshalJSON implementa la interfaz json.Unmarshaler
func (j *JSONDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}

	var lastErr error
	for _, format := range dateFormats {
		t, err := time.Parse(format, s)
		if err == nil {
			*j = JSONDate(t)
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("formato de fecha inválido: %s. Formatos esperados: YYYY-MM-DD o RFC3339. Error: %v", s, lastErr)
}

// MarshalJSON implementa la interfaz json.Marshaler
func (j JSONDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))), nil
}

// ToTime convierte JSONDate a time.Time
func (j JSONDate) ToTime() time.Time {
	return time.Time(j)
}

// IsZero comprueba si la fecha es el valor cero
func (j JSONDate) IsZero() bool {
	return time.Time(j).IsZero()
}
