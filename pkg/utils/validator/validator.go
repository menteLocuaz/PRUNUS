package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate es la instancia global del validador
var Validate *validator.Validate

func init() {
	Validate = validator.New()

	// Registrar función para obtener el nombre de la etiqueta json en los errores
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// FormatErrors formatea los errores de validación en un mapa amigable
func FormatErrors(err error) map[string]string {
	errors := make(map[string]string)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			errors[e.Field()] = getErrorMsg(e)
		}
	}
	return errors
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Este campo es obligatorio"
	case "email":
		return "Formato de email inválido"
	case "min":
		return fmt.Sprintf("Debe tener al menos %s caracteres", fe.Param())
	case "max":
		return fmt.Sprintf("No puede tener más de %s caracteres", fe.Param())
	case "oneof":
		return fmt.Sprintf("Debe ser uno de los siguientes valores: %s", fe.Param())
	}
	return "Valor inválido"
}
