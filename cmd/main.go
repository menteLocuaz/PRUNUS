package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/prunus/pkg/config/database"
	"github.com/prunus/pkg/config/database/migrations"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/routers"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/store"
	transport "github.com/prunus/pkg/transport/http"
)

func main() {
	// Base de datos
	db, err := database.Conexion()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Redis
	rdb, err := database.RedisConexion()
	if err != nil {
		log.Printf("Aviso: No se pudo conectar a Redis: %v. Cache desactivado.", err)
	}
	cacheStore := store.NewRedisStore(rdb)

	// Ejecutar migraciones
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Iniciando migracion de la base de datos")

	// Registrar todos los handlers
	h := RegisterHandlers(db, cacheStore)

	// Configurar router principal
	router := routers.NewMainRouter(h)

	// Iniciar servidor
	fmt.Println("✅ Iniciando servidor en :9090")
	if err := http.ListenAndServe(":9090", router); err != nil {
		log.Fatal(err)
	}
}

// RegisterHandlers centraliza la inyección de dependencias y registro de handlers
func RegisterHandlers(db *sql.DB, cacheStore models.CacheStore) *routers.Handlers {
	// 1. Stores (Repositorios)
	empresaStore := store.NewEmpresa(db)
	sucursalStore := store.NewSucursal(db)
	rolStore := store.NewRol(db)
	usuarioStore := store.NewUsuario(db)
	categoriaStore := store.NewCategoria(db)
	clienteStore := store.NewCliente(db)
	medidaStore := store.NewUnidad(db)
	monedaStore := store.NewMoneda(db)
	productoStore := store.NewProducto(db)
	proveedorStore := store.NewProveedor(db)
	estatusStore := store.NewEstatus(db)
	posStore := store.NewPOSStore(db)

	// 2. Services (Lógica de Negocio)
	empresaServices := services.NewServiceEmpresa(empresaStore)
	sucursalServices := services.NewServiceSucursal(sucursalStore)
	rolService := services.NewServiceRol(rolStore, cacheStore)
	usuarioService := services.NewServiceUsuario(usuarioStore)
	categoriaService := services.NewServiceCategoria(categoriaStore, cacheStore)
	clienteService := services.NewServiceCliente(clienteStore)
	medidaService := services.NewServiceUnidad(medidaStore)
	monedaService := services.NewServiceMoneda(monedaStore)
	productoService := services.NewServiceProducto(productoStore)
	proveedorService := services.NewServiceProveedor(proveedorStore)
	estatusService := services.NewServiceEstatus(estatusStore, cacheStore)
	posService := services.NewServicePOS(posStore)

	// 3. Handlers (Controladores)
	return &routers.Handlers{
		Empresa:   transport.NewEmpresaHandler(empresaServices),
		Sucursal:  transport.NewSucursalHandler(sucursalServices),
		Rol:       transport.NewRolHandler(rolService),
		Usuario:   transport.NewUsuarioHandler(usuarioService),
		Auth:      transport.NewAuthHandler(usuarioService),
		Categoria: transport.NewCategoriaHandler(categoriaService),
		Cliente:   transport.NewClienteHandler(clienteService),
		Medida:    transport.NewMedidaHandler(medidaService),
		Moneda:    transport.NewMonedaHandler(monedaService),
		Producto:  transport.NewProductoHandler(productoService),
		Proveedor: transport.NewProveedorHandler(proveedorService),
		Estatus:   transport.NewEstatusHandler(estatusService),
		POS:       transport.NewPOSHandler(posService),
	}
}
