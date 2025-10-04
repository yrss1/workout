package http

import (
	"errors"
	"net/http"

	"github.com/yrss1/workout/internal/pkg/errs"
	"github.com/yrss1/workout/internal/pkg/log"
)

type PanicRecoverer struct {
	logger logger
}

func NewPanicRecoverer(logger logger) *PanicRecoverer {
	return &PanicRecoverer{
		logger: logger,
	}
}

func (m *PanicRecoverer) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer m.recoverAndRespond(w)
			next.ServeHTTP(w, r)
		})
	}
}

func (m *PanicRecoverer) recoverAndRespond(w http.ResponseWriter) {
	rec := recover()
	if rec == nil {
		return
	}
	if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
		return
	}
	m.logger.Error("Panic while serving HTTP request.", log.Any("details", rec))
	_ = errs.NewError(errs.ErrorCodeInternal, "").RenderToHTTPResponse(w)
}
