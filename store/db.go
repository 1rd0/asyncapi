package store

import (
	"asyncapi/config"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func NewPostgresDb(conf *config.Config) (db *sql.DB, err error) { // Принимаем конфиг в процессе создаем подлкючение к бд
	dsn := conf.DatabaseURL()           //  ссылка к базе данных
	db, err = sql.Open("postgres", dsn) // метод опен подклюения

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // создакем контекс который оставноит приложение если не сможем
	//пингануть к бд
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
