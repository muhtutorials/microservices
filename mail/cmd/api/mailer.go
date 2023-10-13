package main

import (
	"bytes"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"time"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	msg.From = m.FromAddress
	msg.FromName = m.FromName
	msg.DataMap = map[string]any{
		"message": msg.Data,
	}

	htmlMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainTextMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.
		SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject).
		SetBody(mail.TextPlain, plainTextMessage).
		AddAlternative(mail.TextHTML, htmlMessage)

	if len(msg.Attachments) > 0 {
		for _, att := range msg.Attachments {
			email.AddAttachment(att)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	filename := "./templates/mail.html.gohtml"

	t, err := template.New("email_html").ParseFiles(filename)
	if err != nil {
		return "", err
	}

	var tmpl bytes.Buffer
	// body is "{{define "body"}}" in the template
	if err = t.ExecuteTemplate(&tmpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	htmlMessage := tmpl.String()
	htmlMessage, err = m.inlineCSS(htmlMessage)
	if err != nil {
		return "", err
	}

	return htmlMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	filename := "./templates/mail.txt.gohtml"

	t, err := template.New("email_plain").ParseFiles(filename)
	if err != nil {
		return "", err
	}

	var tmpl bytes.Buffer
	if err = t.ExecuteTemplate(&tmpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainTextMessage := tmpl.String()

	return plainTextMessage, nil
}

func (m *Mail) inlineCSS(str string) (string, error) {
	opts := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(str, &opts)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncryption(str string) mail.Encryption {
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
