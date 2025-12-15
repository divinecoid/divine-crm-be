package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	host     string
	port     int
	user     string
	pass     string
	from     string
	fromName string
}

// TemplateData matches the .hbs template variables (snake_case)
type TemplateData map[string]interface{}

func NewMailer() *Mailer {
	port, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	return &Mailer{
		host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		port:     port,
		user:     os.Getenv("SMTP_USER"),
		pass:     os.Getenv("SMTP_PASS"),
		from:     getEnv("SMTP_FROM", "Divine CRM <noreply@divine.com>"),
		fromName: getEnv("SMTP_FROM_NAME", "Divine CRM"),
	}
}

func (m *Mailer) IsConfigured() bool {
	return m.user != "" && m.pass != ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getBaseData returns common template data used in all emails
func getBaseData() TemplateData {
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")

	return TemplateData{
		"company_name":    getEnv("COMPANY_NAME", "Divine CRM"),
		"company_address": getEnv("COMPANY_ADDRESS", "Jakarta, Indonesia"),
		"website_url":     getEnv("WEBSITE_URL", frontendURL),
		"support_url":     getEnv("SUPPORT_URL", frontendURL+"/support"),
		"support_email":   getEnv("SUPPORT_EMAIL", "support@divinecrm.com"),
		"privacy_url":     getEnv("PRIVACY_URL", frontendURL+"/privacy"),
		"current_year":    time.Now().Year(),
		"dashboard_url":   frontendURL + "/dashboard",
		"docs_url":        getEnv("DOCS_URL", frontendURL+"/docs"),
		"login_url":       frontendURL + "/login",
	}
}

// generateCodeFromToken generates a 6-digit code from token
func generateCodeFromToken(token string) string {
	if len(token) < 6 {
		return "000000"
	}

	code := ""
	for _, c := range token {
		if len(code) >= 6 {
			break
		}
		if c >= '0' && c <= '9' {
			code += string(c)
		} else if c >= 'a' && c <= 'z' {
			code += fmt.Sprintf("%d", (c-'a')%10)
		} else if c >= 'A' && c <= 'Z' {
			code += fmt.Sprintf("%d", (c-'A')%10)
		}
	}

	for len(code) < 6 {
		code += "0"
	}

	return code[:6]
}

func (m *Mailer) SendVerificationEmail(email, fullName, token string) error {
	// Use backend URL for verification endpoint
	appURL := getEnv("APP_URL", "http://localhost:3001")

	data := getBaseData()
	data["user_name"] = fullName
	data["user_email"] = email
	data["verification_url"] = fmt.Sprintf("%s/api/v1/auth/verify-email/%s", appURL, token)
	data["verification_code"] = generateCodeFromToken(token)
	data["expire_time"] = "24 hours"

	return m.sendEmail(email, "Verify Your Email - Divine CRM", "verification.hbs", data)
}

func (m *Mailer) SendAccountCreatedEmail(email, fullName, role, tempPassword, token string) error {
	// Use backend URL for verification endpoint
	appURL := getEnv("APP_URL", "http://localhost:3001")

	data := getBaseData()
	data["user_name"] = fullName
	data["user_email"] = email
	data["user_role"] = role
	data["created_date"] = time.Now().Format("January 2, 2006")

	if tempPassword != "" {
		data["temp_password"] = tempPassword
	}

	if token != "" {
		data["verification_url"] = fmt.Sprintf("%s/api/v1/auth/verify-email/%s", appURL, token)
	}

	return m.sendEmail(email, "Your Account Has Been Created - Divine CRM", "account-created.hbs", data)
}

func (m *Mailer) SendWelcomeEmail(email, fullName, role string) error {
	data := getBaseData()
	data["user_name"] = fullName
	data["user_email"] = email
	data["user_role"] = role

	return m.sendEmail(email, "Welcome to Divine CRM", "welcome.hbs", data)
}

func (m *Mailer) SendPasswordResetEmail(email, fullName, token string) error {
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")

	data := getBaseData()
	data["user_name"] = fullName
	data["user_email"] = email
	data["reset_url"] = fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
	data["reset_code"] = generateCodeFromToken(token)
	data["expire_time"] = "1 hour"

	return m.sendEmail(email, "Reset Your Password - Divine CRM", "password-reset.hbs", data)
}

func (m *Mailer) sendEmail(to, subject, templateFile string, data TemplateData) error {
	templatePath := filepath.Join("templates", templateFile)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("❌ Failed to parse email template %s: %v", templateFile, err)
		return m.sendFallbackEmail(to, subject, data)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Printf("❌ Failed to execute email template: %v", err)
		return m.sendFallbackEmail(to, subject, data)
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body.String())

	dialer := gomail.NewDialer(m.host, m.port, m.user, m.pass)

	if err := dialer.DialAndSend(msg); err != nil {
		log.Printf("❌ Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("✅ Email sent to %s: %s", to, subject)
	return nil
}

func (m *Mailer) sendFallbackEmail(to, subject string, data TemplateData) error {
	userName := getStringValue(data, "user_name", "User")
	companyName := getStringValue(data, "company_name", "Divine CRM")
	currentYear := getIntValue(data, "current_year", time.Now().Year())

	var body string

	// Verification email
	if verifyURL, ok := data["verification_url"].(string); ok && verifyURL != "" && data["temp_password"] == nil && data["reset_url"] == nil {
		expireTime := getStringValue(data, "expire_time", "24 hours")

		body = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
<body style="margin: 0; padding: 0; background-color: #f4f4f5; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;">
    <div style="max-width: 600px; margin: 40px auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1);">
        <div style="background: linear-gradient(135deg, #f59e0b 0%%, #d97706 100%%); padding: 40px; text-align: center;">
            <h1 style="color: #ffffff; font-size: 24px; margin: 0;">%s</h1>
        </div>
        <div style="padding: 48px;">
            <h2 style="color: #18181b; font-size: 20px; margin: 0 0 24px;">Verify Your Email Address</h2>
            <p style="color: #52525b; font-size: 15px; line-height: 1.7;">Dear %s,</p>
            <p style="color: #52525b; font-size: 15px; line-height: 1.7;">Please verify your email address by clicking the button below:</p>
            
            
            <div style="text-align: center; margin: 32px 0;">
                <a href="%s" style="display: inline-block; background: linear-gradient(135deg, #f59e0b 0%%, #d97706 100%%); color: #ffffff; text-decoration: none; padding: 14px 40px; border-radius: 8px; font-size: 15px; font-weight: 600;">Verify Email Address</a>
            </div>
            
            <p style="color: #71717a; font-size: 14px; text-align: center;">This link expires in %s</p>
            
            <hr style="border: none; border-top: 1px solid #e4e4e7; margin: 32px 0;">
            <p style="color: #9ca3af; font-size: 13px;">If the button doesn't work, copy and paste this URL:</p>
            <p style="color: #6b7280; font-size: 13px; word-break: break-all; background: #f4f4f5; padding: 12px; border-radius: 6px;">%s</p>
        </div>
        <div style="background: #fafafa; padding: 24px; text-align: center; border-top: 1px solid #e4e4e7;">
            <p style="color: #a1a1aa; font-size: 13px; margin: 0;">&copy; %d %s. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, companyName, userName, verifyURL, expireTime, verifyURL, currentYear, companyName)

	} else if resetURL, ok := data["reset_url"].(string); ok && resetURL != "" {
		// Password reset email
		resetCode := getStringValue(data, "reset_code", "")
		expireTime := getStringValue(data, "expire_time", "1 hour")

		body = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
<body style="margin: 0; padding: 0; background-color: #f4f4f5; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;">
    <div style="max-width: 600px; margin: 40px auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1);">
        <div style="background: linear-gradient(135deg, #f59e0b 0%%, #d97706 100%%); padding: 40px; text-align: center;">
            <h1 style="color: #ffffff; font-size: 24px; margin: 0;">%s</h1>
        </div>
        <div style="padding: 48px;">
            <h2 style="color: #18181b; font-size: 20px; margin: 0 0 24px;">Password Reset Request</h2>
            <p style="color: #52525b; font-size: 15px; line-height: 1.7;">Dear %s,</p>
            <p style="color: #52525b; font-size: 15px; line-height: 1.7;">We received a request to reset your password. Click the button below to proceed:</p>
            
            <div style="background: #f4f4f5; border: 2px dashed #d4d4d8; border-radius: 8px; padding: 24px; margin: 28px 0; text-align: center;">
                <p style="color: #71717a; font-size: 13px; margin: 0 0 12px;">Your verification code</p>
                <p style="font-size: 32px; font-weight: 700; font-family: monospace; letter-spacing: 8px; color: #18181b; margin: 0;">%s</p>
            </div>
            
            <div style="text-align: center; margin: 32px 0;">
                <a href="%s" style="display: inline-block; background: linear-gradient(135deg, #f59e0b 0%%, #d97706 100%%); color: #ffffff; text-decoration: none; padding: 14px 40px; border-radius: 8px; font-size: 15px; font-weight: 600;">Reset Your Password</a>
            </div>
            
            <p style="color: #71717a; font-size: 14px; text-align: center;">This link expires in %s</p>
            
            <div style="background: #fef2f2; border: 1px solid #fecaca; border-left: 4px solid #dc2626; border-radius: 8px; padding: 20px; margin: 28px 0;">
                <p style="color: #991b1b; font-size: 14px; font-weight: 600; margin: 0 0 8px;">Did not request this?</p>
                <p style="color: #b91c1c; font-size: 14px; margin: 0;">If you did not request a password reset, please ignore this email.</p>
            </div>
            
            <hr style="border: none; border-top: 1px solid #e4e4e7; margin: 32px 0;">
            <p style="color: #9ca3af; font-size: 13px;">If the button doesn't work, copy and paste this URL:</p>
            <p style="color: #6b7280; font-size: 13px; word-break: break-all; background: #f4f4f5; padding: 12px; border-radius: 6px;">%s</p>
        </div>
        <div style="background: #fafafa; padding: 24px; text-align: center; border-top: 1px solid #e4e4e7;">
            <p style="color: #a1a1aa; font-size: 13px; margin: 0;">&copy; %d %s. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, companyName, userName, resetCode, resetURL, expireTime, resetURL, currentYear, companyName)

	} else if verifyURL, ok := data["verification_url"].(string); ok && verifyURL != "" {
		userEmail := getStringValue(data, "user_email", "")
		userRole := getStringValue(data, "user_role", "User")
		createdDate := getStringValue(data, "created_date", time.Now().Format("January 2, 2006"))
		tempPassword := getStringValue(data, "temp_password", "")

		tempPasswordHTML := ""
		if tempPassword != "" {
			tempPasswordHTML = fmt.Sprintf(`
            <div style="background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 20px; margin: 20px 0;">
                <p style="color: #991b1b; font-size: 14px; margin: 0;">Your temporary password is:</p>
                <span style="font-family: monospace; font-size: 18px; font-weight: 700; color: #dc2626; background: #ffffff; padding: 8px 16px; border-radius: 6px; display: inline-block; margin-top: 12px; letter-spacing: 1px;">%s</span>
                <p style="margin-top: 12px; font-size: 13px; color: #991b1b;">Please change your password immediately after your first login.</p>
            </div>`, tempPassword)
		}

		body = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
<body style="margin: 0; padding: 0; background-color: #f4f4f5; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;">
    <div style="max-width: 600px; margin: 40px auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1);">
        <div style="background: linear-gradient(135deg, #f59e0b 0%%, #d97706 100%%); padding: 40px; text-align: center;">
            <h1 style="color: #ffffff; font-size: 24px; margin: 0;">%s</h1>
        </div>
        <div style="padding: 48px;">
            <h2 style="color: #18181b; font-size: 20px; margin: 0 0 24px;">Your Account Has Been Created</h2>
            <p style="color: #52525b; font-size: 15px; line-height: 1.7;">Dear %s,</p>
            <p style="color: #52525b; font-size: 15px; line-height: 1.7;">Your account has been successfully created. Please verify your email to activate your account.</p>
            
            <div style="background: #fefce8; border: 1px solid #fef08a; border-radius: 8px; padding: 24px; margin: 28px 0;">
                <table style="width: 100%%;">
                    <tr><td style="color: #78716c; font-size: 14px; padding: 6px 0;">Name</td><td style="color: #18181b; font-size: 14px; font-weight: 600;">%s</td></tr>
                    <tr><td style="color: #78716c; font-size: 14px; padding: 6px 0;">Email</td><td style="color: #18181b; font-size: 14px; font-weight: 600;">%s</td></tr>
                    <tr><td style="color: #78716c; font-size: 14px; padding: 6px 0;">Role</td><td style="color: #18181b; font-size: 14px; font-weight: 600;">%s</td></tr>
                    <tr><td style="color: #78716c; font-size: 14px; padding: 6px 0;">Created</td><td style="color: #18181b; font-size: 14px; font-weight: 600;">%s</td></tr>
                </table>
            </div>
            
            %s
            
            <div style="background: linear-gradient(135deg, #ecfdf5 0%%, #d1fae5 100%%); border: 1px solid #6ee7b7; border-radius: 12px; padding: 32px; margin: 28px 0; text-align: center;">
                <h3 style="color: #065f46; font-size: 16px; font-weight: 600; margin: 0 0 12px;">Verify Your Email Address</h3>
                <p style="color: #047857; font-size: 14px; margin: 0 0 20px;">Click the button below to verify your email and activate your account.</p>
                <a href="%s" style="display: inline-block; background: linear-gradient(135deg, #10b981 0%%, #059669 100%%); color: #ffffff; text-decoration: none; padding: 14px 32px; border-radius: 8px; font-size: 15px; font-weight: 600;">Verify Email Address</a>
            </div>
            
            <p style="color: #71717a; font-size: 14px; text-align: center;">This verification link will expire in 24 hours.</p>
            
            <hr style="border: none; border-top: 1px solid #e4e4e7; margin: 32px 0;">
            <p style="color: #9ca3af; font-size: 13px;">If the button doesn't work, copy and paste this URL:</p>
            <p style="color: #6b7280; font-size: 12px; word-break: break-all; background: #f4f4f5; padding: 12px; border-radius: 6px;">%s</p>
        </div>
        <div style="background: #fafafa; padding: 24px; text-align: center; border-top: 1px solid #e4e4e7;">
            <p style="color: #a1a1aa; font-size: 13px; margin: 0;">&copy; %d %s. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, companyName, userName, userName, userEmail, userRole, createdDate, tempPasswordHTML, verifyURL, verifyURL, currentYear, companyName)

	} else {
		// Welcome email
		dashboardURL := getStringValue(data, "dashboard_url", "")

		body = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
<body style="margin: 0; padding: 0; background-color: #f4f4f5; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;">
    <div style="max-width: 600px; margin: 40px auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1);">
        <div style="background: linear-gradient(135deg, #f59e0b 0%%, #d97706 100%%); padding: 48px; text-align: center;">
            <h1 style="color: #ffffff; font-size: 28px; margin: 0 0 8px;">Welcome to %s</h1>
            <p style="color: rgba(255,255,255,0.9); font-size: 16px; margin: 0;">We are excited to have you on board</p>
        </div>
        <div style="padding: 48px;">
            <h2 style="color: #18181b; font-size: 22px; margin: 0 0 24px;">Hello, %s</h2>
            <p style="color: #52525b; font-size: 15px; line-height: 1.7;">Thank you for joining %s. Your account has been successfully set up and you are now part of our community.</p>
            
            <div style="background: linear-gradient(135deg, #fefce8 0%%, #fef3c7 100%%); border-radius: 12px; padding: 32px; margin: 28px 0; text-align: center;">
                <h3 style="color: #78350f; font-size: 18px; margin: 0 0 12px;">Your journey starts here</h3>
                <p style="color: #92400e; font-size: 15px; margin: 0;">Explore our platform and discover all the powerful features designed to help you succeed.</p>
            </div>
            
            <div style="text-align: center; margin: 36px 0;">
                <a href="%s" style="display: inline-block; background: linear-gradient(135deg, #f59e0b 0%%, #d97706 100%%); color: #ffffff; text-decoration: none; padding: 16px 40px; border-radius: 8px; font-size: 16px; font-weight: 600;">Go to Dashboard</a>
            </div>
        </div>
        <div style="background: #fafafa; padding: 24px; text-align: center; border-top: 1px solid #e4e4e7;">
            <p style="color: #a1a1aa; font-size: 13px; margin: 0;">&copy; %d %s. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, companyName, userName, companyName, dashboardURL, currentYear, companyName)
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dialer := gomail.NewDialer(m.host, m.port, m.user, m.pass)

	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	log.Printf("📧 Fallback email sent to %s: %s", to, subject)
	return nil
}

// Helper functions to safely get values from TemplateData
func getStringValue(data TemplateData, key, defaultValue string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return defaultValue
}

func getIntValue(data TemplateData, key string, defaultValue int) int {
	if val, ok := data[key].(int); ok {
		return val
	}
	return defaultValue
}
