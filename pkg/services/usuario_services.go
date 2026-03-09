package services

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"github.com/prunus/pkg/helper"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

// ServiceUsuario servicio que encapsula la lógica de negocio para usuario
type ServiceUsuario struct {
	store store.StoreUsuario
}

// NewServiceUsuario crea una nueva instancia del servicio de usuario
func NewServiceUsuario(s store.StoreUsuario) *ServiceUsuario {
	return &ServiceUsuario{store: s}
}

// GetAllUsuarios obtiene todos los usuarios del sistema
func (s *ServiceUsuario) GetAllUsuarios() ([]*models.Usuario, error) {
	return s.store.GetAllUsuarios()
}

// GetUsuarioByID obtiene un usuario por su ID
func (s *ServiceUsuario) GetUsuarioByID(id uuid.UUID) (*models.Usuario, error) {
	if id == uuid.Nil {
		return nil, errors.New("el ID del usuario es requerido")
	}
	return s.store.GetUsuarioByID(id)
}

// CreateUsuario crea un nuevo usuario con validaciones de negocio
func (s *ServiceUsuario) CreateUsuario(usuario models.Usuario) (*models.Usuario, error) {
	// Validar campos obligatorios
	if usuario.Email == "" {
		return nil, errors.New("el email del usuario es requerido")
	}
	if usuario.UsuNombre == "" {
		return nil, errors.New("el nombre del usuario es requerido")
	}
	if usuario.UsuDNI == "" {
		return nil, errors.New("el DNI del usuario es requerido")
	}
	if usuario.Password == "" {
		return nil, errors.New("la contraseña del usuario es requerida")
	}
	if usuario.IDSucursal == uuid.Nil {
		return nil, errors.New("el ID de la sucursal es requerido")
	}
	if usuario.IDRol == uuid.Nil {
		return nil, errors.New("el ID del rol es requerido")
	}

	// Validar formato de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(usuario.Email) {
		return nil, errors.New("el formato del email es inválido")
	}

	// Aqui se hashea la contraseña
	hashearPassword, err := helper.HashPassword(usuario.Password)
	if err != nil {
		return nil, errors.New("error al generar hash de la contraseña")
	}
	usuario.Password = hashearPassword

	return s.store.CreateUsuario(&usuario)
}

// UpdateUsuario actualiza un usuario existente con validaciones
func (s *ServiceUsuario) UpdateUsuario(id uuid.UUID, usuario models.Usuario) (*models.Usuario, error) {
	if id == uuid.Nil {
		return nil, errors.New("el ID del usuario es requerido")
	}
	if usuario.Email == "" {
		return nil, errors.New("el email del usuario es requerido")
	}
	if usuario.UsuNombre == "" {
		return nil, errors.New("el nombre del usuario es requerido")
	}
	if usuario.UsuDNI == "" {
		return nil, errors.New("el DNI del usuario es requerido")
	}
	if usuario.IDSucursal == uuid.Nil {
		return nil, errors.New("el ID de la sucursal es requerido")
	}

	// Validar formato de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(usuario.Email) {
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

	return s.store.UpdateUsuario(id, &usuario)
}

// DeleteUsuario elimina un usuario (soft delete)
func (s *ServiceUsuario) DeleteUsuario(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("el ID del usuario es requerido")
	}
	return s.store.DeleteUsuario(id)
}

// AuthenticateUsuario valida las credenciales del usuario y retorna el usuario autenticado
func (s *ServiceUsuario) AuthenticateUsuario(email, password string) (*models.Usuario, error) {
	// Validar que se proporcionen ambos campos
	if email == "" {
		return nil, errors.New("el email es requerido")
	}
	if password == "" {
		return nil, errors.New("la contraseña es requerida")
	}

	// Validar formato de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return nil, errors.New("el formato del email es inválido")
	}

	// Buscar usuario por email
	usuario, err := s.store.GetUsuarioByEmail(email)
	if err != nil {
		// No revelar si el usuario existe o no por seguridad
		return nil, errors.New("credenciales inválidas")
	}

	// Verificar la contraseña usando bcrypt
	err = helper.CheckPassword(password, usuario.Password)
	if err != nil {
		// Password incorrecta
		return nil, errors.New("credenciales inválidas")
	}

	// Limpiar el password del objeto antes de retornarlo
	usuario.Password = ""

	return usuario, nil
}
