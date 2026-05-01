package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/helper"
	"github.com/prunus/pkg/middleware"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	zaplogger "github.com/prunus/pkg/utils/logger"
	"go.uber.org/zap"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ServiceUsuario servicio que encapsula la lógica de negocio para usuario
type ServiceUsuario struct {
	store      store.StoreUsuario
	rolService *ServiceRol
	logsStore  store.StoreLogs
	logger     *zap.Logger
}

// NewServiceUsuario crea una nueva instancia del servicio de usuario
func NewServiceUsuario(s store.StoreUsuario, r *ServiceRol, l store.StoreLogs, logger *zap.Logger) *ServiceUsuario {
	return &ServiceUsuario{
		store:      s,
		rolService: r,
		logsStore:  l,
		logger:     logger,
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

// GetAllUsuarios obtiene una lista paginada de usuarios del sistema
func (s *ServiceUsuario) GetAllUsuarios(ctx context.Context, params dto.PaginationParams) ([]*models.Usuario, error) {
	return s.store.GetAllUsuarios(ctx, params)
}

// GetUsuarioByID obtiene un usuario por su ID con sus permisos cargados (vía cache)
func (s *ServiceUsuario) GetUsuarioByID(ctx context.Context, id uuid.UUID) (*models.Usuario, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de obtener usuario con ID nulo")
		return nil, errors.New("el ID del usuario es requerido")
	}
	usuario, err := s.store.GetUsuarioByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if permisos, err := s.rolService.GetPermisosByRol(ctx, usuario.IDRol); err == nil {
		usuario.Permisos = permisos
	}

	return usuario, nil
}

// CreateUsuario crea un nuevo usuario con validaciones de negocio
func (s *ServiceUsuario) CreateUsuario(ctx context.Context, usuario models.Usuario) (*models.Usuario, error) {
	if err := s.validateUser(&usuario, false); err != nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Fallo de validación al crear usuario",
			zap.String("email", usuario.Email),
			zap.Error(err),
		)
		return nil, err
	}

	hashearPassword, err := helper.HashPassword(usuario.Password)
	if err != nil {
		return nil, errors.New("error al generar hash de la contraseña")
	}
	usuario.Password = hashearPassword

	result, err := s.store.CreateUsuario(ctx, &usuario)
	if err != nil {
		return nil, err
	}
	result.ClearSensitiveFields()
	return result, nil
}

// UpdateUsuario actualiza un usuario existente con validaciones
func (s *ServiceUsuario) UpdateUsuario(ctx context.Context, id uuid.UUID, usuario models.Usuario) (*models.Usuario, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de actualización de usuario con ID nulo")
		return nil, errors.New("el ID del usuario es requerido")
	}

	if err := s.validateUser(&usuario, true); err != nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Fallo de validación al actualizar usuario",
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	if usuario.Password != "" {
		hashearPasword, err := helper.HashPassword(usuario.Password)
		if err != nil {
			return nil, errors.New("error al hashear la contraseña")
		}
		usuario.Password = hashearPasword
	}

	result, err := s.store.UpdateUsuario(ctx, id, &usuario)
	if err != nil {
		return nil, err
	}
	result.ClearSensitiveFields()
	return result, nil
}

// DeleteUsuario elimina un usuario (soft delete)
func (s *ServiceUsuario) DeleteUsuario(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de eliminación de usuario con ID nulo")
		return errors.New("el ID del usuario es requerido")
	}
	return s.store.DeleteUsuario(ctx, id)
}

// AdministrarUsuario gestiona la creación/actualización integral del usuario, incluyendo accesos multi-sucursal
func (s *ServiceUsuario) AdministrarUsuario(ctx context.Context, usuario models.Usuario, adminID uuid.UUID) (*models.Usuario, error) {
	var result *models.Usuario
	var err error

	if usuario.Password != "" {
		hp, err := helper.HashPassword(usuario.Password)
		if err != nil {
			return nil, fmt.Errorf("error al hashear password: %w", err)
		}
		usuario.Password = hp
	}

	if usuario.IDUsuario == uuid.Nil {
		result, err = s.store.CreateUsuario(ctx, &usuario)
	} else {
		result, err = s.store.UpdateUsuario(ctx, usuario.IDUsuario, &usuario)
	}

	if err != nil {
		return nil, err
	}

	if len(usuario.Sucursales) > 0 {
		if err := s.store.AssignSucursales(ctx, result.IDUsuario, usuario.Sucursales); err != nil {
			zaplogger.WithContext(ctx, s.logger).Error("Error asignando sucursales", zap.Error(err))
		}
	}

	s.logsStore.CreateLog(ctx, &models.LogSistema{
		IDUsuario:  adminID,
		Accion:     "ADMINISTRAR_USUARIO",
		Tabla:      "usuario",
		RegistroID: result.IDUsuario,
		IP:         middleware.GetClientIP(ctx),
	})

	result.ClearSensitiveFields()
	return result, nil
}

// AuthenticateUsuario valida las credenciales del usuario y retorna el usuario autenticado
func (s *ServiceUsuario) AuthenticateUsuario(ctx context.Context, req models.LoginRequest) (*models.Usuario, error) {
	var usuario *models.Usuario
	var err error

	zaplogger.WithContext(ctx, s.logger).Info("[LOGIN] Intento de inicio de sesión",
		zap.String("email", req.Email),
		zap.String("username", req.Username),
	)

	if req.Pin != "" {
		usuario, err = s.store.GetUsuarioByPin(ctx, req.Pin)
	} else if req.Email != "" {
		usuario, err = s.store.GetUsuarioByEmail(ctx, strings.TrimSpace(req.Email))
	} else if req.Username != "" {
		usuario, err = s.store.GetUsuarioByUsername(ctx, strings.TrimSpace(req.Username))
	} else {
		return nil, errors.New("debe proporcionar email, username o pin")
	}

	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Warn("[LOGIN] Usuario no encontrado o error en DB", zap.Error(err))
		return nil, errors.New("credenciales inválidas")
	}

	if usuario.IDStatus != models.EstatusGlobalActivo {
		zaplogger.WithContext(ctx, s.logger).Warn("[LOGIN] Usuario inactivo",
			zap.String("id_usuario", usuario.IDUsuario.String()),
			zap.String("status_actual", usuario.IDStatus.String()),
			zap.String("status_esperado", models.EstatusGlobalActivo.String()),
		)
		return nil, errors.New("su cuenta no está activa")
	}

	if req.Pin == "" {
		if req.Password == "" {
			return nil, errors.New("la contraseña es requerida")
		}
		if err := helper.CheckPassword(req.Password, usuario.Password); err != nil {
			zaplogger.WithContext(ctx, s.logger).Warn("[LOGIN] Contraseña incorrecta", zap.String("id_usuario", usuario.IDUsuario.String()))
			return nil, errors.New("credenciales inválidas")
		}
	}

	zaplogger.WithContext(ctx, s.logger).Info("[LOGIN] Autenticación exitosa", zap.String("id_usuario", usuario.IDUsuario.String()))

	usuario.ClearSensitiveFields()
	if permisos, err := s.rolService.GetPermisosByRol(ctx, usuario.IDRol); err == nil {
		usuario.Permisos = permisos
	}

	return usuario, nil
}
