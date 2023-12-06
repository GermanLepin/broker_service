package dto

type RequestPayload struct {
	Action                string                `json:"action"`
	AuthenticationPayload AuthenticationPayload `json:"authentication,omitempty"`
	LoggingPayload        LoggingPayload        `json:"logging,omitempty"`
	MailPayload           MailPayload           `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthenticationPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoggingPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
