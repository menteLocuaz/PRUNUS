package main

import (
	"fmt"
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
	q := `CREATE TABLE IF NOT EXISTS sucursal (
    id_sucursal     SERIAL PRIMARY KEY,
    id_empresa      INTEGER NOT NULL,
    nombre_sucursal VARCHAR(255) NOT NULL,
    estado           INTEGER NOT NULL DEFAULT 1,

    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMP NULL,

    CONSTRAINT fk_sucursal_empresa
        FOREIGN KEY (id_empresa)
        REFERENCES empresa(id_empresa)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);
`
	if _, err := db.Exec(q); err != nil {

		log.Fatal(err.Error())
	}

	// inyetar depedencia emmpresa
	empresaStore := store.NewEmpresa(db)
	empresaServices := services.NewServiceEmpresa(empresaStore)
	empresahandler := transport.NewEmpresaHandler(empresaServices)

	// inyetar depedencia emmpresa
	sucusalStore := store.NewSucursal(db)
	sucursalServices := services.NewServiceSucursal(sucusalStore)
	sucursalHandler := transport.NewSucursalHandler(sucursalServices)

	// configura rutas
	router := routers.NewRouter(empresahandler)
	router = routers.NewRouterSucursal(sucursalHandler)
	// empeazar y escuchar el servidor
	fmt.Println("✅ Iniciando servidor")
	http.ListenAndServe(":9090", router)
}
