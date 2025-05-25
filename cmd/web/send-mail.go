package main

import (
	"log"
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

	emial.SetBody(mail.TextHTML, m.Content)

	err = emial.Send(client)
	if err != nil {
		log.Println("Error sending email:", err)
	}

	log.Println("Email sent")
}

