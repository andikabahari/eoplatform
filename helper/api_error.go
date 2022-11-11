package helper

type APIError interface {
	Error() string
	APIError() (int, string)
}

type apiError struct {
	status  int
	message string
}

func NewAPIError(status int, message string) APIError {
	return &apiError{status, message}
}

func (e apiError) Error() string {
	return e.message
}

func (e apiError) APIError() (int, string) {
	return e.status, e.message
}
