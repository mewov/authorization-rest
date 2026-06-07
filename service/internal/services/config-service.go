package services

import (
	"os"

	"github.com/mewov/authorization-rest/models"
)

func NewConfig(model *models.LocalConfig) {
	*model = models.LocalConfig{
		SERVER_ADDRESS:    os.Getenv("SERVER_ADDRESS"),
		SERVER_PASSWORD:   os.Getenv("SERVER_PASSWORD"),
		POSTGRES_ADDRESS:  os.Getenv("POSTGRES_ADDRESS"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
	}
}
