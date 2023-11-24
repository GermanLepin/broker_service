package dto

type RequestPayload struct {
	Action      string                `json:"action"`
	AuthPayload AuthenticationPayload `json:"authentication,omitempty"`
}

type AuthenticationPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
