package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/prunus/pkg/config"
	"github.com/prunus/pkg/config/database"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/store"
	"github.com/spf13/cobra"
)

// registerCmd representa el comando base para registro de entidades
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registra entidades en el sistema (estatus, empresa, sucursal, rol, usuario)",
}

// Variables para flags
var (
	// Estatus flags
	estDescription string
	estType        string
	estModuloID    int

	// Empresa flags
	empNombre   string
	empRut      string
	empIDStatus string

	// Sucursal flags
	sucIDEmpresa string
	sucNombre    string
	sucIDStatus  string

	// Rol flags
	rolIDSucursal string
	rolNombre     string
	rolIDStatus   string

	// Usuario flags
	usuIDSucursal string
	usuIDRol      string
	usuEmail      string
	usuNombre     string
	usuDni        string
	usuPassword   string
	usuIDStatus   string
)

func init() {
	rootCmd.AddCommand(registerCmd)

	// Subcomando Estatus
	registerEstatusCmd := &cobra.Command{
		Use:   "estatus",
		Short: "Registra un nuevo estatus",
		Run:   runRegisterEstatus,
	}
	registerEstatusCmd.Flags().StringVar(&estDescription, "desc", "", "Descripción del estatus")
	registerEstatusCmd.Flags().StringVar(&estType, "tipo", "", "Tipo de estado")
	registerEstatusCmd.Flags().IntVar(&estModuloID, "modulo", 0, "ID del módulo (1:Empresa, 2:Sucursal, 3:Usuario, etc.)")
	registerEstatusCmd.MarkFlagRequired("desc")
	registerEstatusCmd.MarkFlagRequired("tipo")
	registerEstatusCmd.MarkFlagRequired("modulo")
	registerCmd.AddCommand(registerEstatusCmd)

	// Subcomando Empresa
	registerEmpresaCmd := &cobra.Command{
		Use:   "empresa",
		Short: "Registra una nueva empresa",
		Run:   runRegisterEmpresa,
	}
	registerEmpresaCmd.Flags().StringVar(&empNombre, "nombre", "", "Nombre de la empresa")
	registerEmpresaCmd.Flags().StringVar(&empRut, "rut", "", "RUT de la empresa")
	registerEmpresaCmd.Flags().StringVar(&empIDStatus, "status", "", "ID del estatus (UUID)")
	registerEmpresaCmd.MarkFlagRequired("nombre")
	registerEmpresaCmd.MarkFlagRequired("rut")
	registerEmpresaCmd.MarkFlagRequired("status")
	registerCmd.AddCommand(registerEmpresaCmd)

	// Subcomando Sucursal
	registerSucursalCmd := &cobra.Command{
		Use:   "sucursal",
		Short: "Registra una nueva sucursal",
		Run:   runRegisterSucursal,
	}
	registerSucursalCmd.Flags().StringVar(&sucIDEmpresa, "empresa", "", "ID de la empresa (UUID)")
	registerSucursalCmd.Flags().StringVar(&sucNombre, "nombre", "", "Nombre de la sucursal")
	registerSucursalCmd.Flags().StringVar(&sucIDStatus, "status", "", "ID del estatus (UUID)")
	registerSucursalCmd.MarkFlagRequired("empresa")
	registerSucursalCmd.MarkFlagRequired("nombre")
	registerSucursalCmd.MarkFlagRequired("status")
	registerCmd.AddCommand(registerSucursalCmd)

	// Subcomando Rol
	registerRolCmd := &cobra.Command{
		Use:   "rol",
		Short: "Registra un nuevo rol",
		Run:   runRegisterRol,
	}
	registerRolCmd.Flags().StringVar(&rolIDSucursal, "sucursal", "", "ID de la sucursal (UUID)")
	registerRolCmd.Flags().StringVar(&rolNombre, "nombre", "", "Nombre del rol")
	registerRolCmd.Flags().StringVar(&rolIDStatus, "status", "", "ID del estatus (UUID)")
	registerRolCmd.MarkFlagRequired("sucursal")
	registerRolCmd.MarkFlagRequired("nombre")
	registerRolCmd.MarkFlagRequired("status")
	registerCmd.AddCommand(registerRolCmd)

	// Subcomando Usuario
	registerUsuarioCmd := &cobra.Command{
		Use:   "usuario",
		Short: "Registra un nuevo usuario",
		Run:   runRegisterUsuario,
	}
	registerUsuarioCmd.Flags().StringVar(&usuIDSucursal, "sucursal", "", "ID de la sucursal (UUID)")
	registerUsuarioCmd.Flags().StringVar(&usuIDRol, "rol", "", "ID del rol (UUID)")
	registerUsuarioCmd.Flags().StringVar(&usuEmail, "email", "", "Email del usuario")
	registerUsuarioCmd.Flags().StringVar(&usuNombre, "nombre", "", "Nombre completo")
	registerUsuarioCmd.Flags().StringVar(&usuDni, "dni", "", "DNI/Documento")
	registerUsuarioCmd.Flags().StringVar(&usuPassword, "password", "", "Contraseña")
	registerUsuarioCmd.Flags().StringVar(&usuIDStatus, "status", "", "ID del estatus (UUID)")
	registerUsuarioCmd.MarkFlagRequired("sucursal")
	registerUsuarioCmd.MarkFlagRequired("rol")
	registerUsuarioCmd.MarkFlagRequired("email")
	registerUsuarioCmd.MarkFlagRequired("nombre")
	registerUsuarioCmd.MarkFlagRequired("dni")
	registerUsuarioCmd.MarkFlagRequired("password")
	registerUsuarioCmd.MarkFlagRequired("status")
	registerCmd.AddCommand(registerUsuarioCmd)
}

// Funciones auxiliares para inicializar dependencias
func initCLI() (*sql.DB, models.CacheStore, *slog.Logger) {
	if err := config.Validate("DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME"); err != nil {
		log.Fatalf("❌ Error de configuración: %v", err)
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("❌ Error conectando a la base de datos: %v", err)
	}

	// No inicializamos Redis para operaciones rápidas de CLI a menos que sea necesario
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return db, nil, logger
}

func runRegisterEstatus(cmd *cobra.Command, args []string) {
	db, cache, logger := initCLI()
	defer db.Close()

	estatusStore := store.NewEstatus(db)
	estatusService := services.NewServiceEstatus(estatusStore, cache, logger)

	model := models.Estatus{
		StdDescripcion: estDescription,
		StdTipoEstado:  estType,
		MdlID:          estModuloID,
	}

	result, err := estatusService.CreateEstatus(context.Background(), model)
	if err != nil {
		fmt.Printf("❌ Error registrando estatus: %v\n", err)
		return
	}
	fmt.Printf("✅ Estatus registrado con ID: %s\n", result.IDStatus)
}

func runRegisterEmpresa(cmd *cobra.Command, args []string) {
	db, _, logger := initCLI()
	defer db.Close()

	statusID, err := uuid.Parse(empIDStatus)
	if err != nil {
		fmt.Printf("❌ ID de status inválido: %v\n", err)
		return
	}

	empresaStore := store.NewEmpresa(db)
	empresaService := services.NewServiceEmpresa(empresaStore, logger)

	model := models.Empresa{
		Nombre:   empNombre,
		RUT:      empRut,
		IDStatus: statusID,
	}

	result, err := empresaService.CrearEmpresa(context.Background(), model)
	if err != nil {
		fmt.Printf("❌ Error registrando empresa: %v\n", err)
		return
	}
	fmt.Printf("✅ Empresa registrada con ID: %s\n", result.IDEmpresa)
}

func runRegisterSucursal(cmd *cobra.Command, args []string) {
	db, _, logger := initCLI()
	defer db.Close()

	empID, err := uuid.Parse(sucIDEmpresa)
	if err != nil {
		fmt.Printf("❌ ID de empresa inválido: %v\n", err)
		return
	}

	statusID, err := uuid.Parse(sucIDStatus)
	if err != nil {
		fmt.Printf("❌ ID de status inválido: %v\n", err)
		return
	}

	sucursalStore := store.NewSucursal(db)
	sucursalService := services.NewServiceSucursal(sucursalStore, logger)

	model := models.Sucursal{
		IDEmpresa:      empID,
		NombreSucursal: sucNombre,
		IDStatus:       statusID,
	}

	result, err := sucursalService.CreateSucursal(context.Background(), model)
	if err != nil {
		fmt.Printf("❌ Error registrando sucursal: %v\n", err)
		return
	}
	fmt.Printf("✅ Sucursal registrada con ID: %s\n", result.IDSucursal)
}

func runRegisterRol(cmd *cobra.Command, args []string) {
	db, cache, logger := initCLI()
	defer db.Close()

	sucID, err := uuid.Parse(rolIDSucursal)
	if err != nil {
		fmt.Printf("❌ ID de sucursal inválido: %v\n", err)
		return
	}

	statusID, err := uuid.Parse(rolIDStatus)
	if err != nil {
		fmt.Printf("❌ ID de status inválido: %v\n", err)
		return
	}

	rolStore := store.NewRol(db)
	rolService := services.NewServiceRol(rolStore, cache, logger)

	model := models.Rol{
		IDSucursal: sucID,
		RolNombre:  rolNombre,
		IDStatus:   statusID,
	}

	result, err := rolService.CreateRol(context.Background(), model)
	if err != nil {
		fmt.Printf("❌ Error registrando rol: %v\n", err)
		return
	}
	fmt.Printf("✅ Rol registrado con ID: %s\n", result.IDRol)
}

func runRegisterUsuario(cmd *cobra.Command, args []string) {
	db, _, logger := initCLI()
	defer db.Close()

	sucID, err := uuid.Parse(usuIDSucursal)
	if err != nil {
		fmt.Printf("❌ ID de sucursal inválido: %v\n", err)
		return
	}

	rolID, err := uuid.Parse(usuIDRol)
	if err != nil {
		fmt.Printf("❌ ID de rol inválido: %v\n", err)
		return
	}

	statusID, err := uuid.Parse(usuIDStatus)
	if err != nil {
		fmt.Printf("❌ ID de status inválido: %v\n", err)
		return
	}

	usuarioStore := store.NewUsuario(db)
	logsStore := store.NewLogs(db)
	usuarioService := services.NewServiceUsuario(usuarioStore, logsStore, logger)

	model := models.Usuario{
		IDSucursal: sucID,
		IDRol:      rolID,
		Email:      usuEmail,
		UsuNombre:  usuNombre,
		UsuDNI:     usuDni,
		Password:   usuPassword,
		IDStatus:   statusID,
	}

	result, err := usuarioService.CreateUsuario(context.Background(), model)
	if err != nil {
		fmt.Printf("❌ Error registrando usuario: %v\n", err)
		return
	}
	fmt.Printf("✅ Usuario registrado con ID: %s\n", result.IDUsuario)
}
