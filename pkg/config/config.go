package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Get retorna el valor de una variable de configuración como string.
// Busca primero en el archivo .env cargado por viper, luego en variables de entorno del sistema.
func Get(key string) string {
	return viper.GetString(key)
}

// GetInt retorna el valor de una variable de configuración como int.
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetDefault retorna el valor de una variable de configuración o un default si está vacía.
func GetDefault(key, defaultValue string) string {
	if v := viper.GetString(key); v != "" {
		return v
	}
	return defaultValue
}

// Validate verifica que todas las claves requeridas tengan valor.
// Retorna un error descriptivo listando todas las claves faltantes.
func Validate(requiredKeys ...string) error {
	var missing []string
	for _, key := range requiredKeys {
		if viper.GetString(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("variables de configuración requeridas no encontradas: [%s]", strings.Join(missing, ", "))
	}
	return nil
}
