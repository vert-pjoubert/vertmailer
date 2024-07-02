package vertmailer

import (
	"testing"
)

func TestValidateMail(t *testing.T) {
	tests := []struct {
		name    string
		mail    Mail
		wantErr bool
	}{
		{
			name: "Valid email",
			mail: Mail{
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "Valid Subject",
				Body:    "Valid Body",
			},
			wantErr: false,
		},
		{
			name: "Invalid sender email",
			mail: Mail{
				From:    "invalid-email",
				To:      []string{"recipient@example.com"},
				Subject: "Valid Subject",
				Body:    "Valid Body",
			},
			wantErr: true,
		},
		{
			name: "Invalid recipient email",
			mail: Mail{
				From:    "sender@example.com",
				To:      []string{"invalid-email"},
				Subject: "Valid Subject",
				Body:    "Valid Body",
			},
			wantErr: true,
		},
		{
			name: "Empty subject",
			mail: Mail{
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "",
				Body:    "Valid Body",
			},
			wantErr: true,
		},
		{
			name: "Empty body",
			mail: Mail{
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "Valid Subject",
				Body:    "",
			},
			wantErr: true,
		},
		{
			name: "Whitespace-only subject",
			mail: Mail{
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "   ",
				Body:    "Valid Body",
			},
			wantErr: true,
		},
		{
			name: "Whitespace-only body",
			mail: Mail{
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "Valid Subject",
				Body:    "   ",
			},
			wantErr: true,
		},
		{
			name: "Multiple valid recipients",
			mail: Mail{
				From:    "sender@example.com",
				To:      []string{"recipient1@example.com", "recipient2@example.com"},
				Subject: "Valid Subject",
				Body:    "Valid Body",
			},
			wantErr: false,
		},
		{
			name: "Mixed valid and invalid recipients",
			mail: Mail{
				From:    "sender@example.com",
				To:      []string{"recipient1@example.com", "invalid-email"},
				Subject: "Valid Subject",
				Body:    "Valid Body",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateMail(tt.mail); (err != nil) != tt.wantErr {
				t.Errorf("ValidateMail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
