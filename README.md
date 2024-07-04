# VertMailer

VertMailer is a Go package that provides a simple interface for sending emails using an SMTP server. It includes support for TLS/SSL connections and can be configured to use a custom CA certificate for secure connections.

## Features

- Send emails using SMTP.
- Support for plain text and HTML email bodies.
- TLS/SSL support with custom CA certificate loading.
- Easy configuration through a `MailServer` struct.

## Installation

To install VertMailer, use `go get`:

```sh
go get github.com/yourusername/vertmailer
```

## Usage

Here's an example of how to use the VertMailer package to send an email:

```go
package main

import (
	"log"
	"github.com/vert-pjoubert/vertmailer"
)

func main() {
	server := vertmailer.MailServer{
		Host:       "smtp.example.com",
		Port:       "587",
		Username:   "your-username",
		Password:   "your-password",
		CACertPath: "/path/to/ca-cert.pem",
		UseTLS:     true,
	}

	mailer := vertmailer.NewMailerService(server)

	mail := vertmailer.Mail{
		From:    "your-email@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    `<h1>This is a test email body.</h1><p style="color:red;">This is a paragraph.</p>`,
	}

	if err := mailer.SendMail(mail); err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully!")
}
```

## MailServer Configuration

The `MailServer` struct is used to configure the SMTP server details:

- `Host`: The SMTP server host.
- `Port`: The SMTP server port.
- `Username`: The username for SMTP authentication.
- `Password`: The password for SMTP authentication.
- `CACertPath`: The path to the CA certificate file for TLS/SSL connections.
- `UseTLS`: A boolean indicating whether to use TLS/SSL for the connection.

## Mail Struct

The `Mail` struct represents the email details:

- `From`: The sender's email address.
- `To`: A slice of recipient email addresses.
- `Subject`: The subject of the email.
- `Body`: The body of the email, which can include HTML content.

## Testing

The package includes a test for sending emails using a mock SMTP server. The mock server supports basic authentication and prints the email message to the console for verification.

To run the tests, use the following command:

```sh
go test -v ./...
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

---
