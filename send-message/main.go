package handler

import (
	"fmt"
	"net/smtp"
	"os"

	ort "github.com/open-runtimes/types-for-go/v4/openruntimes"
)

type ResponseBody struct {
	Message string `json:"message"`
}

type RequestBody struct {
	Name     string `json:"name"`
	Subject  string `json:"subject"`
	Contact  string `json:"contact"`
	Email    string `json:"email"`
	WhatsApp string `json:"whatsapp"`
	Message  string `json:"message"`
}

func Main(Context ort.Context) ort.Response {
	if Context.Req.Path == "/" && Context.Req.Method == "POST" {
		var reqBody RequestBody
		err := Context.Req.BodyJson(&reqBody)
		if err != nil {
			Context.Res.WithStatusCode(400)
			return Context.Res.Json(ResponseBody{Message: "Bad request"})
		} else {
			emailFrom := os.Getenv("APPWRITE_FUNCTION_SEND_MESSAGE_EMAIL_FROM")
			emailTo := os.Getenv("APPWRITE_FUNCTION_SEND_MESSAGE_EMAIL_TO")
			emailSecret := os.Getenv("APPWRITE_FUNCTION_SEND_MESSAGE_EMAIL_SECRET")
			smtpHost := os.Getenv("APPWRITE_FUNCTION_SEND_MESSAGE_SMTP_SERVER_HOST")
			smtpPort := os.Getenv("APPWRITE_FUNCTION_SEND_MESSAGE_SMTP_SERVER_PORT")

			auth := smtp.PlainAuth("", emailFrom, emailSecret, smtpHost)

			body := fmt.Sprintf("Name: %v\nContact: %v\nEmail: %v\nWhatsApp: %v\n\n%v",
				reqBody.Name, reqBody.Contact, reqBody.Email, reqBody.WhatsApp, reqBody.Message)
			msg := fmt.Sprintf("From: %v\nTo: %v\nSubject: %v\n\n%v", emailFrom, emailTo, reqBody.Subject, body)

			err := smtp.SendMail(fmt.Sprintf("%v:%v", smtpHost, smtpPort), auth, emailFrom, []string{emailTo}, []byte(msg))
			if err != nil {
				Context.Error(err)
			}

			Context.Res.WithStatusCode(201)
			return Context.Res.Json(ResponseBody{Message: "Created"})

		}

	} else {
		Context.Res.WithStatusCode(400)
		return Context.Res.Json(ResponseBody{Message: "Bad request"})

	}
}
