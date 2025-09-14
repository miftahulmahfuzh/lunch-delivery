# 📧 Gmail SMTP Setup Guide

Your SMTP test script is working correctly! The authentication error is expected because Gmail requires special setup.

## Current Status ✅
- ✅ SMTP script created successfully
- ✅ Configuration loaded properly
- ✅ Connection to Gmail SMTP server successful
- ❌ Authentication failed (needs App Password)

## Gmail App Password Setup

### Step 1: Enable 2-Factor Authentication
1. Go to [Google Account Security](https://myaccount.google.com/security)
2. Under "Signing in to Google", click "2-Step Verification"
3. Follow the setup process if not already enabled

### Step 2: Generate App Password
1. Go to [Google App Passwords](https://myaccount.google.com/apppasswords)
2. Select "Mail" from the dropdown
3. Generate password
4. Copy the 16-character password (like: `abcd efgh ijkl mnop`)

### Step 3: Update .env File
Replace your current password with the App Password:

```bash
# Replace this line in your .env file:
SMTP_PASSWORD=pondokku78

# With the App Password (remove spaces):
SMTP_PASSWORD=abcdefghijklmnop
```

### Step 4: Test Again
```bash
go run scripts/smtp/send.go
```

## Alternative: Different Email Provider

If Gmail is too restrictive, you could use:

### Outlook/Hotmail
```bash
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USERNAME=your-email@hotmail.com
SMTP_PASSWORD=your-regular-password
```

### SendGrid (Recommended for production)
```bash
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your-sendgrid-api-key
```

## Quick Test with Different Provider

If you want to test immediately, try with Outlook:

```bash
# Temporarily in .env file
SMTP_HOST=smtp-mail.outlook.com
SMTP_USERNAME=mahfuzh74@hotmail.co.uk
SMTP_PASSWORD=your-hotmail-password
```

---

**Next Steps:**
1. Set up Gmail App Password, OR
2. Try with Outlook/Hotmail credentials, OR
3. Continue development with current setup (emails will be logged to console)