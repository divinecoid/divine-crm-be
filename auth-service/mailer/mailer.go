package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed templates/*.hbs
var templateFS embed.FS

// Config holds SMTP configuration
type Config struct {
	Host        string
	Port        string
	Username    string
	Password    string
	FromEmail   string
	FromName    string
	AppURL      string
	FrontendURL string
}

// Mailer handles email sending
type Mailer struct {
	config    Config
	templates map[string]*template.Template
}

// EmailData holds common email data
type EmailData struct {
	// Common fields
	AppName     string
	AppURL      string
	FrontendURL string
	Year        int

	// Recipient info
	RecipientName  string
	RecipientEmail string

	// Template specific
	VerifyURL    string
	ResetURL     string
	Token        string
	ExpiryHours  int
	Role         string
	TempPassword string

	// Custom data
	Extra map[string]interface{}
}

// NewMailer creates a new mailer instance
func NewMailer() *Mailer {
	config := Config{
		Host:        getEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:        getEnv("SMTP_PORT", "587"),
		Username:    os.Getenv("SMTP_USER"),
		Password:    os.Getenv("SMTP_PASS"),
		FromEmail:   os.Getenv("SMTP_USER"),
		FromName:    getEnv("SMTP_FROM_NAME", "Divine CRM"),
		AppURL:      getEnv("APP_URL", "http://localhost:3001"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}

	// Parse SMTP_FROM if it contains name and email
	if fromStr := os.Getenv("SMTP_FROM"); fromStr != "" {
		if strings.Contains(fromStr, "<") && strings.Contains(fromStr, ">") {
			parts := strings.Split(fromStr, "<")
			config.FromName = strings.TrimSpace(parts[0])
			config.FromEmail = strings.Trim(parts[1], ">")
		} else {
			config.FromEmail = fromStr
		}
	}

	m := &Mailer{
		config:    config,
		templates: make(map[string]*template.Template),
	}

	m.loadTemplates()

	return m
}

// loadTemplates loads all HBS templates
func (m *Mailer) loadTemplates() {
	// Define template functions
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
		"now":   func() time.Time { return time.Now() },
		"year":  func() int { return time.Now().Year() },
	}

	// Read embedded templates
	entries, err := templateFS.ReadDir("templates")
	if err != nil {
		log.Printf("⚠️ Failed to read templates directory: %v", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".hbs") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".hbs")
		content, err := templateFS.ReadFile(filepath.Join("templates", entry.Name()))
		if err != nil {
			log.Printf("⚠️ Failed to read template %s: %v", entry.Name(), err)
			continue
		}

		tmpl, err := template.New(name).Funcs(funcMap).Parse(string(content))
		if err != nil {
			log.Printf("⚠️ Failed to parse template %s: %v", entry.Name(), err)
			continue
		}

		m.templates[name] = tmpl
		log.Printf("📧 Loaded email template: %s", name)
	}
}

// IsConfigured checks if SMTP is properly configured
func (m *Mailer) IsConfigured() bool {
	return m.config.Username != "" && m.config.Password != ""
}

// renderTemplate renders a template with data
func (m *Mailer) renderTemplate(templateName string, data EmailData) (string, error) {
	tmpl, ok := m.templates[templateName]
	if !ok {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	// Set default values
	if data.AppName == "" {
		data.AppName = m.config.FromName
	}
	if data.AppURL == "" {
		data.AppURL = m.config.AppURL
	}
	if data.FrontendURL == "" {
		data.FrontendURL = m.config.FrontendURL
	}
	if data.Year == 0 {
		data.Year = time.Now().Year()
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %v", err)
	}

	return buf.String(), nil
}

// Send sends an email
func (m *Mailer) Send(to, subject, body string) error {
	if !m.IsConfigured() {
		return fmt.Errorf("SMTP not configured")
	}

	from := fmt.Sprintf("%s <%s>", m.config.FromName, m.config.FromEmail)

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
	addr := fmt.Sprintf("%s:%s", m.config.Host, m.config.Port)

	return smtp.SendMail(addr, auth, m.config.FromEmail, []string{to}, msg.Bytes())
}

// SendVerificationEmail sends email verification
func (m *Mailer) SendVerificationEmail(email, fullName, token string) error {
	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify-email/%s", m.config.AppURL, token)

	data := EmailData{
		RecipientName:  fullName,
		RecipientEmail: email,
		VerifyURL:      verifyURL,
		Token:          token,
		ExpiryHours:    24,
	}

	body, err := m.renderTemplate("verification", data)
	if err != nil {
		return err
	}

	subject := "🔐 Verify Your Email - Divine CRM"

	if err := m.Send(email, subject, body); err != nil {
		// Log the verification link for manual verification
		log.Printf("⚠️ Failed to send email to %s: %v", email, err)
		log.Printf("🔗 Manual verification link: %s", verifyURL)
		return err
	}

	log.Printf("📧 Verification email sent to %s", email)
	return nil
}

// SendWelcomeEmail sends welcome email after verification
func (m *Mailer) SendWelcomeEmail(email, fullName, role string) error {
	data := EmailData{
		RecipientName:  fullName,
		RecipientEmail: email,
		Role:           role,
	}

	body, err := m.renderTemplate("welcome", data)
	if err != nil {
		return err
	}

	subject := "🎉 Welcome to Divine CRM!"

	if err := m.Send(email, subject, body); err != nil {
		log.Printf("⚠️ Failed to send welcome email to %s: %v", email, err)
		return err
	}

	log.Printf("📧 Welcome email sent to %s", email)
	return nil
}

// SendPasswordResetEmail sends password reset email
func (m *Mailer) SendPasswordResetEmail(email, fullName, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", m.config.FrontendURL, token)

	data := EmailData{
		RecipientName:  fullName,
		RecipientEmail: email,
		ResetURL:       resetURL,
		Token:          token,
		ExpiryHours:    1,
	}

	body, err := m.renderTemplate("password-reset", data)
	if err != nil {
		return err
	}

	subject := "🔑 Reset Your Password - Divine CRM"

	if err := m.Send(email, subject, body); err != nil {
		log.Printf("⚠️ Failed to send password reset email to %s: %v", email, err)
		return err
	}

	log.Printf("📧 Password reset email sent to %s", email)
	return nil
}

// SendAccountCreatedEmail sends notification when admin creates an account
func (m *Mailer) SendAccountCreatedEmail(email, fullName, role, tempPassword, token string) error {
	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify-email/%s", m.config.AppURL, token)

	data := EmailData{
		RecipientName:  fullName,
		RecipientEmail: email,
		Role:           role,
		TempPassword:   tempPassword,
		VerifyURL:      verifyURL,
		Token:          token,
		ExpiryHours:    24,
	}

	body, err := m.renderTemplate("account-created", data)
	if err != nil {
		// Fallback to verification template if account-created doesn't exist
		return m.SendVerificationEmail(email, fullName, token)
	}

	subject := "🎊 Your Divine CRM Account Has Been Created"

	if err := m.Send(email, subject, body); err != nil {
		log.Printf("⚠️ Failed to send account created email to %s: %v", email, err)
		log.Printf("🔗 Manual verification link: %s", verifyURL)
		return err
	}

	log.Printf("📧 Account created email sent to %s", email)
	return nil
}

// Helper function
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
