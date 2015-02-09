package passwordless

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/smtp"
	"os"
	"strings"
	"text/template"
)

var (
	auth      smtp.Auth
	authParts map[string]string
	fromEmail string
)

var emailTemplate = template.Must(template.New("emailTemplate").Parse(`From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}`))

func SendMail(to []string, subject, message string) error {
	var doc bytes.Buffer

	ctx := struct {
		From    string
		To      string
		Subject string
		Body    string
	}{
		fromEmail,
		strings.Join(to, ", "),
		subject,
		message,
	}

	if err := emailTemplate.Execute(&doc, ctx); err != nil {
		return err
	}

	return smtp.SendMail(
		fmt.Sprintf("%v:%v", authParts["EMAIL_HOST"], authParts["EMAIL_PORT"]),
		auth,
		fromEmail,
		to,
		doc.Bytes())
}

func init() {
	authParts = make(map[string]string)
	// get from address
	fromEmail = os.Getenv("EMAIL_FROM")
	if fromEmail == "" {
		log.Panic("Missing required environment variable 'EMAIL_FROM'.")
	}

	// initialize smtp auth
	for _, part := range []string{"EMAIL_HOST_USER", "EMAIL_HOST_PASSWORD", "EMAIL_HOST", "EMAIL_PORT"} {
		envPart := os.Getenv(part)
		if envPart == "" {
			log.Panicf("Missing required environment variable '%s'.", part)
		}
		authParts[part] = envPart
	}
	auth = smtp.PlainAuth("", authParts["EMAIL_HOST_USER"], authParts["EMAIL_HOST_PASSWORD"], authParts["EMAIL_HOST"])
}
