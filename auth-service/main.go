package main

import (
	"log"
	"os"
	"time"

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

// Models
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Password  string    `json:"-"` // Hidden from JSON
	Role      string    `json:"role" gorm:"default:'user'"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex"`
	Description string    `json:"description"`
	Permissions string    `json:"permissions" gorm:"type:text"` // JSON array of permissions
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Session struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Token     string    `json:"token" gorm:"uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role"`
}

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func initDB() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=divine_crm port=5432 sslmode=disable"
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate
	db.AutoMigrate(&User{}, &Role{}, &Session{})

	// Create default roles
	createDefaultRoles()

	// Create default admin user
	createDefaultAdmin()

	log.Println("✅ Auth Service: Database connected and migrated")
}

func createDefaultRoles() {
	roles := []Role{
		{
			Name:        "admin",
			Description: "Administrator with full access",
			Permissions: `["*"]`,
		},
		{
			Name:        "manager",
			Description: "Manager with limited admin access",
			Permissions: `["read:*", "write:contacts", "write:products", "read:analytics"]`,
		},
		{
			Name:        "agent",
			Description: "Customer service agent",
			Permissions: `["read:contacts", "write:chats", "read:products"]`,
		},
		{
			Name:        "user",
			Description: "Regular user with basic access",
			Permissions: `["read:own", "write:own"]`,
		},
	}

	for _, role := range roles {
		var existingRole Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			db.Create(&role)
			log.Printf("✅ Created role: %s", role.Name)
		}
	}
}

func createDefaultAdmin() {
	var admin User
	if err := db.Where("email = ?", "admin@divine.com").First(&admin).Error; err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin = User{
			Name:     "Admin",
			Email:    "admin@divine.com",
			Password: string(hashedPassword),
			Role:     "admin",
			Active:   true,
		}
		db.Create(&admin)
		log.Println("✅ Created default admin user")
		log.Println("   Email: admin@divine.com")
		log.Println("   Password: admin123")
	}
}

func main() {
	godotenv.Load()
	initDB()

	// JWT Secret
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("divine-crm-secret-key-change-in-production")
	}

	app := fiber.New(fiber.Config{
		AppName:      "Auth Service",
		ErrorHandler: errorHandler,
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:3001,http://localhost:3002,http://localhost:8080",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth-service",
			"time":    time.Now(),
		})
	})

	api := app.Group("/api/v1")

	// Public routes (no authentication)
	auth := api.Group("/auth")
	auth.Post("/register", register)
	auth.Post("/login", login)
	auth.Post("/refresh-token", refreshToken)

	// Protected routes (require authentication)
	protected := api.Group("/auth", jwtMiddleware)
	protected.Get("/me", getCurrentUser)
	protected.Post("/logout", logout)
	protected.Put("/profile", updateProfile)
	protected.Put("/change-password", changePassword)

	// Admin routes
	admin := api.Group("/admin", jwtMiddleware, adminMiddleware)
	admin.Get("/users", getUsers)
	admin.Get("/users/:id", getUserByID)
	admin.Put("/users/:id", updateUser)
	admin.Delete("/users/:id", deleteUser)
	admin.Put("/users/:id/activate", activateUser)
	admin.Put("/users/:id/deactivate", deactivateUser)

	// Role management
	admin.Get("/roles", getRoles)
	admin.Post("/roles", createRole)
	admin.Put("/roles/:id", updateRole)
	admin.Delete("/roles/:id", deleteRole)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Printf("🚀 Auth Service starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

// ==================== AUTH HANDLERS ====================

func register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Validate email
	var existingUser User
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Email already registered",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to hash password",
		})
	}

	// Set default role
	if req.Role == "" {
		req.Role = "user"
	}

	// Create user
	user := User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		Active:   true,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create user",
		})
	}

	// Generate JWT token
	token, err := generateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate token",
		})
	}

	// Save session
	session := Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	db.Create(&session)

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Registration successful",
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

func login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Find user
	var user User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Invalid email or password",
		})
	}

	// Check if user is active
	if !user.Active {
		return c.Status(403).JSON(fiber.Map{
			"success": false,
			"message": "Account is deactivated",
		})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Invalid email or password",
		})
	}

	// Generate JWT token
	token, err := generateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate token",
		})
	}

	// Save session
	session := Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	db.Create(&session)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

func logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "No token provided",
		})
	}

	// Remove "Bearer " prefix
	if len(token) > 7 {
		token = token[7:]
	}

	// Delete session
	db.Where("token = ?", token).Delete(&Session{})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logout successful",
	})
}

func getCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func updateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	db.Save(&user)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
		"data":    user,
	})
}

func changePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Invalid old password",
		})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to hash password",
		})
	}

	user.Password = string(hashedPassword)
	db.Save(&user)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password changed successfully",
	})
}

func refreshToken(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Verify old token
	claims, err := verifyToken(req.Token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Invalid token",
		})
	}

	// Get user
	var user User
	if err := db.First(&user, claims.UserID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	// Generate new token
	newToken, err := generateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate token",
		})
	}

	// Update session
	db.Model(&Session{}).Where("token = ?", req.Token).Updates(map[string]interface{}{
		"token":      newToken,
		"expires_at": time.Now().Add(24 * time.Hour),
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"token": newToken,
		},
	})
}

// ==================== ADMIN HANDLERS ====================

func getUsers(c *fiber.Ctx) error {
	var users []User
	db.Find(&users)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    users,
	})
}

func getUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var user User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func updateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	db.Save(&user)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User updated successfully",
		"data":    user,
	})
}

func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := db.Delete(&User{}, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User deleted successfully",
	})
}

func activateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	user.Active = true
	db.Save(&user)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User activated successfully",
	})
}

func deactivateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	user.Active = false
	db.Save(&user)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User deactivated successfully",
	})
}

// ==================== ROLE HANDLERS ====================

func getRoles(c *fiber.Ctx) error {
	var roles []Role
	db.Find(&roles)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    roles,
	})
}

func createRole(c *fiber.Ctx) error {
	var role Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if err := db.Create(&role).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create role",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    role,
	})
}

func updateRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var role Role
	if err := db.First(&role, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Role not found",
		})
	}

	if err := c.BodyParser(&role); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	db.Save(&role)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    role,
	})
}

func deleteRole(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := db.Delete(&Role{}, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Role not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Role deleted successfully",
	})
}

// ==================== JWT UTILITIES ====================

func generateToken(user User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
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

// ==================== MIDDLEWARE ====================

func jwtMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "No authorization header",
		})
	}

	// Remove "Bearer " prefix
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	// Verify token
	claims, err := verifyToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Invalid or expired token",
		})
	}

	// Check if session exists
	var session Session
	if err := db.Where("token = ? AND expires_at > ?", tokenString, time.Now()).First(&session).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Session expired",
		})
	}

	// Store user info in context
	c.Locals("userID", claims.UserID)
	c.Locals("email", claims.Email)
	c.Locals("role", claims.Role)

	return c.Next()
}

func adminMiddleware(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"success": false,
			"message": "Admin access required",
		})
	}

	return c.Next()
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": err.Error(),
	})
}
