package services

import (
	"errors"
	"regexp"

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
func (s *ServiceUsuario) GetUsuarioByID(id uint) (*models.Usuario, error) {
	if id == 0 {
		return nil, errors.New("el ID del usuario es requerido")
	}
	return s.store.GetUsuarioByID(id)
}

// CreateUsuario crea un nuevo usuario con validaciones de negocio
func (s *ServiceUsuario) CreateUsuario(usuario models.Usuario) (*models.Usuario, error) {
	// Validar campos obligatorios
	if usuario.UsuEmail == "" {
		return nil, errors.New("el email del usuario es requerido")
	}
	if usuario.UsuNombre == "" {
		return nil, errors.New("el nombre del usuario es requerido")
	}
	if usuario.UsuDni == "" {
		return nil, errors.New("el DNI del usuario es requerido")
	}
	if usuario.UsuPassword == "" {
		return nil, errors.New("la contraseña del usuario es requerida")
	}
	if usuario.IDSucursal == 0 {
		return nil, errors.New("el ID de la sucursal es requerido")
	}

	// Validar formato de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(usuario.UsuEmail) {
		return nil, errors.New("el formato del email es inválido")
	}

	// Establecer estado por defecto si no está definido
	if usuario.Estado == 0 {
		usuario.Estado = 1
	}

	// Aqui se hashea la contraseña
	hashearPassword, err := helper.HashPassword(usuario.UsuPassword)
	if err != nil {
		return nil, errors.New("error al generar hash de la contraseña")
	}
	usuario.UsuPassword = hashearPassword

	return s.store.CreateUsuario(&usuario)
}

// UpdateUsuario actualiza un usuario existente con validaciones
func (s *ServiceUsuario) UpdateUsuario(id uint, usuario models.Usuario) (*models.Usuario, error) {
	if id == 0 {
		return nil, errors.New("el ID del usuario es requerido")
	}
	if usuario.UsuEmail == "" {
		return nil, errors.New("el email del usuario es requerido")
	}
	if usuario.UsuNombre == "" {
		return nil, errors.New("el nombre del usuario es requerido")
	}
	if usuario.UsuDni == "" {
		return nil, errors.New("el DNI del usuario es requerido")
	}
	if usuario.IDSucursal == 0 {
		return nil, errors.New("el ID de la sucursal es requerido")
	}

	// Validar formato de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(usuario.UsuEmail) {
		return nil, errors.New("el formato del email es inválido")
	}

	// SOLO si viene contraseña nueva → hashear
	if usuario.UsuPassword != "" {
		hashearPasword, err := helper.HashPassword(usuario.UsuPassword)
		if err != nil {
			return nil, errors.New("error al hashear la contraseña")
		}
		usuario.UsuPassword = hashearPasword
	} else {
		// Evita sobreescribir contraseña en DB
		usuario.UsuPassword = ""
	}

	return s.store.UpdateUsuario(id, &usuario)
}

// DeleteUsuario elimina un usuario (soft delete)
func (s *ServiceUsuario) DeleteUsuario(id uint) error {
	if id == 0 {
		return errors.New("el ID del usuario es requerido")
	}
	return s.store.DeleteUsuario(id)
}
