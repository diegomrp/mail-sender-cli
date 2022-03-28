package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/diegomrp/mail-sender-cli/struts"
	mail "github.com/xhit/go-simple-mail/v2"
)

func main() {

	var mailVars = struts.MailVars{}
	var r struts.Recipient
	var attFiles []mail.File

	flag.StringVar(&mailVars.Host, "host", "mail.central.inditex.grp", "SMTP host")
	flag.IntVar(&mailVars.Port, "port", 25, "SMTP port")
	flag.StringVar(&mailVars.Template, "template", "./resources/template/template.html", "Email template to send")
	flag.StringVar(&mailVars.PathImage, "image", "./resources/imgs", "Image directory that the program will use to attach the images to the email")
	flag.StringVar(&mailVars.PathDocument, "document", "./resources/documents", "Document directory that the program will use to attach the docuents to the email")
	flag.StringVar(&mailVars.From, "from", "Test <noreply@inditex.com>", "From email user")
	flag.StringVar(&mailVars.Subject, "subject", "Test", "Email Subject")
	flag.StringVar(&mailVars.RecipientsFile, "recipients", "./resources/csv/recipients.csv", "File with the recipients of the email")
	flag.Parse()

	server := mail.NewSMTPClient()

	server.Host = mailVars.Host
	server.Port = mailVars.Port
	server.Encryption = mail.EncryptionSTARTTLS

	// Variable to keep alive connection
	server.KeepAlive = true

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// SMTP client
	smtpClient, err := server.Connect()

	if err != nil {
		log.Fatal(err)
	}

	htmlBody := readTemplate(mailVars.Template)

	if htmlBody != "" {
		recipients, err := r.GetRecipients(mailVars.RecipientsFile)

		if err != nil {
			log.Fatal("Error reading recipients")
		}

		if mailVars.PathImage != "" {
			attachFiles(mailVars.PathImage, true, &attFiles)
		}

		if mailVars.PathDocument != "" {
			attachFiles(mailVars.PathDocument, false, &attFiles)
		}

		for _, to := range recipients {
			// New email simple html with inline
			email := mail.NewMSG()

			email.SetFrom(mailVars.From).
				AddTo(to.Email).
				SetSubject(mailVars.Subject)

			email.SetBody(mail.TextHTML, htmlBody)

			// Attach the common images and documents
			for _, attf := range attFiles {
				email.Attach(&attf)
			}

			// always check error after send
			if email.Error != nil {
				log.Fatal(email.Error)
			}

			// Call Send and pass the client
			err = email.Send(smtpClient)
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Email Sent to", to.Email)
			}
		}
	} else {
		println("No content to send")
	}
}

func readTemplate(path string) string {
	content, err := ioutil.ReadFile(path)

	if err == nil {
		return string(content)
	}
	return ""
}

func attachFiles(path string, inline bool, attFiles *[]mail.File) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Println(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			*attFiles = append(*attFiles, mail.File{FilePath: path + "/" + file.Name(), Name: file.Name(), Inline: inline})
		}
	}
}
