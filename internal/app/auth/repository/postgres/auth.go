package postgres

import (
	"time"
)

const (
	timeout = 500 * time.Millisecond
)

type AuthRepository struct {
	logger   logger
	database database
}

func NewAuthRepository(logger logger, database database) *AuthRepository {
	return &AuthRepository{logger: logger, database: database}
}
