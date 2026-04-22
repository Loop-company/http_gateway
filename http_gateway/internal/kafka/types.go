package kafka

type LoginCommand struct {
	RequestID string `json:"request_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginResponse struct {
	RequestID    string `json:"request_id"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Error        string `json:"error,omitempty"`
}
