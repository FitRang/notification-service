package apperror

type Code string

const (
    Unauthenticated Code = "UNAUTHENTICATED"
    Forbidden       Code = "FORBIDDEN"
    NotFound        Code = "NOT_FOUND"
    Conflict        Code = "CONFLICT"
    BadInput        Code = "BAD_USER_INPUT"
    Internal        Code = "INTERNAL_SERVER_ERROR"
)

type AppError struct {
    Code    Code
    Message string
    Err     error
}

func (e *AppError) Error() string {
    return e.Message
}

func New(code Code, msg string) *AppError {
    return &AppError{Code: code, Message: msg}
}

func Wrap(code Code, msg string, err error) *AppError {
    return &AppError{Code: code, Message: msg, Err: err}
}
