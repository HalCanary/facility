package email

// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.

import (
	"bytes"
	"io"
	"mime"
	"mime/quotedprintable"
	"net/mail"
	"strings"
	"testing"
	"time"

	"github.com/HalCanary/facility/expect"
)

const testdata = `Lorem ipsum dolor sit amet, consectetur adipiscing elit.  Nunc imperdiet elit
eu sapien accumsan, quis accumsan nisl porta. Donec congue dui in dignissim
tincidunt. Phasellus vel ligula lobortis tortor iaculis vulputate. Proin non
augue quis est molestie dignissim eget sit amet enim.  Donec ut purus ac enim
hendrerit ornare. Pellentesque egestas tempor sodales.  Pellentesque eget
auctor mauris.`

const expected = `Date: Sat, 01 Jan 2022 00:00:00 +0000
Subject: a quick note =?utf-8?q?(=E2=99=A0=E2=99=A5=E2=99=A6=E2=99=A3)?=
From: =?utf-8?q?Z_=E2=86=90=E2=86=91=E2=86=92=E2=86=93?= <z@example.com>
To: "A" <a@example.com>,
 "B" <b@example.com>
Cc: "C" <c@example.com>,
 "D" <d@example.com>
Mime-Version: 1.0
Content-Type: multipart/mixed; boundary="================"

--================
Content-Transfer-Encoding: quoted-printable
Content-Type: text/plain; charset="UTF-8"

Hello, World!
--================
Content-Disposition: attachment; filename="foo.txt"
Content-Transfer-Encoding: base64
Content-Type: text/plain; charset=utf-8
MIME-Version: 1.0

TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4g
IE51bmMgaW1wZXJkaWV0IGVsaXQKZXUgc2FwaWVuIGFjY3Vtc2FuLCBxdWlzIGFjY3Vtc2FuIG5p
c2wgcG9ydGEuIERvbmVjIGNvbmd1ZSBkdWkgaW4gZGlnbmlzc2ltCnRpbmNpZHVudC4gUGhhc2Vs
bHVzIHZlbCBsaWd1bGEgbG9ib3J0aXMgdG9ydG9yIGlhY3VsaXMgdnVscHV0YXRlLiBQcm9pbiBu
b24KYXVndWUgcXVpcyBlc3QgbW9sZXN0aWUgZGlnbmlzc2ltIGVnZXQgc2l0IGFtZXQgZW5pbS4g
IERvbmVjIHV0IHB1cnVzIGFjIGVuaW0KaGVuZHJlcml0IG9ybmFyZS4gUGVsbGVudGVzcXVlIGVn
ZXN0YXMgdGVtcG9yIHNvZGFsZXMuICBQZWxsZW50ZXNxdWUgZWdldAphdWN0b3IgbWF1cmlzLg==

--================--
`

func toQuotedPrintable(s string) string {
	var b bytes.Buffer
	qpw := quotedprintable.NewWriter(&b)
	io.WriteString(qpw, s)
	qpw.Close()
	return b.String()
}

func TestEmail(t *testing.T) {
	mail := Email{
		Date:    time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:      []Address{Address{"A", "a@example.com"}, Address{"B", "b@example.com"}},
		Cc:      []Address{Address{"C", "c@example.com"}, Address{"D", "d@example.com"}},
		Bcc:     []Address{Address{"E", "e@example.com"}, Address{"F", "f@example.com"}},
		From:    Address{"Z ←↑→↓", "z@example.com"},
		Subject: "a quick note (♠♥♦♣)",
		Content: "Hello, World!",
		Attachments: []Attachment{
			Attachment{
				Filename: "foo.txt",
				Data:     []byte(testdata),
			},
		},
		Headers: map[string]string{},
	}
	var buffer bytes.Buffer
	mail.Make(&buffer)
	expect.Equal(t, buffer.String(), strings.ReplaceAll(expected, "\n", "\r\n"))
}

const testmessage2 = "Lorem ipsum dolor sit amet, consectetur adipiscing elit.  Nunc imperdiet elit eu sapien accumsan, quis accumsan nisl porta. Donec congue dui in dignissim tincidunt. Phasellus vel ligula lobortis tortor iaculis vulputate.\n\nProin non augue quis est molestie dignissim eget sit amet enim.  Donec ut purus ac enim hendrerit ornare. Pellentesque egestas tempor sodales.  Pellentesque eget auctor mauris."

const expected2 = `Date: Sat, 01 Jan 2022 00:00:00 +0000
Subject: test2
From: "Z" <z@example.com>
To: "A" <a@example.com>
Mime-Version: 1.0
Content-Type: text/plain; charset="UTF-8"
Content-Transfer-Encoding: quoted-printable

Lorem ipsum dolor sit amet, consectetur adipiscing elit.  Nunc imperdiet el=
it eu sapien accumsan, quis accumsan nisl porta. Donec congue dui in dignis=
sim tincidunt. Phasellus vel ligula lobortis tortor iaculis vulputate.

Proin non augue quis est molestie dignissim eget sit amet enim.  Donec ut p=
urus ac enim hendrerit ornare. Pellentesque egestas tempor sodales.  Pellen=
tesque eget auctor mauris.
`

func TestEmail2(t *testing.T) {
	m := Email{
		Date:    time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:      []Address{Address{"A", "a@example.com"}},
		From:    Address{"Z", "z@example.com"},
		Subject: "test2",
		Content: testmessage2,
		Headers: map[string]string{},
	}
	var buffer bytes.Buffer
	m.Make(&buffer)
	expect.Equal(t, buffer.String(), strings.ReplaceAll(expected2, "\n", "\r\n"))
}

func addressList(header mail.Header, key string) ([]Address, error) {
	var result []Address
	list, err := header.AddressList(key)
	if err == nil {
		result = make([]Address, 0, len(list))
		for _, a := range list {
			if a != nil {
				result = append(result, *a)
			}
		}
	}
	return result, err
}

func TestEmail3(t *testing.T) {
	m := Email{
		Date:    time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:      []Address{Address{"A ← ↑ → ↓", "a@example.com"}},
		From:    Address{"Z ← ↑ → ↓", "z@example.com"},
		Subject: "test2 ←↑ →↓",
		Content: testmessage2,
		Headers: nil,
	}
	e, _ := decodeMessage(encodeBytes(m))

	expect.True(t, m.Date.Equal(e.Date))
	expect.DeepEqual(t, m.To, e.To)
	expect.DeepEqual(t, m.From, e.From)
	expect.DeepEqual(t, m.Cc, e.Cc)
	expect.DeepEqual(t, m.Subject, e.Subject)
	expect.DeepEqual(t, m.Content, e.Content)
	expect.DeepEqual(t, m.Attachments, e.Attachments)
	expect.DeepEqual(t, m.Headers, e.Headers)
	expect.True(t, e.Equal(m))
}

func equalAddressArrays(u, v []Address) bool {
	if len(u) != len(v) {
		return false
	}
	for i, a := range u {
		if a != v[i] {
			return false
		}
	}
	return true
}

func (e Email) Equal(m Email) bool {
	return e.Date.Equal(m.Date) &&
		equalAddressArrays(e.To, m.To) &&
		equalAddressArrays(e.Cc, m.Cc) &&
		equalAddressArrays(e.Bcc, m.Bcc) &&
		e.From == m.From &&
		e.Subject == m.Subject &&
		e.Content == m.Content
}

func encodeBytes(m Email) []byte {
	var buffer bytes.Buffer
	m.Make(&buffer)
	return buffer.Bytes()
}

func encodeString(m Email) string {
	var buffer bytes.Buffer
	m.Make(&buffer)
	return buffer.String()
}

func decodeMessage(message []byte) (Email, error) {
	var email Email
	msg, err := mail.ReadMessage(bytes.NewReader(message))
	if err != nil || msg == nil {
		return email, err
	}
	email.Date, err = msg.Header.Date()
	if err != nil {
		return email, err
	}
	email.To, err = addressList(msg.Header, "To")
	email.Cc, err = addressList(msg.Header, "Cc")
	email.Bcc, err = addressList(msg.Header, "Bcc")

	fromList, err := msg.Header.AddressList("From")
	if len(fromList) > 0 && fromList[0] != nil {
		email.From = *(fromList[0])
	}
	wordDecoder := mime.WordDecoder{}
	email.Subject, _ = wordDecoder.DecodeHeader(msg.Header.Get("Subject"))

	knownHeaders := map[string]struct{}{
		"Bcc":                       struct{}{},
		"Cc":                        struct{}{},
		"Content-Transfer-Encoding": struct{}{},
		"Content-Type":              struct{}{},
		"Date":                      struct{}{},
		"From":                      struct{}{},
		"Mime-Version":              struct{}{},
		"Subject":                   struct{}{},
		"To":                        struct{}{},
	}

	for k, _ := range msg.Header {
		_, ok := knownHeaders[k]
		if !ok {
			if email.Headers == nil {
				email.Headers = map[string]string{}
			}
			email.Headers[k] = msg.Header.Get(k)
		}
	}

	body, _ := io.ReadAll(msg.Body)
	if msg.Header.Get("Content-Transfer-Encoding") == "quoted-printable" {
		qpr := quotedprintable.NewReader(bytes.NewReader(body))
		body, _ = io.ReadAll(qpr)
	}
	email.Content = strings.ReplaceAll(strings.TrimSpace(string(body)), "\r\n", "\n")

	// 	Content     string
	// 	Attachments []Attachment
	return email, nil
}

func decodeMimeHeader(s string) (string, error) {
	d := mime.WordDecoder{}
	return d.DecodeHeader(s)
}

func qencodeString(s string) string {
	var b bytes.Buffer
	qencode(&b, s)
	return b.String()
}

func TestMimeHeader(t *testing.T) {
	for _, s := range []string{
		"HELLO WORLD",
		"test2 ←↑ →↓",
		"HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO HELLO ←↑ →↓",
	} {
		q := qencodeString(s)
		v, _ := decodeMimeHeader(q)
		expect.Equal(t, s, v)
	}

}

// //
//
//	func foo() {
//		mediaType, _, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
//		if err != nil {
//			t.Fatal(err)
//		}
//		if strings.HasPrefix(mediaType, "multipart/") {
//		} else {
//			switch strings.ToUpper(msg.Header.Get("Content-Transfer-Encoding")) {
//			case "BASE64":
//				encoded, _ := io.ReadAll(msg.Body)
//				content, _ := base64.StdEncoding.DecodeString(string(encoded))
//				t.Logf("\n%s\n", content)
//			case "QUOTED-PRINTABLE":
//				content, _ := io.ReadAll(quotedprintable.NewReader(msg.Body))
//				t.Logf("\n%s\n", content)
//			default:
//				content, _ := io.ReadAll(msg.Body)
//				t.Logf("\n%s\n", content)
//			}
//		}
//
// }

const ics = `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//ical.marudot.com//iCal Event Maker
CALSCALE:GREGORIAN
BEGIN:VTIMEZONE
TZID:America/New_York
LAST-MODIFIED:20201011T015911Z
TZURL:http://tzurl.org/zoneinfo-outlook/America/New_York
X-LIC-LOCATION:America/New_York
BEGIN:DAYLIGHT
TZNAME:EDT
TZOFFSETFROM:-0500
TZOFFSETTO:-0400
DTSTART:19700308T020000
RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=2SU
END:DAYLIGHT
BEGIN:STANDARD
TZNAME:EST
TZOFFSETFROM:-0400
TZOFFSETTO:-0500
DTSTART:19701101T020000
RRULE:FREQ=YEARLY;BYMONTH=11;BYDAY=1SU
END:STANDARD
END:VTIMEZONE
BEGIN:VEVENT
DTSTAMP:20230424T131500Z
UID:1682342079292-55595@ical.marudot.com
DTSTART;TZID=America/New_York:20230429T120000
DTEND;TZID=America/New_York:20230429T130000
SUMMARY:FOO BAR
DESCRIPTION:This is a test.
LOCATION:9 UPTON CT\, DURHAM NC 27713-7573
END:VEVENT
END:VCALENDAR
`

var icsExpected = toDos(`Date: Sat, 01 Jan 2022 00:00:00 +0000
Subject: test
From: "Z" <z@example.com>
To: "A" <a@example.com>
Mime-Version: 1.0
Content-Type: multipart/mixed; boundary="================"

--================
Content-Transfer-Encoding: quoted-printable
Content-Type: text/plain; charset="UTF-8"

`) + toQuotedPrintable(testmessage2) + toDos(`
--================
Content-Disposition: attachment; filename="invite.ics"
Content-Transfer-Encoding: quoted-printable
Content-Type: text/calendar
MIME-Version: 1.0

`) + toQuotedPrintable(ics) + toDos(`
--================--
`)

func toDos(s string) string {
	return strings.ReplaceAll(s, "\n", "\r\n")
}

func TestEmail4(t *testing.T) {
	m := Email{
		Date:    time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:      []Address{Address{"A", "a@example.com"}},
		From:    Address{"Z", "z@example.com"},
		Subject: "test",
		Content: testmessage2,
		Headers: map[string]string{},
		Attachments: []Attachment{
			Attachment{
				Filename:    "invite.ics",
				ContentType: "text/calendar",
				Data:        []byte(ics),
				Textual:     true,
			},
		},
	}
	var buffer bytes.Buffer
	m.Make(&buffer)
	expect.Equal(t, buffer.String(), icsExpected)
}
