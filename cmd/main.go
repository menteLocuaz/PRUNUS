package main

import (
	"log"
	"net/http"

	"github.com/prunus/pkg/config/database"
	"github.com/prunus/pkg/routers"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/store"
	transport "github.com/prunus/pkg/transport/http"
)

func main() {

	// base de datos
	db, err := database.Conexion()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// crear la tabla si no existe
	q := `CREATE TABLE IF NOT EXISTS empresa (
    id_empresa SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    rut VARCHAR(20) NOT NULL UNIQUE,
    estado INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
`
	if _, err := db.Exec(q); err != nil {

		log.Fatal(err.Error())
	}

	// inyetar depebdicia
	empresaStore := store.NewEmpresa(db)
	empresaServices := services.NewServiceEmpresa(empresaStore)
	empresahandler := transport.NewEmpresaHandler(empresaServices)
	// configura rutas
	router := routers.NewRouter(empresahandler)

	// empeazar y escuchar el servidor

	http.ListenAndServe(":9090", router)
}
