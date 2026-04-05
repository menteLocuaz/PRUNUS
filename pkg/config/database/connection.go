package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prunus/pkg/config"
	"github.com/redis/go-redis/v9"
)

// InitDB abre la conexión a la base de datos usando la configuración de Viper.
func InitDB() (*sql.DB, error) {
	host := config.Get("DB_HOST")
	user := config.Get("DB_USER")
	password := config.Get("DB_PASSWORD")
	dbname := config.Get("DB_NAME")
	port := config.GetDefault("DB_PORT", "5432")
	sslmode := config.GetDefault("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	maxOpenConns, _ := strconv.Atoi(config.GetDefault("DB_POOL_MAX_OPEN", "25"))
	maxIdleConns, _ := strconv.Atoi(config.GetDefault("DB_POOL_MAX_IDLE", "5"))
	connMaxLifetimeMin, _ := strconv.Atoi(config.GetDefault("DB_POOL_CONN_MAX_LIFETIME_MIN", "60"))

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(connMaxLifetimeMin) * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al conectar (ping) a PostgreSQL: %w", err)
	}

	return db, nil
}

// InitRedis abre la conexión a Redis usando la configuración de Viper.
func InitRedis() (*redis.Client, error) {
	host := config.GetDefault("REDIS_HOST", "localhost")
	port := config.GetDefault("REDIS_PORT", "6379")
	password := config.Get("REDIS_PASSWORD")
	db := config.GetInt("REDIS_DB")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return rdb, rdb.Ping(ctx).Err()
}
