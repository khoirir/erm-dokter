package domain

type BusinessError struct {
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

func NewBusinessError(msg string) error {
	return &BusinessError{Message: msg}
}
