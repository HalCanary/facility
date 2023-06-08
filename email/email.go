// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.
package email

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/http"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/HalCanary/facility/humanize"
)

type Address = mail.Address

// Data structure representing instructions for connecting to SMTP server.
// Headers are additional headers to be added to outgoing email.
type EmailSecrets struct {
	SmtpHost string            // example: "smtp.gmail.com"
	SmtpUser string            // example: "foobar@gmail.com"
	SmtpPass string            // for gmail, is a App Password
	From     Address           // example: {"Foo Bar", "foobar@gmail.com"}
	Headers  map[string]string // extra headers to be added to email.
}

// Attachment for an email.
type Attachment struct {
	Filename    string // Optional
	ContentType string // If empty, determined via http.DetectContentType
	Data        []byte
	Textual     bool // If true and Data is valid UTF-8, then encode as quoted-printable over base64
}

// An electric mail message.
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

// Read email secrets from the given file.
func GetSecrets(path string) (EmailSecrets, error) {
	var v EmailSecrets
	b, err := os.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(b, &v)
	}
	return v, err
}

// Send the given email using the provided SMTP secrets.
func (m Email) Send(secrets EmailSecrets) error {
	to := make([]string, 0, len(m.To)+len(m.Cc)+len(m.Bcc))
	for _, list := range [...][]Address{m.To, m.Cc, m.Bcc} {
		for _, a := range list {
			to = append(to, a.Address)
		}
	}
	if len(secrets.Headers) > 0 {
		old := m.Headers
		m.Headers = make(map[string]string, len(old)+len(secrets.Headers))
		for _, d := range [...]map[string]string{old, secrets.Headers} {
			for k, v := range d {
				m.Headers[k] = v
			}
		}
	}
	var buffer bytes.Buffer
	m.Make(&buffer)
	msg := buffer.Bytes()
	auth := smtp.PlainAuth("", secrets.SmtpUser, secrets.SmtpPass, secrets.SmtpHost)
	return smtp.SendMail(secrets.SmtpHost+":587", auth, secrets.SmtpUser, to, msg)
}

func qencode(out io.Writer, s string) {
	//out.Write([]byte(mime.QEncoding.Encode("utf-8", s)))
	for s != "" {
		idx := strings.Index(s, " ") + 1
		if idx <= 0 {
			idx = len(s)
		}
		out.Write([]byte(mime.QEncoding.Encode("utf-8", s[:idx])))
		s = s[idx:]
	}
}

var (
	space      = [...]byte{' '}
	comma      = [...]byte{','}
	colon      = [...]byte{':'}
	colonspace = [...]byte{':', ' '}
	crlf       = [...]byte{'\r', '\n'}
)

func encodeHeader(out io.Writer, key, s string) {
	if s == "" {
		return
	}
	io.WriteString(out, textproto.CanonicalMIMEHeaderKey(key))
	out.Write(colonspace[:])
	qencode(out, s)
	out.Write(crlf[:])
}

func encodeMultiheader(out io.Writer, key string, values []Address) {
	if len(values) > 0 {
		io.WriteString(out, textproto.CanonicalMIMEHeaderKey(key))
		out.Write(colon[:])
		for i, val := range values {
			out.Write(space[:])
			io.WriteString(out, val.String())
			if i+1 != len(values) {
				out.Write(comma[:])
			}
			out.Write(crlf[:])
		}
	}
}

func (mail Email) makeHeader(out io.Writer) {
	if mail.Date.IsZero() {
		mail.Date = time.Now()
	}
	encodeHeader(out, "Date", mail.Date.Format(time.RFC1123Z))
	encodeHeader(out, "Subject", mail.Subject)
	encodeMultiheader(out, "From", []Address{mail.From})
	encodeMultiheader(out, "To", mail.To)
	encodeMultiheader(out, "Cc", mail.Cc)
	for key, value := range mail.Headers {
		encodeHeader(out, key, value)
	}
}

// Make, but do not send an email message.
func (mail Email) Make(out io.Writer) {
	const boundary = "================"
	mail.makeHeader(out)
	encodeHeader(out, "MIME-Version", "1.0")

	if len(mail.Attachments) == 0 {
		encodeHeader(out, "Content-Type", `text/plain; charset="UTF-8"`)
		encodeHeader(out, "Content-Transfer-Encoding", "quoted-printable")
		out.Write(crlf[:]) // end of header
		quotedprintableWrite([]byte(mail.Content), out)
		out.Write(crlf[:])
		return
	}

	encodeHeader(out, "Content-Type", `multipart/mixed; boundary="`+boundary+`"`)
	out.Write(crlf[:]) // end of header

	mw := multipart.NewWriter(out)
	mw.SetBoundary(boundary)
	if mail.Content != "" {
		w, _ := mw.CreatePart(textproto.MIMEHeader{
			"Content-Type":              []string{`text/plain; charset="UTF-8"`},
			"Content-Transfer-Encoding": []string{"quoted-printable"},
		})
		quotedprintableWrite([]byte(mail.Content), w)
	}
	for _, attachment := range mail.Attachments {
		contentType := attachment.ContentType
		if contentType == "" {
			contentType = http.DetectContentType(attachment.Data)
		}
		if attachment.Textual && utf8.Valid(attachment.Data) {
			w, _ := mw.CreatePart(textproto.MIMEHeader{
				"Content-Type":              []string{contentType},
				"Content-Transfer-Encoding": []string{"quoted-printable"},
				"Content-Disposition":       []string{contentDisposition(attachment.Filename)},
				"MIME-Version":              []string{"1.0"},
			})
			quotedprintableWrite(attachment.Data, w)
			continue
		}
		w, _ := mw.CreatePart(textproto.MIMEHeader{
			"Content-Type":              []string{contentType},
			"Content-Transfer-Encoding": []string{"base64"},
			"Content-Disposition":       []string{contentDisposition(attachment.Filename)},
			"MIME-Version":              []string{"1.0"},
		})
		base64Write(attachment.Data, w)
	}
	mw.Close()
}

// Send a file to a single destination.
func SendFile(dst Address, path, contentType string, secrets EmailSecrets) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	base := filepath.Base(path)
	subject := fmt.Sprintf("(%s) %s", humanize.Humanize(int64(len(data))), base)
	return Email{
		From:    secrets.From,
		To:      []Address{dst},
		Subject: subject,
		Content: "☺",
		Attachments: []Attachment{
			Attachment{
				Data:        data,
				ContentType: contentType,
				Filename:    base,
			},
		},
	}.Send(secrets)
}

func contentDisposition(filename string) string {
	if filename == "" {
		return "attachment"
	}
	return fmt.Sprintf("attachment; filename=%q", mime.QEncoding.Encode("utf-8", filename))
}

func quotedprintableWrite(src []byte, dst io.Writer) {
	qpw := quotedprintable.NewWriter(dst)
	qpw.Write(src)
	qpw.Close()
}

func base64Write(src []byte, dst io.Writer) {
	const linelength = 57
	const bufferlength = 78 // base64.StdEncoding.EncodedLen(57) + 2
	var buffer [bufferlength]byte
	for len(src) > 0 {
		l := len(src)
		if l > linelength {
			l = linelength
		}
		el := base64.StdEncoding.EncodedLen(l)
		base64.StdEncoding.Encode(buffer[:el], src[:l])
		src = src[l:]
		copy(buffer[el:el+2], crlf[:])
		dst.Write(buffer[:el+2])
	}
}
