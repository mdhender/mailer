/*******************************************************************************
mailer: interfaces to send and fetch e-mail: https://github.com/mdhender/mailer

Copyright (c) 2021 Michael D Henderson.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
******************************************************************************/

// Package mailer defines interfaces to send and fetch messages. It also
// provides a simple template engine and transformers to "fix up" line
// endings for message bodies.
package mailer

import (
	"net/mail"
)

// Fetcher is the interface for retrieving e-mails from the server and storing
// them as files in the mailbox directory.
type Fetcher interface {
	Get() error
	Errors() []error
}

// Sender is the interface for generating e-mails and connecting to the server
// to send them.
type Sender interface {
	Send(e *Envelope, t *Template, name string, data interface{}, attachments ...*Attachment) []error
	SendHTML(e *Envelope, t *Template, name string, data interface{}, attachments ...*Attachment) []error
	Errors() []error
}

// Envelope specifies the e-mail envelope information.
type Envelope struct {
	Sender  mail.Address
	To      []mail.Address
	Cc      []mail.Address
	Bcc     []mail.Address
	Subject string
}

// Attachment implements e-mail attachments.
type Attachment struct {
	Name string
	Data []byte
}
