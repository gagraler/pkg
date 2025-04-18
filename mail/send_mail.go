package mail

import (
	"github.com/gagraler/pkg/log"
	"gopkg.in/gomail.v2"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/6/15 0:21
 * @file: send_mail.go
 * @description: 发送邮件
 */

var dFlag bool

type Mail struct {
	From     string
	To       string
	Subject  string
	Body     string
	Host     string
	User     string
	Password string
	Port     int
}

// SendMail SendMail 发送邮件
func (m *Mail) SendMail() {
	mail := gomail.NewMessage()
	mail.SetHeader("From", m.From)
	mail.SetHeader("To", m.To)
	mail.SetHeader("Subject", m.Subject)
	mail.SetBody("text/html", m.Body)

	d := gomail.NewDialer(m.Host, m.Port, m.User, m.Password)
	if err := d.DialAndSend(mail); err != nil {
		log.Errorf("send email error: %s", err.Error())
	}

	if dFlag {
		log.Debug("sending email...")
	}

	if dFlag {
		log.Debug("send email success!")
	}
	return
}

// SendAttachmentMail SendAttachmentMail 发送带附件的邮件
func (m *Mail) SendAttachmentMail(filename string) {
	mail := gomail.NewMessage()
	mail.SetHeader("From", m.From)
	mail.SetHeader("To", m.To)
	mail.SetHeader("Subject", m.Subject)
	mail.SetBody("text/html", m.Body)

	mail.Attach(filename)

	d := gomail.NewDialer(m.Host, m.Port, m.User, m.Password)
	if err := d.DialAndSend(mail); err != nil {
		log.Errorf("send email error: %s", err.Error())
	}

	if dFlag {
		log.Debug("sending email...")
	}

	if err := d.DialAndSend(mail); err != nil {
		if dFlag {
			log.Errorf("send email error: %s", err.Error())
		}
		return
	}

	if dFlag {
		log.Debug("send email success!")
	}
}
