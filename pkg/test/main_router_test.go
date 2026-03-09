package routers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prunus/pkg/routers"
)

func TestMainRouter_PluralRoutes(t *testing.T) {
	// Setup - Usamos nil para los handlers ya que solo queremos probar que la ruta existe
	// (la respuesta será 401 Unauthorized por el middleware o pánico si llega al handler nil,
	// pero lo que nos importa es si el router reconoce la ruta pluralizada).

	router := routers.NewMainRouter(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/usuarios"},
		{"GET", "/api/v1/sucursales"},
		{"GET", "/api/v1/roles"},
		{"GET", "/api/v1/productos"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			// Esperamos 401 porque no estamos autenticados, pero 404 significaría que la ruta no existe
			if rr.Code == http.StatusNotFound {
				t.Errorf("Ruta %s no encontrada (404), se esperaba pluralización correcta", tt.path)
			}
		})
	}
}
