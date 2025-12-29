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
	q := `CREATE TABLE IF NOT EXISTS usuario (
    id_usuario    SERIAL PRIMARY KEY,
    id_sucursal   INTEGER NOT NULL,
    id_rol        INTEGER NOT NULL,

    email         VARCHAR(150) NOT NULL UNIQUE,
    usu_nombre    VARCHAR(150) NOT NULL,
    usu_dni       VARCHAR(30)  NOT NULL UNIQUE,
    usu_telefono  VARCHAR(30),
    password      TEXT NOT NULL,
    estado        INTEGER NOT NULL DEFAULT 1,

    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP NULL,

    CONSTRAINT fk_usuario_sucursal
        FOREIGN KEY (id_sucursal)
        REFERENCES sucursal(id_sucursal)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_usuario_rol
        FOREIGN KEY (id_rol)
        REFERENCES rol(id_rol)
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

	// inyetar depedencia sucursal
	sucusalStore := store.NewSucursal(db)
	sucursalServices := services.NewServiceSucursal(sucusalStore)
	sucursalHandler := transport.NewSucursalHandler(sucursalServices)
	// inyetar depedencia rol
	rolStore := store.NewRol(db)
	rolService := services.NewServiceRol(rolStore)
	rolHandler := transport.NewRolHandler(rolService)
	// inyectar dependencia usuario
	usuarioStore := store.NewUsuario(db)
	usuarioService := services.NewServiceUsuario(usuarioStore)
	usuarioHandler := transport.NewUsuarioHandler(usuarioService)

	// configura rutas - combinar todos los handlers en un solo router
	router := routers.NewMainRouter(empresahandler, sucursalHandler, rolHandler, usuarioHandler)

	// empezar y escuchar el servidor
	fmt.Println("✅ Iniciando servidor")
	http.ListenAndServe(":9090", router)
}
