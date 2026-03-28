package main

import (
	"database/sql"
	"log/slog"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/routers"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/store"
	transport "github.com/prunus/pkg/transport/http"
)

// RegisterHandlers centraliza la inyección de dependencias y registro de handlers.
func RegisterHandlers(db *sql.DB, cacheStore models.CacheStore, logger *slog.Logger) *routers.Handlers {
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
	inventarioStore := store.NewInventario(db)
	agregadoresStore := store.NewAgregadores(db)
	cajaStore := store.NewCaja(db)
	facturaStore := store.NewFactura(db)
	ordenPedidoStore := store.NewOrdenPedido(db)
	dispositivoPosStore := store.NewDispositivoPosStore(db)
	estacionPosStore := store.NewEstacionPosStore(db)
	compraStore := store.NewCompra(db)
	periodoStore := store.NewPeriodoStore(db)
	configuracionStore := store.NewConfiguracion(db)
	logsStore := store.NewLogs(db)

	// 2. Services (Lógica de Negocio)
	empresaServices := services.NewServiceEmpresa(empresaStore, logger)
	sucursalServices := services.NewServiceSucursal(sucursalStore, logger)
	rolService := services.NewServiceRol(rolStore, cacheStore, logger)
	usuarioService := services.NewServiceUsuario(usuarioStore, logsStore, logger)
	categoriaService := services.NewServiceCategoria(categoriaStore, cacheStore, logger)
	clienteService := services.NewServiceCliente(clienteStore, logger)
	medidaService := services.NewServiceUnidad(medidaStore, cacheStore, logger)
	monedaService := services.NewServiceMoneda(monedaStore, cacheStore, logger)
	productoService := services.NewServiceProducto(productoStore, inventarioStore, logger)
	proveedorService := services.NewServiceProveedor(proveedorStore, logger)
	estatusService := services.NewServiceEstatus(estatusStore, cacheStore, logger)
	posService := services.NewServicePOS(posStore, logsStore, logger)
	inventarioService := services.NewServiceInventario(inventarioStore, logger)
	agregadoresService := services.NewServiceAgregadores(agregadoresStore, logger)
	cajaService := services.NewServiceCaja(cajaStore, logger)
	facturaService := services.NewServiceFactura(facturaStore, cacheStore, logger)
	ordenPedidoService := services.NewServiceOrdenPedido(ordenPedidoStore, logger)
	dispositivoPosService := services.NewServiceDispositivoPos(dispositivoPosStore, logger)
	estacionPosService := services.NewServiceEstacionPos(estacionPosStore, logger)
	compraService := services.NewServiceCompra(compraStore, inventarioService)
	periodoService := services.NewServicePeriodo(periodoStore, posStore, logger)
	configuracionService := services.NewServiceConfiguracion(configuracionStore)

	// 3. Handlers (Controladores)
	return &routers.Handlers{
		Empresa:        transport.NewEmpresaHandler(empresaServices),
		Sucursal:       transport.NewSucursalHandler(sucursalServices),
		Rol:            transport.NewRolHandler(rolService),
		Usuario:        transport.NewUsuarioHandler(usuarioService),
		Auth:           transport.NewAuthHandler(usuarioService),
		Categoria:      transport.NewCategoriaHandler(categoriaService),
		Cliente:        transport.NewClienteHandler(clienteService),
		Medida:         transport.NewMedidaHandler(medidaService),
		Moneda:         transport.NewMonedaHandler(monedaService),
		Producto:       transport.NewProductoHandler(productoService),
		Proveedor:      transport.NewProveedorHandler(proveedorService),
		Estatus:        transport.NewEstatusHandler(estatusService),
		POS:            transport.NewPOSHandler(posService),
		Inventario:     transport.NewInventarioHandler(inventarioService),
		Agregadores:    transport.NewAgregadoresHandler(agregadoresService),
		Caja:           transport.NewCajaHandler(cajaService),
		Factura:        transport.NewFacturaHandler(facturaService),
		OrdenPedido:    transport.NewOrdenPedidoHandler(ordenPedidoService),
		DispositivoPos: transport.NewDispositivoPosHandler(dispositivoPosService),
		EstacionPos:    transport.NewEstacionPosHandler(estacionPosService),
		Compra:         transport.NewCompraHandler(compraService),
		Periodo:        transport.NewPeriodoHandler(periodoService),
		Configuracion:  transport.NewConfiguracionHandler(configuracionService),
	}
}
