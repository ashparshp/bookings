package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ashparshp/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
    server.Host = app.MailConfig.Host
    server.Port = app.MailConfig.Port
    server.Username = app.MailConfig.Username
    server.Password = app.MailConfig.Password
    server.KeepAlive = false
    server.ConnectTimeout = 10 * time.Second
    server.SendTimeout = 10 * time.Second

	switch strings.ToLower(app.MailConfig.Encryption) {
	case "starttls", "tls":
		server.Encryption = mail.EncryptionSTARTTLS
	case "ssl":
		server.Encryption = mail.EncryptionSSLTLS
	case "none":
		server.Encryption = mail.EncryptionNone
	default:
		log.Printf("Unknown mail encryption type: %s. Defaulting to none.", app.MailConfig.Encryption)
		server.Encryption = mail.EncryptionNone
	}

	if app.MailConfig.Host == "smtp.gmail.com" && app.MailConfig.Port == 465 {
        server.Encryption = mail.EncryptionSSLTLS
        server.TLSConfig = &tls.Config{InsecureSkipVerify: false}
    }

	if app.MailConfig.Username != "" && app.MailConfig.Password != "" {
        server.Authentication = mail.AuthPlain
    }

	client, err := server.Connect()
	if err != nil {
		errorLog.Println("Error connecting to mail server:", err)
	}

	email := mail.NewMSG()
	fromAddress := app.MailConfig.FromAddress
    if m.From != "" {
        fromAddress = m.From
    }
	
	email.SetFrom(fromAddress).
        AddTo(m.To).
        SetSubject(m.Subject)

    if m.Template == "" {
        email.SetBody(mail.TextHTML, m.Content)
    } else {
        data, err := os.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
        if err != nil {
            app.ErrorLog.Println("Error reading template file:", err)
            return
        }
        
        mailTemplate := string(data)
        msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
        email.SetBody(mail.TextHTML, msgToSend)
    }

    err = email.Send(client)
    if err != nil {
        log.Println("Error sending email:", err)
        return
    }

	log.Println("Email sent")
}
