package services

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"go.uber.org/zap"
)

// --- MOCKS MANUALES ---

type mockPeriodoStore struct {
	activePeriodo *models.Periodo
	createErr     error
	getCalled     int
	createCalled  int
}

func (m *mockPeriodoStore) GetActivePeriodo(ctx context.Context, sucursalID uuid.UUID) (*models.Periodo, error) {
	m.getCalled++
	return m.activePeriodo, nil
}

func (m *mockPeriodoStore) CreatePeriodo(ctx context.Context, p *models.Periodo) (*models.Periodo, error) {
	m.createCalled++
	if m.createErr != nil {
		return nil, m.createErr
	}
	return p, nil
}

func (m *mockPeriodoStore) CerrarPeriodo(ctx context.Context, id uuid.UUID, idUsuarioCierre uuid.UUID, ip string) error {
	return nil
}

func (m *mockPeriodoStore) GetPeriodoByID(ctx context.Context, id uuid.UUID) (*models.Periodo, error) { return nil, nil }
func (m *mockPeriodoStore) GetAllPeriodos(ctx context.Context) ([]*models.Periodo, error)           { return nil, nil }
func (m *mockPeriodoStore) UpdatePeriodo(ctx context.Context, id uuid.UUID, p *models.Periodo) (*models.Periodo, error) {
	return nil, nil
}
func (m *mockPeriodoStore) DeletePeriodo(ctx context.Context, id uuid.UUID) error { return nil }

type mockPOSStore struct {
	activeControls int
}

func (m *mockPOSStore) GetTotalActiveControls(ctx context.Context) (int, error) {
	return m.activeControls, nil
}

func (m *mockPOSStore) GetActiveControlByEstacion(ctx context.Context, id uuid.UUID) (*models.ControlEstacion, error) {
	return nil, nil
}
func (m *mockPOSStore) GetEstadoCompletoEstacion(ctx context.Context, id uuid.UUID) (*dto.EstadoCajaDTO, error) {
	return nil, nil
}
func (m *mockPOSStore) CreateControlEstacion(ctx context.Context, c *models.ControlEstacion) (*models.ControlEstacion, error) {
	return nil, nil
}
func (m *mockPOSStore) UpdateControlEstacion(ctx context.Context, c *models.ControlEstacion) error { return nil }
func (m *mockPOSStore) GetEstacionByID(ctx context.Context, id uuid.UUID) (*models.EstacionPos, error) {
	return nil, nil
}
func (m *mockPOSStore) UpdateEstacionStatus(ctx context.Context, id uuid.UUID, s uuid.UUID) error {
	return nil
}
func (m *mockPOSStore) GetActivePeriodo(ctx context.Context) (*models.Periodo, error) { return nil, nil }
func (m *mockPOSStore) DesmontarCajero(ctx context.Context, id uuid.UUID, s1, s2, s3 uuid.UUID, mot string) error {
	return nil
}
func (m *mockPOSStore) UpdateValoresDeclarados(ctx context.Context, c, f, u uuid.UUID, v float64, t int, s1, s2, s3 uuid.UUID) error {
	return nil
}

// --- PRUEBAS CON ASSERT NATIVO ---

func TestAbrirNuevoPeriodo_Exito(t *testing.T) {
	store := &mockPeriodoStore{activePeriodo: nil}
	service := NewServicePeriodo(store, &mockPOSStore{}, zap.NewNop())

	uID := uuid.New()
	sID := uuid.New()
	ip := "192.168.1.1"
	motivo := "Apertura de turno mañana"

	res, err := service.AbrirNuevoPeriodo(context.Background(), uID, sID, ip, motivo)

	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo %v", err)
	}
	if res == nil {
		t.Fatal("se esperaba resultado no nil")
	}
	if res.IDSucursal != sID {
		t.Errorf("sucursal incorrecta: %v != %v", res.IDSucursal, sID)
	}
	if res.PrdIPApertura != ip {
		t.Errorf("IP incorrecta: %s != %s", res.PrdIPApertura, ip)
	}
	if store.createCalled != 1 {
		t.Errorf("se esperaba 1 llamada a Create, se obtuvieron %d", store.createCalled)
	}
}

func TestAbrirNuevoPeriodo_YaExisteActivo(t *testing.T) {
	existente := &models.Periodo{IDPeriodo: uuid.New(), IDSucursal: uuid.New()}
	store := &mockPeriodoStore{activePeriodo: existente}
	service := NewServicePeriodo(store, &mockPOSStore{}, zap.NewNop())

	res, err := service.AbrirNuevoPeriodo(context.Background(), uuid.New(), existente.IDSucursal, "1.1.1.1", "")

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.IDPeriodo != existente.IDPeriodo {
		t.Error("debería devolver el periodo ya existente")
	}
	if store.createCalled != 0 {
		t.Error("no debería haber creado un nuevo periodo")
	}
}

func TestFinalizarPeriodo_ErrorEstacionesAbiertas(t *testing.T) {
	store := &mockPeriodoStore{}
	posStore := &mockPOSStore{activeControls: 5}
	service := NewServicePeriodo(store, posStore, zap.NewNop())

	err := service.FinalizarPeriodo(context.Background(), uuid.New(), uuid.New(), "1.1.1.1")

	if err == nil {
		t.Fatal("se esperaba error por cajas abiertas")
	}
	if !strings.Contains(err.Error(), "hay cajas o estaciones abiertas") {
		t.Errorf("mensaje de error incorrecto: %v", err)
	}
}
