package services

import (
	"database/sql"
	"errors"

	"github.com/mewov/authorization-rest/internal/repository"
	"github.com/mewov/authorization-rest/models"
	"github.com/mewov/authorization-rest/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthorizathionService struct {
		db *repository.Postgres
	}
)

func NewAuthorizathion(db *repository.Postgres, model *AuthorizathionService) {
	*model = AuthorizathionService{
		db: db,
	}
}

func (auth *AuthorizathionService) CreateUser(login, email, password, client, role string) (int64, error) {
	if err := validator.ValidateLogin(login); err != nil {
		return 0, err
	}
	if err := validator.Email(email); err != nil {
		return 0, err
	}
	if err := validator.ValidatePassword(password); err != nil {
		return 0, err
	}

	if client == "" {
		client = "none"
	}
	if role == "" {
		role = "user"
	}

	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	lastId, err := auth.db.RegisterUser(login, email, string(bcryptHash), client, role)
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (auth *AuthorizathionService) SearchUser(login string, password string) (*models.PostgresUser, error) {
	if err := validator.ValidateLogin(login); err != nil {
		return nil, err
	}
	if err := validator.ValidatePassword(password); err != nil {
		return nil, err
	}

	user, err := auth.db.SearchUser(login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (auth *AuthorizathionService) SearchUserByID(id int64) (*models.PostgresUser, error) {
	user, err := auth.db.SearchUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (auth *AuthorizathionService) RemoveUser(id int64) error {
	return auth.db.RemoveUser(id)
}
