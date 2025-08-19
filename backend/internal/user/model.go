package user

type RegistrationRequest struct {
	Email    string
	Username string
	Password string
}

type LoginRequest struct {
	Username string
	Password string
}

type Response struct {
	Result string `json:"result"`
}
