package exception

type InternalServerError struct {
	Message string
}

func (e InternalServerError) Error() string {
	return e.Message
}

func NewInternalServerError(message string) InternalServerError {
	return InternalServerError{Message: message}
}