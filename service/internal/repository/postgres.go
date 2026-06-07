package repository

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/mewov/authorization-rest/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type (
	Postgres struct {
		conn   *sqlx.DB
		logger *slog.Logger
	}
)

func NewPostgres(config *models.LocalConfig, logger *slog.Logger, model *Postgres) {
	pcs := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.POSTGRES_USER, config.POSTGRES_PASSWORD,
		config.POSTGRES_ADDRESS, config.POSTGRES_DB,
	)

	var connect *sqlx.DB
	var err error
	for range 10 {
		connect, err = sqlx.Connect("postgres", pcs)
		if err == nil {
			logger.Info("postgres connect:", slog.String("addr", config.POSTGRES_ADDRESS))
			fmt.Println("[+] postgres.connect:", config.POSTGRES_ADDRESS)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		logger.Error(err.Error(), slog.String("addr", config.POSTGRES_ADDRESS))
		fmt.Println("[-] postgres.connect:", err.Error())
		return
	}

	err = connect.Ping()
	if err != nil {
		logger.Error(err.Error(), slog.String("addr", config.POSTGRES_ADDRESS))
		fmt.Println("[-] postgres.ping:", err.Error())
		return
	}

	logger.Info("postgres ping:", slog.String("addr", config.POSTGRES_ADDRESS))
	fmt.Println("[+] postgres.ping:", config.POSTGRES_ADDRESS)

	*model = Postgres{
		conn:   connect,
		logger: logger,
	}
}

func (p *Postgres) Migration() {
	_, err := p.conn.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		client VARCHAR(255),
		role VARCHAR(255) DEFAULT 'user',
		created_at INTEGER
    )`)
	if err != nil {
		fmt.Println("[-] create table:", err.Error())
	}

	_, err = p.conn.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		token VARCHAR(255) UNIQUE NOT NULL,
    	client VARCHAR(255) NOT NULL,
    	expires_at INTEGER NOT NULL,
    	created_at INTEGER NOT NULL
	)`)
	if err != nil {
		fmt.Println("[-] create table:", err.Error())
	}
}
