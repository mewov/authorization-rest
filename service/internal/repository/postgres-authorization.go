package repository

import (
	"time"

	"github.com/mewov/authorization-rest/models"
)

func (p *Postgres) RegisterUser(login, email, password, client, role string) (int64, error) {
	var id int64
	err := p.conn.QueryRow("INSERT INTO users (login, email, password, client, role, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;",
		login, email, password, client, role, time.Now().Unix()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *Postgres) SearchUser(login string) (*models.PostgresUser, error) {
	var user models.PostgresUser
	err := p.conn.Get(&user, "SELECT * FROM users WHERE login = $1", login)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *Postgres) SearchUserByID(id int64) (*models.PostgresUser, error) {
	var user models.PostgresUser
	err := p.conn.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) RemoveUser(id int64) error {
	_, err := p.conn.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
