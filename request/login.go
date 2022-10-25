package request

type BasicLogin struct {
	Email    string
	Password string
}

type LoginRequest struct {
	BasicLogin
}
