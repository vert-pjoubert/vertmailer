package vertmailer

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/microcosm-cc/bluemonday"
)

// ValidateMail validates the email fields.
func ValidateMail(mail Mail) error {
	if err := validation.Validate(mail.From, validation.Required, is.Email); err != nil {
		return errors.New("invalid sender email address")
	}

	for _, to := range mail.To {
		if err := validation.Validate(to, validation.Required, is.Email); err != nil {
			return errors.New("invalid recipient email address")
		}
	}

	if strings.TrimSpace(mail.Subject) == "" {
		return errors.New("subject cannot be empty")
	}

	if strings.TrimSpace(mail.Body) == "" {
		return errors.New("body cannot be empty")
	}

	return nil
}

// SanitizeMail sanitizes the email fields.
func SanitizeMail(mail Mail) Mail {
	mail.Subject = sanitizeString(mail.Subject)
	mail.Body = sanitizeHTML(mail.Body)
	return mail
}

// sanitizeString removes any unwanted characters from a string.
func sanitizeString(str string) string {
	return strings.TrimSpace(str)
}

// sanitizeHTML sanitizes HTML content to remove potentially harmful content while preserving formatting.
func sanitizeHTML(html string) string {
	p := bluemonday.UGCPolicy()
	return p.Sanitize(html)
}
