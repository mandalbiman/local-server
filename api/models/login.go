package models

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Roles *[]string `json:"roles,omitempty"`
	Token *string   `json:"token,omitempty"`
	User  *string   `json:"user,omitempty"`
}

type N400BadResponse struct {
	Error   *string `json:"error,omitempty"`
	Message *string `json:"message,omitempty"`
}
