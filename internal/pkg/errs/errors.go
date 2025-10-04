package errs

import (
	"cmp"
	"encoding/json"
	"errors"
	"reflect"
	"slices"

	"go.uber.org/zap/zapcore"
)

type ErrorCode uint

const (
	ErrorCodeOK ErrorCode = iota
	ErrorCodeCanceled
	ErrorCodeUnknown
	ErrorCodeInvalidArgument
	ErrorCodeDeadlineExceeded
	ErrorCodeNotFound
	ErrorCodeAlreadyExists
	ErrorCodePermissionDenied
	ErrorCodeResourceExhausted
	ErrorCodeFailedPrecondition
	ErrorCodeAborted
	ErrorCodeOutOfRange
	ErrorCodeUnimplemented
	ErrorCodeInternal
	ErrorCodeUnavailable
	ErrorCodeDataLoss
	ErrorCodeUnauthenticated
	ErrorCodeClosedRequest
)

type Param struct {
	Key   string `description:"Ключ параметра"                      json:"key"`
	Value string `description:"Значение параметра"                  json:"value"`
} // @name Param

type Params []Param // @name Params

func (p Param) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString(p.Key, p.Value)
	return nil
}

func (p Params) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	for _, param := range p {
		encoder.AddString(param.Key, param.Value)
	}

	return nil
}

type Error struct {
	Code    ErrorCode `description:"Код ошибки (3=Неверный аргумент, 5=Не найдено, и т.д.)"                                      json:"code"`
	Message string    `description:"Читаемое сообщение об ошибке"                                                                json:"message"`
	Params  Params    `description:"Дополнительные параметры ошибки"                                                             json:"params"`
	Err     error     `json:"-"` // Internal error, not exposed in API
} // @name Error

func (e *Error) Error() string {
	data, err := json.Marshal(e)
	if err != nil {
		return "Error: " + e.Message
	}
	return string(data)
}

func (e *Error) AddParam(key, value string) {
	e.Params = append(e.Params, Param{Key: key, Value: value})
}

func NewError(code ErrorCode, message string) *Error {
	return &Error{Code: code, Message: message, Params: nil, Err: nil}
}

func (e *Error) WithParam(key, value string) *Error {
	e.AddParam(key, value)
	return e
}

func (e *Error) WithCause(err error) *Error {
	e.Err = err
	return e
}

func (e *Error) WithParams(params ...Param) *Error {
	if len(e.Params) == 0 {
		e.Params = params
	} else {
		e.Params = append(e.Params, params...)
	}
	return e
}

func (e *Error) Cause() error {
	return e.Err
}

func (e *Error) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("message", e.Message)
	encoder.AddUint("code", uint(e.Code))
	if err := encoder.AddObject("params", e.Params); err != nil {
		return err
	}
	return nil
}

func NewUnexpectedBehaviourError(details string) *Error {
	err := NewError(ErrorCodeInternal, "Unexpected behavior.")
	err.AddParam("details", details)
	return err
}

func NewInvalidFormError() *Error {
	return NewError(
		ErrorCodeInvalidArgument,
		"The request data is invalid.",
	)
}

func NewInvalidParameter(message string) *Error {
	return NewError(ErrorCodeInvalidArgument, message)
}

func NewEntityNotFoundError() *Error {
	return NewError(ErrorCodeNotFound, "Entity not found.")
}

func NewBadTokenError() *Error {
	return NewError(ErrorCodePermissionDenied, "Bad token.")
}

func NewPermissionDeniedError() *Error {
	return NewError(ErrorCodePermissionDenied, "Permission denied.")
}

func NewUnauthenticatedError() *Error {
	return NewError(ErrorCodeFailedPrecondition, "Unauthenticated error.")
}

func NewEntityAlreadyExistsError() *Error {
	return NewError(ErrorCodeAlreadyExists, "Entity already exists.")
}

func NewJSONDeserializationError(err error) *Error {
	return NewError(ErrorCodeInvalidArgument, "Failed to deserialize JSON.").WithCause(err)
}

func NewJSONSerializationError(err error) *Error {
	return NewError(ErrorCodeInternal, "Failed to serialize JSON.").WithCause(err)
}

func NewBase64EncodingError(err error) *Error {
	return NewError(ErrorCodeInternal, "Failed to encode base64.").WithCause(err)
}

func NewBase64DecodingError(err error) *Error {
	return NewError(ErrorCodeInvalidArgument, "Failed to decode base64.").WithCause(err)
}

func NewInvalidCursorError(message string) *Error {
	return NewError(ErrorCodeInvalidArgument, message)
}

func paramCompare(a, b Param) int {
	return cmp.Compare(a.Key, b.Key)
}

func (e *Error) Is(tgt error) bool {
	var target *Error
	if ok := errors.As(tgt, &target); !ok {
		return false
	}
	target.Err = nil

	slices.SortFunc(target.Params, paramCompare)
	err := *e
	err.Err = nil
	slices.SortFunc(err.Params, paramCompare)
	eq := reflect.DeepEqual(&err, target)
	return eq
}
