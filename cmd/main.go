// Package main es el punto de entrada de la aplicación servidor Prunus.
// Inicializa la conexión a la base de datos, ejecuta migraciones,
// configura la inyección de dependencias y levanta el servidor HTTP.
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

// main es la función principal de la aplicación. Orquesta el arranque completo
// del servidor: conexión a la base de datos, migraciones, inyección de
// dependencias por dominio y registro de rutas HTTP.
func main() {

	// Establece la conexión con la base de datos. Si falla, la aplicación
	// termina inmediatamente para evitar un estado inconsistente.
	db, err := database.Conexion()
	if err != nil {
		log.Fatal(err)
	}
	// Garantiza el cierre de la conexión al finalizar la ejecución del programa.
	defer db.Close()

	// Ejecuta las migraciones pendientes sobre el esquema de la base de datos.
	// Asegura que la estructura de tablas esté actualizada antes de iniciar el servidor.
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Iniciado migracion de la base de datos")

	// --- Inyección de dependencias por dominio ---
	// Cada dominio sigue el patrón: Store (acceso a datos) → Service (lógica de negocio) → Handler (transporte HTTP).

	// Dominio: Empresa
	empresaStore := store.NewEmpresa(db)
	empresaServices := services.NewServiceEmpresa(empresaStore)
	empresahandler := transport.NewEmpresaHandler(empresaServices)

	// Dominio: Sucursal
	sucusalStore := store.NewSucursal(db)
	sucursalServices := services.NewServiceSucursal(sucusalStore)
	sucursalHandler := transport.NewSucursalHandler(sucursalServices)

	// Dominio: Rol — gestión de roles y permisos del sistema.
	rolStore := store.NewRol(db)
	rolService := services.NewServiceRol(rolStore)
	rolHandler := transport.NewRolHandler(rolService)

	// Dominio: Usuario — administración de cuentas de usuario.
	usuarioStore := store.NewUsuario(db)
	usuarioService := services.NewServiceUsuario(usuarioStore)
	usuarioHandler := transport.NewUsuarioHandler(usuarioService)

	// Dominio: Autenticación — reutiliza el servicio de usuario para manejar
	// login, logout y validación de credenciales.
	authHandler := transport.NewAuthHandler(usuarioService)

	// Dominio: Categoría — clasificación de productos o servicios.
	categoriaStore := store.NewCategoria(db)
	categoriaService := services.NewServiceCategoria(categoriaStore)
	categoriaHandler := transport.NewCategoriaHandler(categoriaService)

	// Dominio: Cliente — gestión del catálogo de clientes.
	clienteStore := store.NewCliente(db)
	clienteService := services.NewServiceCliente(clienteStore)
	clienteHandler := transport.NewClienteHandler(clienteService)

	// Dominio: Unidad de medida — manejo de unidades para productos (kg, lt, pz, etc.).
	medidaStore := store.NewUnidad(db)
	medidaService := services.NewServiceUnidad(medidaStore)
	medidaHandler := transport.NewMedidaHandler(medidaService)

	// Dominio: Moneda — soporte para múltiples divisas en transacciones.
	monedaStore := store.NewMoneda(db)
	monedaService := services.NewServiceMoneda(monedaStore)
	monedaHandler := transport.NewMonedaHandler(monedaService)

	// Dominio: Producto — catálogo principal de productos del sistema.
	productoStore := store.NewProducto(db)
	productoService := services.NewServiceProducto(productoStore)
	productoHandler := transport.NewProductoHandler(productoService)

	// Dominio: Proveedor — gestión del catálogo de proveedores.
	proveedorStore := store.NewProveedor(db)
	proveedorService := services.NewServiceProveedor(proveedorStore)
	proveedorHandler := transport.NewProveedorHandler(proveedorService)

	// Registra todas las rutas HTTP consolidando los handlers de cada dominio
	// en un único router principal. Centraliza la configuración de endpoints.
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
