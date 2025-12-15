package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"divine-crm/auth-service/mailer"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var jwtSecret []byte
var emailMailer *mailer.Mailer

type User struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	Username          string     `json:"username" gorm:"uniqueIndex;not null"`
	Email             string     `json:"email" gorm:"uniqueIndex;not null"`
	Password          string     `json:"-" gorm:"not null"`
	FullName          string     `json:"full_name"`
	Role              string     `json:"role" gorm:"default:'Agent'"`
	Active            bool       `json:"active" gorm:"default:false"`
	EmailVerified     bool       `json:"email_verified" gorm:"default:false"`
	VerificationToken string     `json:"-" gorm:"index"`
	TokenExpiry       *time.Time `json:"-"`
	ResetToken        string     `json:"-" gorm:"index"`
	ResetTokenExpiry  *time.Time `json:"-"`
	LatestLogin       *time.Time `json:"latest_login"`
	CreatedBy         *uint      `json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name"`
	Active   *bool  `json:"active"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

func initDB() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "postgres"),
			getEnv("DB_NAME", "divine_crm"),
			getEnv("DB_SSLMODE", "disable"),
		)
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	db.AutoMigrate(&User{})
	createDefaultSuperadmin()

	log.Println("✅ Auth Service: Database connected")
}

func createDefaultSuperadmin() {
	var count int64
	db.Model(&User{}).Where("role = ?", "Superadmin").Count(&count)

	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("superadmin123"), bcrypt.DefaultCost)
		superadmin := User{
			Username:      "superadmin",
			Email:         "superadmin@divine.com",
			Password:      string(hashedPassword),
			FullName:      "Super Administrator",
			Role:          "Superadmin",
			Active:        true,
			EmailVerified: true,
		}
		db.Create(&superadmin)
		log.Println("✅ Default Superadmin created")
	}
}

func generateToken(user *User) (string, error) {
	expirationHours := 24
	if envHours := os.Getenv("JWT_EXPIRATION_HOURS"); envHours != "" {
		fmt.Sscanf(envHours, "%d", &expirationHours)
	}

	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "divine-crm-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func verifyToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "No authorization header"})
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	claims, err := verifyToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid or expired token"})
	}

	c.Locals("userID", claims.UserID)
	c.Locals("username", claims.Username)
	c.Locals("email", claims.Email)
	c.Locals("fullName", claims.FullName)
	c.Locals("role", claims.Role)

	return c.Next()
}

func requireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role").(string)

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{"success": false, "message": "Access denied. Insufficient permissions."})
	}
}

func login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Email and password are required"})
	}

	var user User
	if err := db.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid email or password"})
	}

	if !user.Active {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Account is deactivated"})
	}

	if !user.EmailVerified && user.Role != "Superadmin" {
		return c.Status(401).JSON(fiber.Map{
			"success":            false,
			"message":            "Please verify your email first. Check your inbox for verification link.",
			"email_not_verified": true,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid email or password"})
	}

	now := time.Now()
	user.LatestLogin = &now
	db.Save(&user)

	token, err := generateToken(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to generate token"})
	}

	log.Printf("✅ Login: %s (%s)", user.Email, user.Role)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"token":   token,
		"user": fiber.Map{
			"id":             user.ID,
			"username":       user.Username,
			"email":          user.Email,
			"full_name":      user.FullName,
			"role":           user.Role,
			"email_verified": user.EmailVerified,
		},
	})
}

func register(c *fiber.Ctx) error {
	currentRole := c.Locals("role").(string)
	currentUserID := c.Locals("userID").(uint)

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Username, email, and password are required"})
	}

	if req.Role == "" {
		req.Role = "Agent"
	}

	switch currentRole {
	case "Superadmin":
		if req.Role != "Admin" && req.Role != "Agent" {
			return c.Status(403).JSON(fiber.Map{"success": false, "message": "Superadmin can only create Admin or Agent"})
		}
	case "Admin":
		if req.Role != "Agent" {
			return c.Status(403).JSON(fiber.Map{"success": false, "message": "Admin can only create Agent"})
		}
	default:
		return c.Status(403).JSON(fiber.Map{"success": false, "message": "You don't have permission to create users"})
	}

	var existingUser User
	if err := db.Where("email = ?", strings.ToLower(req.Email)).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Email already registered"})
	}
	if err := db.Where("username = ?", strings.ToLower(req.Username)).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Username already taken"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to hash password"})
	}

	verificationToken := generateRandomToken()
	tokenExpiry := time.Now().Add(24 * time.Hour)

	user := User{
		Username:          strings.ToLower(req.Username),
		Email:             strings.ToLower(req.Email),
		Password:          string(hashedPassword),
		FullName:          req.FullName,
		Role:              req.Role,
		Active:            false,
		EmailVerified:     false,
		VerificationToken: verificationToken,
		TokenExpiry:       &tokenExpiry,
		CreatedBy:         &currentUserID,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to create user"})
	}

	go sendAccountCreatedEmail(user.Email, user.FullName, user.Role, verificationToken)

	log.Printf("📧 User created (pending verification): %s (%s) by user %d", user.Email, user.Role, currentUserID)

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "User created. Verification email has been sent to " + user.Email,
		"user": fiber.Map{
			"id":             user.ID,
			"username":       user.Username,
			"email":          user.Email,
			"full_name":      user.FullName,
			"role":           user.Role,
			"email_verified": user.EmailVerified,
			"active":         user.Active,
		},
	})
}

func generateRandomToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func sendVerificationEmail(email, fullName, token string) {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:3001"
	}

	if !emailMailer.IsConfigured() {
		log.Printf("⚠️ SMTP not configured. Verification link: %s/api/v1/auth/verify-email/%s", appURL, token)
		return
	}

	if err := emailMailer.SendVerificationEmail(email, fullName, token); err != nil {
		log.Printf("❌ Failed to send verification email: %v", err)
		log.Printf("🔗 Manual verification link: %s/api/v1/auth/verify-email/%s", appURL, token)
	}
}

func sendAccountCreatedEmail(email, fullName, role, token string) {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:3001"
	}

	if !emailMailer.IsConfigured() {
		log.Printf("⚠️ SMTP not configured. Verification link: %s/api/v1/auth/verify-email/%s", appURL, token)
		return
	}

	if err := emailMailer.SendAccountCreatedEmail(email, fullName, role, "", token); err != nil {
		log.Printf("❌ Failed to send account created email: %v", err)
		sendVerificationEmail(email, fullName, token)
	}
}

func sendWelcomeEmail(email, fullName, role string) {
	if !emailMailer.IsConfigured() {
		return
	}

	if err := emailMailer.SendWelcomeEmail(email, fullName, role); err != nil {
		log.Printf("⚠️ Failed to send welcome email: %v", err)
	}
}

func sendPasswordResetEmail(email, fullName, token string) {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	if !emailMailer.IsConfigured() {
		log.Printf("⚠️ SMTP not configured. Password reset link: %s", resetURL)
		return
	}

	if err := emailMailer.SendPasswordResetEmail(email, fullName, token); err != nil {
		log.Printf("❌ Failed to send password reset email: %v", err)
		log.Printf("🔗 Manual reset link: %s", resetURL)
	} else {
		log.Printf("📧 Password reset email sent to %s", email)
	}
}

func verifyEmail(c *fiber.Ctx) error {
	token := c.Params("token")

	var user User
	if err := db.Where("verification_token = ?", token).First(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid or expired verification token",
		})
	}

	if user.TokenExpiry != nil && time.Now().After(*user.TokenExpiry) {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Verification token has expired. Please request a new one.",
		})
	}

	user.Active = true
	user.EmailVerified = true
	user.VerificationToken = ""
	user.TokenExpiry = nil
	db.Save(&user)

	log.Printf("✅ Email verified: %s (%s)", user.Email, user.Role)

	go sendWelcomeEmail(user.Email, user.FullName, user.Role)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Email Verified - Divine CRM</title>
    <style>
        body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); }
        .card { background: white; padding: 50px; border-radius: 20px; box-shadow: 0 25px 50px rgba(0,0,0,0.2); text-align: center; max-width: 450px; }
        .icon { font-size: 64px; margin-bottom: 20px; }
        h1 { color: #1f2937; margin-bottom: 10px; }
        p { color: #6b7280; margin-bottom: 25px; }
        .email { color: #f59e0b; font-weight: bold; }
        .button { display: inline-block; background: linear-gradient(135deg, #f59e0b 0%%, #d97706 100%%); color: white; padding: 15px 40px; text-decoration: none; border-radius: 10px; font-weight: bold; box-shadow: 0 4px 15px rgba(245, 158, 11, 0.4); }
    </style>
</head>
<body>
    <div class="card">
        <h1>Email Verified!</h1>
        <p>Your account <span class="email">%s</span> has been verified and activated.</p>
        <p>You can now log in to Divine CRM.</p>
        <a href="%s/login" class="button">Go to Login</a>
    </div>
</body>
</html>
`, user.Email, frontendURL)

	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}

func resendVerification(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	c.BodyParser(&req)

	if req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Email is required"})
	}

	var user User
	if err := db.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "User not found"})
	}

	if user.EmailVerified {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Email is already verified"})
	}

	user.VerificationToken = generateRandomToken()
	tokenExpiry := time.Now().Add(24 * time.Hour)
	user.TokenExpiry = &tokenExpiry
	db.Save(&user)

	go sendVerificationEmail(user.Email, user.FullName, user.VerificationToken)

	log.Printf("📧 Resent verification email to %s", user.Email)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Verification email has been resent to " + user.Email,
	})
}

func forgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	if req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Email is required"})
	}

	var user User
	if err := db.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "If the email exists, a password reset link has been sent.",
		})
	}

	if !user.Active {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "If the email exists, a password reset link has been sent.",
		})
	}

	user.ResetToken = generateRandomToken()
	resetExpiry := time.Now().Add(1 * time.Hour)
	user.ResetTokenExpiry = &resetExpiry
	db.Save(&user)

	go sendPasswordResetEmail(user.Email, user.FullName, user.ResetToken)

	log.Printf("🔑 Password reset requested for %s", user.Email)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "If the email exists, a password reset link has been sent.",
	})
}

func resetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	if req.Token == "" || req.NewPassword == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Token and new password are required"})
	}

	if len(req.NewPassword) < 6 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Password must be at least 6 characters"})
	}

	var user User
	if err := db.Where("reset_token = ?", req.Token).First(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid or expired reset token"})
	}

	if user.ResetTokenExpiry != nil && time.Now().After(*user.ResetTokenExpiry) {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Reset token has expired. Please request a new one."})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to hash password"})
	}

	user.Password = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpiry = nil
	db.Save(&user)

	log.Printf("✅ Password reset successful for %s", user.Email)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password has been reset successfully. You can now login with your new password.",
	})
}

func verifyResetToken(c *fiber.Ctx) error {
	token := c.Params("token")

	var user User
	if err := db.Where("reset_token = ?", token).First(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid reset token"})
	}

	if user.ResetTokenExpiry != nil && time.Now().After(*user.ResetTokenExpiry) {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Reset token has expired"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Token is valid",
		"email":   user.Email,
	})
}

func getCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "User not found"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":           user.ID,
			"username":     user.Username,
			"email":        user.Email,
			"full_name":    user.FullName,
			"role":         user.Role,
			"active":       user.Active,
			"latest_login": user.LatestLogin,
			"created_at":   user.CreatedAt,
		},
	})
}

func changePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Old password and new password are required"})
	}

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "User not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Old password is incorrect"})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	db.Save(&user)

	return c.JSON(fiber.Map{"success": true, "message": "Password changed successfully"})
}

func verifyTokenHandler(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "No authorization header"})
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	claims, err := verifyToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid or expired token"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"user_id":   claims.UserID,
			"username":  claims.Username,
			"email":     claims.Email,
			"full_name": claims.FullName,
			"role":      claims.Role,
		},
	})
}

func getUsers(c *fiber.Ctx) error {
	currentRole := c.Locals("role").(string)

	var users []User

	switch currentRole {
	case "Superadmin":
		db.Where("role != ?", "Superadmin").Order("created_at desc").Find(&users)
	case "Admin":
		db.Where("role = ?", "Agent").Order("created_at desc").Find(&users)
	default:
		return c.Status(403).JSON(fiber.Map{"success": false, "message": "Access denied"})
	}

	return c.JSON(fiber.Map{"success": true, "data": users})
}

func getUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	currentRole := c.Locals("role").(string)

	var user User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "User not found"})
	}

	if currentRole == "Admin" && user.Role != "Agent" {
		return c.Status(403).JSON(fiber.Map{"success": false, "message": "Access denied"})
	}

	return c.JSON(fiber.Map{"success": true, "data": user})
}

func updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	currentRole := c.Locals("role").(string)

	var user User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "User not found"})
	}

	if currentRole == "Admin" && user.Role != "Agent" {
		return c.Status(403).JSON(fiber.Map{"success": false, "message": "Access denied"})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Active != nil {
		user.Active = *req.Active
	}
	if req.Role != "" && currentRole == "Superadmin" {
		if req.Role == "Admin" || req.Role == "Agent" {
			user.Role = req.Role
		}
	}
	if req.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	db.Save(&user)

	return c.JSON(fiber.Map{"success": true, "data": user, "message": "User updated successfully"})
}

func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	currentRole := c.Locals("role").(string)
	currentUserID := c.Locals("userID").(uint)

	var user User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "User not found"})
	}

	if user.ID == currentUserID {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Cannot delete your own account"})
	}

	if user.Role == "Superadmin" {
		return c.Status(403).JSON(fiber.Map{"success": false, "message": "Cannot delete Superadmin"})
	}

	if currentRole == "Admin" && user.Role != "Agent" {
		return c.Status(403).JSON(fiber.Map{"success": false, "message": "Access denied"})
	}

	db.Delete(&user)

	return c.JSON(fiber.Map{"success": true, "message": "User deleted successfully"})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	godotenv.Load()
	initDB()

	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("divine-crm-jwt-secret-key-2024")
	}

	emailMailer = mailer.NewMailer()
	if emailMailer.IsConfigured() {
		log.Println("📧 Email mailer configured successfully")
	} else {
		log.Println("⚠️ Email mailer not configured - verification links will be logged to console")
	}

	app := fiber.New(fiber.Config{
		AppName: "Divine CRM Auth Service",
	})

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth-service",
			"time":    time.Now(),
		})
	})

	api := app.Group("/api/v1")

	auth := api.Group("/auth")

	auth.Post("/login", login)
	auth.Post("/verify", verifyTokenHandler)
	auth.Post("/forgot-password", forgotPassword)
	auth.Post("/reset-password", resetPassword)
	auth.Get("/reset-password/:token", verifyResetToken)
	auth.Get("/verify-email/:token", verifyEmail)

	auth.Get("/me", authMiddleware, getCurrentUser)
	auth.Post("/change-password", authMiddleware, changePassword)
	auth.Post("/register", authMiddleware, requireRole("Superadmin", "Admin"), register)
	auth.Post("/resend-verification", authMiddleware, requireRole("Superadmin", "Admin"), resendVerification)

	users := api.Group("/users", authMiddleware, requireRole("Superadmin", "Admin"))
	users.Get("/", getUsers)
	users.Get("/:id", getUserByID)
	users.Put("/:id", updateUser)
	users.Delete("/:id", deleteUser)

	port := getEnv("SERVER_PORT", "3001")

	log.Fatal(app.Listen(":" + port))
}
