package dto

// ConfigResponseDTO es una respuesta genérica para los catálogos
type ConfigResponseDTO struct {
	ID          string `json:"id"`
	Descripcion string `json:"descripcion"`
}

// EstacionConfiguracionDTO agrupa todos los catálogos en un solo objeto si fuera necesario
type EstacionConfiguracionDTO struct {
	Canales    []ConfigResponseDTO `json:"canales"`
	Impresoras []ConfigResponseDTO `json:"impresoras"`
	Puertos    []ConfigResponseDTO `json:"puertos"`
}
