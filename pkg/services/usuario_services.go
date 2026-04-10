package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/prunus/pkg/helper"
	"github.com/prunus/pkg/middleware"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ServiceUsuario servicio que encapsula la lógica de negocio para usuario
type ServiceUsuario struct {
	store     store.StoreUsuario
	logsStore store.StoreLogs
	logger    *slog.Logger
}

// NewServiceUsuario crea una nueva instancia del servicio de usuario
func NewServiceUsuario(s store.StoreUsuario, l store.StoreLogs, logger *slog.Logger) *ServiceUsuario {
	return &ServiceUsuario{
		store:     s,
		logsStore: l,
		logger:    logger,
	}
}

func (s *ServiceUsuario) validateUser(u *models.Usuario, isUpdate bool) error {
	if u.Email == "" {
		return errors.New("el email del usuario es requerido")
	}
	if !emailRegex.MatchString(u.Email) {
		return errors.New("el formato del email es inválido")
	}
	if u.UsuNombre == "" {
		return errors.New("el nombre del usuario es requerido")
	}
	if u.UsuDNI == "" {
		return errors.New("el DNI del usuario es requerido")
	}
	if !isUpdate && u.Password == "" {
		return errors.New("la contraseña del usuario es requerida")
	}
	if u.IDSucursal == uuid.Nil {
		return errors.New("el ID de la sucursal es requerido")
	}
	if !isUpdate && u.IDRol == uuid.Nil {
		return errors.New("el ID del rol es requerido")
	}
	return nil
}

// GetAllUsuarios obtiene todos los usuarios del sistema
func (s *ServiceUsuario) GetAllUsuarios(ctx context.Context) ([]*models.Usuario, error) {
	return s.store.GetAllUsuarios(ctx)
}

// GetUsuarioByID obtiene un usuario por su ID con sus permisos cargados
func (s *ServiceUsuario) GetUsuarioByID(ctx context.Context, id uuid.UUID) (*models.Usuario, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener usuario con ID nulo")
		return nil, errors.New("el ID del usuario es requerido")
	}
	usuario, err := s.store.GetUsuarioByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if permisos, err := s.store.GetPermisosByRol(ctx, usuario.IDRol); err == nil {
		usuario.Permisos = permisos
	}
	return usuario, nil
}

// CreateUsuario crea un nuevo usuario con validaciones de negocio
func (s *ServiceUsuario) CreateUsuario(ctx context.Context, usuario models.Usuario) (*models.Usuario, error) {
	if err := s.validateUser(&usuario, false); err != nil {
		s.logger.WarnContext(ctx, "Fallo de validación al crear usuario", slog.String("email", usuario.Email), slog.Any("error", err))
		return nil, err
	}

	// Aqui se hashea la contraseña
	hashearPassword, err := helper.HashPassword(usuario.Password)
	if err != nil {
		return nil, errors.New("error al generar hash de la contraseña")
	}
	usuario.Password = hashearPassword

	return s.store.CreateUsuario(ctx, &usuario)
}

// UpdateUsuario actualiza un usuario existente con validaciones
func (s *ServiceUsuario) UpdateUsuario(ctx context.Context, id uuid.UUID, usuario models.Usuario) (*models.Usuario, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualización de usuario con ID nulo")
		return nil, errors.New("el ID del usuario es requerido")
	}

	if err := s.validateUser(&usuario, true); err != nil {
		s.logger.WarnContext(ctx, "Fallo de validación al actualizar usuario", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}

	// SOLO si viene contraseña nueva → hashear
	if usuario.Password != "" {
		hashearPasword, err := helper.HashPassword(usuario.Password)
		if err != nil {
			return nil, errors.New("error al hashear la contraseña")
		}
		usuario.Password = hashearPasword
	}

	return s.store.UpdateUsuario(ctx, id, &usuario)
}

// DeleteUsuario elimina un usuario (soft delete)
func (s *ServiceUsuario) DeleteUsuario(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminación de usuario con ID nulo")
		return errors.New("el ID del usuario es requerido")
	}
	return s.store.DeleteUsuario(ctx, id)
}

// AdministrarUsuario gestiona la creación/actualización integral del usuario, incluyendo accesos multi-sucursal (Supermercado)
func (s *ServiceUsuario) AdministrarUsuario(ctx context.Context, usuario models.Usuario, adminID uuid.UUID) (*models.Usuario, error) {
	var result *models.Usuario
	var err error

	// 1. Validar y Procesar Password
	if usuario.Password != "" {
		hp, err := helper.HashPassword(usuario.Password)
		if err != nil {
			return nil, fmt.Errorf("error al hashear password: %w", err)
		}
		usuario.Password = hp
	}

	// 2. Ejecutar Operación Principal (Create o Update)
	if usuario.IDUsuario == uuid.Nil {
		result, err = s.store.CreateUsuario(ctx, &usuario)
	} else {
		result, err = s.store.UpdateUsuario(ctx, usuario.IDUsuario, &usuario)
	}

	if err != nil {
		return nil, err
	}

	// 3. Gestionar Accesos Multi-Sucursal
	if len(usuario.Sucursales) > 0 {
		if err := s.store.AssignSucursales(ctx, result.IDUsuario, usuario.Sucursales); err != nil {
			s.logger.ErrorContext(ctx, "Error asignando sucursales", slog.Any("error", err))
		}
	}

	// 4. Auditoría
	s.logsStore.CreateLog(ctx, &models.LogSistema{
		IDUsuario:  adminID,
		Accion:     "ADMINISTRAR_USUARIO",
		Tabla:      "usuario",
		RegistroID: result.IDUsuario,
		IP:         middleware.GetClientIP(ctx),
	})

	return result, nil
}

// AuthenticateUsuario valida las credenciales del usuario y retorna el usuario autenticado
func (s *ServiceUsuario) AuthenticateUsuario(ctx context.Context, req models.LoginRequest) (*models.Usuario, error) {
	var usuario *models.Usuario
	var err error

	// 1. Identificar el método de búsqueda
	if req.Pin != "" {
		// Login por PIN (Acceso rápido POS)
		usuario, err = s.store.GetUsuarioByPin(ctx, req.Pin)
		if err != nil {
			return nil, errors.New("PIN inválido")
		}
	} else if req.Email != "" {
		// Login por Email
		usuario, err = s.store.GetUsuarioByEmail(ctx, strings.TrimSpace(req.Email))
	} else if req.Username != "" {
		// Login por Username
		usuario, err = s.store.GetUsuarioByUsername(ctx, strings.TrimSpace(req.Username))
	} else {
		return nil, errors.New("debe proporcionar email, username o pin")
	}

	if err != nil {
		s.logger.WarnContext(ctx, "Login fallido: usuario no encontrado", slog.Any("error", err))
		return nil, errors.New("credenciales inválidas")
	}

	// 2. Validar estatus activo
	estatusActivo := uuid.MustParse("3a99d245-b34f-48a5-ac08-a5a010c5822f")
	if usuario.IDStatus != estatusActivo {
		return nil, errors.New("su cuenta no está activa")
	}

	// 3. Validar Password (solo si NO es login por PIN)
	if req.Pin == "" {
		if req.Password == "" {
			return nil, errors.New("la contraseña es requerida")
		}
		if err := helper.CheckPassword(req.Password, usuario.Password); err != nil {
			return nil, errors.New("credenciales inválidas")
		}
	}

	// 4. Limpiar password y cargar permisos
	usuario.Password = ""
	permisos, err := s.store.GetPermisosByRol(ctx, usuario.IDRol)
	if err == nil {
		usuario.Permisos = permisos
	}

	return usuario, nil
}
