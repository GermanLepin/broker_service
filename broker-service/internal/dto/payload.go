package dto

type RequestPayload struct {
	Action                string                `json:"action"`
	AuthenticationPayload AuthenticationPayload `json:"authentication,omitempty"`
	LoggingPayload        LoggingPayload        `json:"logging,omitempty"`
}

type AuthenticationPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoggingPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
