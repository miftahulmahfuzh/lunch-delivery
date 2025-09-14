# SMTP Test Email Sender

This script tests your SMTP configuration by sending a real test email.

## Usage

### 1. Set up your environment variables:

Create or update your `.env` file with SMTP credentials:

```bash
# Required SMTP Configuration
SMTP_HOST=smtp.gmail.com          # or your SMTP server
SMTP_PORT=587                     # or 465 for SSL, 25 for plain
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password   # Use App Password for Gmail
SMTP_FROM=your-email@gmail.com    # Optional, defaults to SMTP_USERNAME

# Test target email
SMTP_TEST_EMAIL_ADDRESS=mahfuzh74@hotmail.co.uk
```

### 2. Run the test:

```bash
# From the project root directory
go run scripts/smtp/send.go
```

## Gmail Setup

For Gmail, you need to:

1. Enable 2-Factor Authentication
2. Generate an App Password:
   - Go to Google Account settings
   - Security → 2-Step Verification → App passwords
   - Generate password for "Mail"
   - Use this password in `SMTP_PASSWORD`

## Common SMTP Settings

### Gmail
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

### Outlook/Hotmail
```
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
```

### Yahoo
```
SMTP_HOST=smtp.mail.yahoo.com
SMTP_PORT=587
```

## Troubleshooting

If the email fails to send, check:

1. ✅ SMTP credentials are correct
2. ✅ App passwords enabled (Gmail)
3. ✅ Less secure app access enabled (if applicable)
4. ✅ Firewall allows SMTP traffic
5. ✅ Recipient email address is valid
6. ✅ SMTP server settings are correct