package http

import (
	"github.com/yrss1/workout/internal/pkg/log"
)

type logger interface {
	Error(msg string, fields ...log.Field)
	Info(msg string, fields ...log.Field)
}

// type authService interface {
// 	ValidateToken(context.Context, string) error
// 	PutUserInContext(context.Context, *model.User) context.Context
// 	GetUserFromToken(context.Context, string) (*model.User, error)
// }
