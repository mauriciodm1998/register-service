package mail

import (
	"bytes"
	"html/template"
	"net/smtp"
)

type Mailer interface {
	Send(subject, body string, to []string) error
	MountHTMLBody(params any) (string, error)
}

type mailer struct {
	from        string
	password    string
	smtpAddress string
}

func NewMailer(from, password, smtpAddress string) Mailer {
	return &mailer{
		from:        from,
		password:    password,
		smtpAddress: smtpAddress,
	}
}

func (m *mailer) Send(subject, body string, to []string) error {
	smtp.PlainAuth("", m.from, m.password, m.smtpAddress)

	err := smtp.SendMail(m.smtpAddress, nil, m.from, to, []byte(body))
	if err != nil {
		return err
	}
	return nil
}

func (m *mailer) MountHTMLBody(params any) (string, error) {

	t, err := template.ParseFiles("template.html")
	if err != nil {
		return "", err
	}
	var body bytes.Buffer
	err = t.Execute(&body, params)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}
