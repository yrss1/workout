package postgres

import (
	"context"
	"database/sql"

	"github.com/yrss1/workout/internal/pkg/log"
)

type database interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type logger interface {
	Error(msg string, fields ...log.Field)
}
