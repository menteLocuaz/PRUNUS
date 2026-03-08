package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prunus/pkg/config/database"
	"github.com/prunus/pkg/config/database/migrations"
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

	// redis
	rdb, err := database.RedisConexion()
	if err != nil {
		log.Printf("Aviso: No se pudo conectar a Redis: %v. Cache desactivado.", err)
	}
	cacheStore := store.NewRedisStore(rdb)

	// ejecutar migraciones
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Iniciado migracion de la base de datos")

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
	rolService := services.NewServiceRol(rolStore, cacheStore)
	rolHandler := transport.NewRolHandler(rolService)
	// inyectar dependencia usuario
	usuarioStore := store.NewUsuario(db)
	usuarioService := services.NewServiceUsuario(usuarioStore)
	usuarioHandler := transport.NewUsuarioHandler(usuarioService)

	// inyectar dependencia autenticación (usa el mismo servicio de usuario)
	authHandler := transport.NewAuthHandler(usuarioService)

	// inyectar dependencia categoria
	categoriaStore := store.NewCategoria(db)
	categoriaService := services.NewServiceCategoria(categoriaStore, cacheStore)
	categoriaHandler := transport.NewCategoriaHandler(categoriaService)

	// inyectar dependencia cliente
	clienteStore := store.NewCliente(db)
	clienteService := services.NewServiceCliente(clienteStore)
	clienteHandler := transport.NewClienteHandler(clienteService)

	// inyectar dependencia medida
	medidaStore := store.NewUnidad(db)
	medidaService := services.NewServiceUnidad(medidaStore)
	medidaHandler := transport.NewMedidaHandler(medidaService)

	// inyectar dependencia moneda
	monedaStore := store.NewMoneda(db)
	monedaService := services.NewServiceMoneda(monedaStore)
	monedaHandler := transport.NewMonedaHandler(monedaService)

	// inyectar dependencia producto
	productoStore := store.NewProducto(db)
	productoService := services.NewServiceProducto(productoStore)
	productoHandler := transport.NewProductoHandler(productoService)

	// inyectar dependencia proveedor
	proveedorStore := store.NewProveedor(db)
	proveedorService := services.NewServiceProveedor(proveedorStore)
	proveedorHandler := transport.NewProveedorHandler(proveedorService)

	// configura rutas - combinar todos los handlers en un solo router
	router := routers.NewMainRouter(
		empresahandler,
		sucursalHandler,
		rolHandler,
		usuarioHandler,
		authHandler,
		categoriaHandler,
		clienteHandler,
		medidaHandler,
		monedaHandler,
		productoHandler,
		proveedorHandler,
	)

	// Inicia el servidor HTTP en el puerto 9090 y queda en escucha de peticiones entrantes.
	// Si el servidor no puede iniciar, log.Fatal detendrá la aplicación con el error correspondiente.
	fmt.Println("✅ Iniciando servidor en :9090")
	if err := http.ListenAndServe(":9090", router); err != nil {
		log.Fatal(err)
	}
}
