package mail_service

import (
	"bytes"
	"net/http"
	"text/template"
	"time"

	"mail-service/internal/dto"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type JsonService interface {
	ErrorJSON(w http.ResponseWriter, err error, status ...int) error
	WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error
	ReadJSON(w http.ResponseWriter, r *http.Request, data any) error
}

func (s *service) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := s.jsonService.ReadJSON(w, r, &requestPayload)
	if err != nil {
		s.jsonService.ErrorJSON(w, err)
		return
	}

	msg := dto.Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = s.SendSMTPMessage(msg)
	if err != nil {
		s.jsonService.ErrorJSON(w, err)
		return
	}

	payload := dto.JsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	s.jsonService.WriteJSON(w, http.StatusAccepted, payload)
}

func (s *service) SendSMTPMessage(msg dto.Message) error {
	if msg.From == "" {
		msg.From = s.mail.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = s.mail.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := s.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := s.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = s.mail.Host
	server.Port = s.mail.Port
	server.Username = s.mail.Username
	server.Password = s.mail.Password
	server.Encryption = s.getEncryption(s.mail.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) buildHTMLMessage(msg dto.Message) (string, error) {
	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = s.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (s *service) buildPlainTextMessage(msg dto.Message) (string, error) {
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (s *service) inlineCSS(str string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(str, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (s *service) getEncryption(str string) mail.Encryption {
	switch str {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

func New(
	jsonService JsonService,
	mail dto.Mail,
) *service {
	return &service{
		jsonService: jsonService,
		mail:        mail,
	}
}

type service struct {
	jsonService JsonService
	mail        dto.Mail
}
