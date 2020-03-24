package smtp

import (
	"fmt"
	"net/smtp"
	"text/template"
)

const (
	_smtp = "smtp.gmail.com"
	_port = "587"
)

type Mailer interface {
	SendEmail(Mail) (bool, error)
}

type Mail struct {
	From     string
	To       []string
	Subject  string
	Template *template.Template
	Body     []byte
	Username string
	Password string
}

func (m Mail) SendEmail(env string) (bool, error) {
	if err := smtp.SendMail(fmt.Sprintf("%s:%s", _smtp, _port), smtp.PlainAuth("", m.Username, m.Password, _smtp), m.From, m.To, m.Body); err != nil {
		return false, err
	}

	return false, nil
}
