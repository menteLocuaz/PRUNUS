package services

import (
	"context"
	"errors"
	"log/slog"
	"regexp"

	"github.com/google/uuid"
	"github.com/prunus/pkg/helper"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

// ServiceUsuario servicio que encapsula la lógica de negocio para usuario
type ServiceUsuario struct {
	store  store.StoreUsuario
	logger *slog.Logger
}

// NewServiceUsuario crea una nueva instancia del servicio de usuario
func NewServiceUsuario(s store.StoreUsuario, logger *slog.Logger) *ServiceUsuario {
	return &ServiceUsuario{
		store:  s,
		logger: logger,
	}
}

// GetAllUsuarios obtiene todos los usuarios del sistema
func (s *ServiceUsuario) GetAllUsuarios(ctx context.Context) ([]*models.Usuario, error) {
	return s.store.GetAllUsuarios(ctx)
}

// GetUsuarioByID obtiene un usuario por su ID
func (s *ServiceUsuario) GetUsuarioByID(ctx context.Context, id uuid.UUID) (*models.Usuario, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener usuario con ID nulo")
		return nil, errors.New("el ID del usuario es requerido")
	}
	return s.store.GetUsuarioByID(ctx, id)
}

// CreateUsuario crea un nuevo usuario con validaciones de negocio
func (s *ServiceUsuario) CreateUsuario(ctx context.Context, usuario models.Usuario) (*models.Usuario, error) {
	// Validar campos obligatorios
	if usuario.Email == "" {
		s.logger.WarnContext(ctx, "Intento de creación de usuario con email vacío")
		return nil, errors.New("el email del usuario es requerido")
	}
	if usuario.UsuNombre == "" {
		s.logger.WarnContext(ctx, "Intento de creación de usuario con nombre vacío", slog.String("email", usuario.Email))
		return nil, errors.New("el nombre del usuario es requerido")
	}
	if usuario.UsuDNI == "" {
		s.logger.WarnContext(ctx, "Intento de creación de usuario con DNI vacío", slog.String("email", usuario.Email))
		return nil, errors.New("el DNI del usuario es requerido")
	}
	if usuario.Password == "" {
		s.logger.WarnContext(ctx, "Intento de creación de usuario con password vacío", slog.String("email", usuario.Email))
		return nil, errors.New("la contraseña del usuario es requerida")
	}
	if usuario.IDSucursal == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de usuario sin sucursal", slog.String("email", usuario.Email))
		return nil, errors.New("el ID de la sucursal es requerido")
	}
	if usuario.IDRol == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de usuario sin rol", slog.String("email", usuario.Email))
		return nil, errors.New("el ID del rol es requerido")
	}

	// Validar formato de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(usuario.Email) {
		s.logger.WarnContext(ctx, "Intento de creación de usuario con formato de email inválido", slog.String("email", usuario.Email))
		return nil, errors.New("el formato del email es inválido")
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
	if usuario.Email == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de usuario con email vacío", slog.String("id", id.String()))
		return nil, errors.New("el email del usuario es requerido")
	}
	if usuario.UsuNombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de usuario con nombre vacío", slog.String("id", id.String()))
		return nil, errors.New("el nombre del usuario es requerido")
	}
	if usuario.UsuDNI == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de usuario con DNI vacío", slog.String("id", id.String()))
		return nil, errors.New("el DNI del usuario es requerido")
	}
	if usuario.IDSucursal == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualización de usuario sin sucursal", slog.String("id", id.String()))
		return nil, errors.New("el ID de la sucursal es requerido")
	}

	// Validar formato de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(usuario.Email) {
		s.logger.WarnContext(ctx, "Intento de actualización de usuario con formato de email inválido", slog.String("id", id.String()), slog.String("email", usuario.Email))
		return nil, errors.New("el formato del email es inválido")
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

// AuthenticateUsuario valida las credenciales del usuario y retorna el usuario autenticado
func (s *ServiceUsuario) AuthenticateUsuario(ctx context.Context, email, password string) (*models.Usuario, error) {
	// Validar que se proporcionen ambos campos
	if email == "" {
		s.logger.WarnContext(ctx, "Intento de login con email vacío")
		return nil, errors.New("el email es requerido")
	}
	if password == "" {
		s.logger.WarnContext(ctx, "Intento de login con password vacío", slog.String("email", email))
		return nil, errors.New("la contraseña es requerida")
	}

	// Validar formato de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		s.logger.WarnContext(ctx, "Intento de login con formato de email inválido", slog.String("email", email))
		return nil, errors.New("el formato del email es inválido")
	}

	// Buscar usuario por email
	usuario, err := s.store.GetUsuarioByEmail(ctx, email)
	if err != nil {
		// No revelar si el usuario existe o no por seguridad
		s.logger.WarnContext(ctx, "Login fallido: usuario no encontrado", slog.String("email", email))
		return nil, errors.New("credenciales inválidas")
	}

	// Verificar la contraseña usando bcrypt
	err = helper.CheckPassword(password, usuario.Password)
	if err != nil {
		// Password incorrecta
		s.logger.WarnContext(ctx, "Login fallido: contraseña incorrecta", slog.String("email", email), slog.String("id_usuario", usuario.IDUsuario.String()))
		return nil, errors.New("credenciales inválidas")
	}

	// Limpiar el password del objeto antes de retornarlo
	usuario.Password = ""

	return usuario, nil
}
