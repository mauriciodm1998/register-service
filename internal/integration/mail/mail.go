package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"register-service/internal/config"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type Mailer interface {
	SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error
	MountHTMLBody(params any) (string, error)
}

type mailer struct {
	name        string
	from        string
	password    string
	smtpAddress string
}

func NewMailer() Mailer {
	return &mailer{
		name:        "Hackaton FIAP",
		from:        config.Get().Mailer.From,
		password:    config.Get().Mailer.Pwd,
		smtpAddress: smtpAuthAddress,
	}
}

func (sender *mailer) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.from)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	return e.Send(smtpServerAddress, smtp.PlainAuth("", sender.from, sender.password, smtpAuthAddress))
}

func (m *mailer) MountHTMLBody(params any) (string, error) {
	tmpl, err := template.New("template").Parse(`<!-- template.html -->
	<!DOCTYPE html>
	<html>
	<body>
		<h3>User: </h3><span>{{.ID}}</span><br/><br/>
		<h3>Total time: </h3><span>{{.Time}}</span><br/><br/>
		<table border=1>
			<tr>
			  <th>Date</th>
			  <th>Clock Ins</th>
			  <th>Time</th>
		  </tr>
			{{.Message}}
		</table>
	</body>
	</html>`)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer

	err = tmpl.Execute(&body, params)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}
