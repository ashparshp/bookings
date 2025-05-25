package main

import (
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
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10*time.Second
	server.SendTimeout = 10*time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println("Error connecting to mail server:", err)
	}

	emial := mail.NewMSG()
	emial.SetFrom(m.From).
		AddTo(m.To).
		SetSubject(m.Subject)

	if m.Template == "" {
	emial.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println("Error reading template file:", err)
			return
		}
		
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		emial.SetBody(mail.TextHTML, msgToSend)
	}

	err = emial.Send(client)
	if err != nil {
		log.Println("Error sending email:", err)
	}

	log.Println("Email sent")
}

