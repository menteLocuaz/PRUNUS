package database

import (
	"database/sql"
	"fmt"
	"os"
	"time" // Importado para configurar el tiempo de vida de las conexiones

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

// Conexion establece y configura la conexión a la base de datos PostgreSQL
func Conexion() (*sql.DB, error) {
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		// No retornamos error aquí por si las variables ya están en el sistema (ej. Docker/Heroku)
		fmt.Println("Aviso: No se pudo cargar el archivo .env, usando variables de entorno del sistema")
	}

	// Leer variables de configuración
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	// Construir la cadena de conexión (DSN)
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	// Abrir la conexión usando el driver pgx
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	// --- Configuración del Pool de Conexiones ---

	// SetMaxOpenConns establece el número máximo de conexiones abiertas a la base de datos.
	db.SetMaxOpenConns(25)

	// SetMaxIdleConns establece el número máximo de conexiones en el pool de conexiones inactivas.
	db.SetMaxIdleConns(25)

	// SetConnMaxLifetime establece el tiempo máximo que una conexión puede ser reutilizada.
	// Esto ayuda a evitar problemas con conexiones que se vuelven inestables con el tiempo.
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar que la conexión sea válida
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al conectar (ping) a PostgreSQL: %w", err)
	}

	fmt.Println("✅ Conectado a PostgreSQL y Pool configurado")
	return db, nil
}
