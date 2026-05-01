package main

import (
	"database/sql"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/routers"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/store"
	transport "github.com/prunus/pkg/transport/http"
	"github.com/prunus/pkg/utils"
	"go.uber.org/zap"
)

// RegisterHandlers centraliza la inyección de dependencias y registro de handlers.
func RegisterHandlers(db *sql.DB, cacheStore models.CacheStore, logger *zap.Logger) *routers.Handlers {
	// 0. Cache Manager (Capa centralizada de caché)
	cacheMgr := utils.NewCacheManager(cacheStore)

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
	cajaStore := store.NewCajaStore(db)
	facturaStore := store.NewFactura(db)
	ordenPedidoStore := store.NewOrdenPedido(db)
	dispositivoPosStore := store.NewDispositivoPosStore(db)
	estacionPosStore := store.NewEstacionPosStore(db)
	compraStore := store.NewCompra(db)
	periodoStore := store.NewPeriodoStore(db)
	configuracionStore := store.NewConfiguracion(db)
	logsStore := store.NewLogs(db)
	dashboardStore := store.NewDashboardStore(db)

	// 2. Services (Lógica de Negocio)
	empresaServices := services.NewServiceEmpresa(empresaStore, logger)
	sucursalServices := services.NewServiceSucursal(sucursalStore, logger)
	rolService := services.NewServiceRol(rolStore, cacheMgr, logger)
	usuarioService := services.NewServiceUsuario(usuarioStore, rolService, logsStore, logger)
	categoriaService := services.NewServiceCategoria(categoriaStore, cacheMgr, logger)
	clienteService := services.NewServiceCliente(clienteStore, logger)
	medidaService := services.NewServiceUnidad(medidaStore, cacheMgr, logger)
	monedaService := services.NewServiceMoneda(monedaStore, cacheMgr, logger)
	productoService := services.NewServiceProducto(productoStore, inventarioStore, cacheMgr, logger)
	proveedorService := services.NewServiceProveedor(proveedorStore, logger)
	estatusService := services.NewServiceEstatus(estatusStore, cacheMgr, logger)
	posService := services.NewServicePOS(posStore, usuarioStore, logsStore, logger)
	inventarioService := services.NewServiceInventario(inventarioStore, logger)
	agregadoresService := services.NewServiceAgregadores(agregadoresStore, logger)
	cajaService := services.NewServiceCaja(cajaStore, usuarioStore, logger)
	facturaService := services.NewServiceFactura(facturaStore, cacheMgr, logger)
	ordenPedidoService := services.NewServiceOrdenPedido(ordenPedidoStore, logger)
	dispositivoPosService := services.NewServiceDispositivoPos(dispositivoPosStore, logger)
	estacionPosService := services.NewServiceEstacionPos(estacionPosStore, logger)
	compraService := services.NewServiceCompra(compraStore, inventarioService)
	periodoService := services.NewServicePeriodo(periodoStore, posStore, logger)
	configuracionService := services.NewServiceConfiguracion(configuracionStore)
	dashboardService := services.NewDashboardService(dashboardStore)

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
		Factura:        transport.NewFacturaHandler(facturaService, logger),
		OrdenPedido:    transport.NewOrdenPedidoHandler(ordenPedidoService),
		DispositivoPos: transport.NewDispositivoPosHandler(dispositivoPosService),
		EstacionPos:    transport.NewEstacionPosHandler(estacionPosService),
		Compra:         transport.NewCompraHandler(compraService),
		Periodo:        transport.NewPeriodoHandler(periodoService),
		Configuracion:  transport.NewConfiguracionHandler(configuracionService),
		Dashboard:      transport.NewDashboardHandler(dashboardService),
	}
}
