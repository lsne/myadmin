/*
 * @Author: Liu Sainan
 * @Date: 2024-01-15 23:44:17
 */

package mailutils

import "gopkg.in/gomail.v2"

type Mail struct {
	sender   string
	receiver []string
	// Cc         []string
	Subject    string
	bodyFormat string // "text/html" or "text/plain" or other
	Body       string
	SMTPServer string
	SMTPPort   uint16
}

func NewMail(sender string, receiver []string, bodyFormat string, SMTPServer string, SMTPPort uint16) *Mail {
	if bodyFormat == "" {
		bodyFormat = "text/plain"
	}
	return &Mail{
		sender:     sender,
		receiver:   receiver,
		bodyFormat: bodyFormat,
		SMTPServer: SMTPServer,
		SMTPPort:   SMTPPort,
	}
}

func (m *Mail) SendMail(subject string) (err error) {
	mail := gomail.NewMessage()
	mail.SetHeader("From", m.sender)
	mail.SetHeader("TO", m.receiver...)
	// m.SetHeader("Cc", c.Cc...)
	mail.SetHeader("Subject", subject)
	mail.SetBody(m.bodyFormat, m.Body)

	d := gomail.Dialer{Host: m.SMTPServer, Port: int(m.SMTPPort)}

	return d.DialAndSend(mail)
}
