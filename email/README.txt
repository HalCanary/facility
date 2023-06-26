package email // import "github.com/HalCanary/facility/email"

Copyright 2022 Hal Canary Use of this program is governed by the file LICENSE.

FUNCTIONS

func SendFile(dst Address, path, contentType string, secrets EmailSecrets) error
    Send a file to a single destination.


TYPES

type Address = mail.Address

type Attachment struct {
	Filename    string // Optional
	ContentType string // If empty, determined via http.DetectContentType
	Data        []byte
	Textual     bool // If true and Data is valid UTF-8, then encode as quoted-printable over base64
}
    Attachment for an email.

type Email struct {
	Date        time.Time // If not set, use time.Now()
	To          []Address
	Cc          []Address
	Bcc         []Address
	From        Address
	Subject     string
	Content     string // Assumed to be text/plain.
	Attachments []Attachment
	Headers     map[string]string // Optional extra headers.
}
    An electric mail message.

func (mail Email) Make(out io.Writer)
    Make, but do not send an email message.

func (m Email) Send(secrets EmailSecrets) error
    Send the given email using the provided SMTP secrets.

type EmailSecrets struct {
	SmtpHost string            // example: "smtp.gmail.com"
	SmtpUser string            // example: "foobar@gmail.com"
	SmtpPass string            // for gmail, is a App Password
	From     Address           // example: {"Foo Bar", "foobar@gmail.com"}
	Headers  map[string]string // extra headers to be added to email.
}
    Data structure representing instructions for connecting to SMTP server.
    Headers are additional headers to be added to outgoing email.

func GetSecrets(path string) (EmailSecrets, error)
    Read email secrets from the given JSON file. It might look something like
    this:

        {
            "SmtpHost": "",
            "SmtpUser": "",
            "SmtpPass": "",
            "FromAddr": "Foo Bar <foobar@example.com>",
            "From": {
                "Name": "Foo Bar",
                "Address": "foobar@example.com"
            },
            "Headers": {
                "X-PGP-Key": "",
                "X-PGP-KeyID": ""
            }
        }

