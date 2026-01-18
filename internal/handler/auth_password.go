// Package handler handles HTTP requests for authentication
package handler

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/auth"
	"github.com/minhtranin/ct/internal/email"
	"github.com/minhtranin/ct/internal/logger"
	"github.com/minhtranin/ct/internal/models"
	"github.com/minhtranin/ct/internal/view/shared/toast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db             *mongo.Client
	sesClient      *email.SESClient
	sessionManager *auth.SessionManager
)

// InitAuth initializes the auth handlers
func InitAuth(database *mongo.Client) {
	db = database
	var err error
	sesClient, err = email.NewSESClient()
	if err != nil {
		logger.Error("Auth", "Failed to create SES client: "+err.Error())
	}
	sessionManager = auth.NewSessionManager()
}

// generateVerificationCode generates a 6-digit verification code
func generateVerificationCode() string {
	code := rand.Intn(900000) + 100000 // 6-digit number between 100000-999999
	return fmt.Sprintf("%06d", code)
}

// SignUp handles POST /api/auth/signup
func SignUp(c *fiber.Ctx) error {
	email := strings.TrimSpace(c.FormValue("email"))
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm-password")

	// Validate input
	if email == "" || password == "" {
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Email and password are required",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}

	if password != confirmPassword {
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Passwords do not match",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}

	if len(password) < 8 {
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Password must be at least 8 characters",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}

	usersCollection := db.Database("ct").Collection("users")

	// Check if user already exists
	var existingUser models.User
	err := usersCollection.FindOne(c.Context(), bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Email already registered",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		logger.Error("Auth", "Failed to hash password: "+err.Error())
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Failed to create account",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}

	// Generate 6-digit verification code
	verifyCode := generateVerificationCode()
	expiresAt := time.Now().Add(15 * time.Minute) // Code expires in 15 minutes

	// Create user (unverified)
	now := time.Now()
	user := models.User{
		Email:           email,
		PasswordHash:    hashedPassword,
		CreatedAt:       now,
		UpdatedAt:       now,
		LastLoginAt:     now,
		EmailVerified:   false,
		VerifyToken:     verifyCode,
		VerifyExpiresAt: expiresAt,
		LastEmailSentAt: now,
	}

	_, err = usersCollection.InsertOne(c.Context(), user)
	if err != nil {
		logger.Error("Auth", "Failed to create user: "+err.Error())
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Failed to create account",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}

	logger.Info("Auth", "User created: "+email)

	// Send verification email
	if sesClient != nil {
		go sendVerificationEmail(email, verifyCode)
	}

	// Redirect to verify email page
	return c.Redirect("/verify-email?email=" + email)
}

// SignIn handles POST /api/auth/signin
func SignIn(c *fiber.Ctx) error {
	email := strings.TrimSpace(c.FormValue("email"))
	password := c.FormValue("password")

	// Validate input
	if email == "" || password == "" {
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Email and password are required",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}

	usersCollection := db.Database("ct").Collection("users")

	// Find user
	var user models.User
	err := usersCollection.FindOne(c.Context(), bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Invalid email or password",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}
	if err != nil {
		logger.Error("Auth", "Database error: "+err.Error())
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Login failed",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}

	// Check if email is verified
	if !user.EmailVerified {
		// Set HTMX redirect header for client-side redirect to verify email page
		c.Set("HX-Redirect", "/verify-email?email="+email)
		return c.SendStatus(fiber.StatusOK)
	}

	// Verify password
	if !auth.VerifyPassword(user.PasswordHash, password) {
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Invalid email or password",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}

	// Update last login
	usersCollection.UpdateOne(c.Context(), bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{"last_login_at": time.Now()},
	})

	logger.Info("Auth", "User signed in: "+email)

	// Set session
	sessionManager.SetSession(c, &user)

	// Set HTMX redirect header for client-side redirect
	c.Set("HX-Redirect", "/dashboard")
	return c.SendStatus(fiber.StatusOK)
}

// VerifyCode handles POST /api/auth/verify-code
func VerifyCode(c *fiber.Ctx) error {
	email := strings.TrimSpace(c.FormValue("email"))
	code := strings.TrimSpace(c.FormValue("code"))

	if email == "" || code == "" {
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Email and code are required",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}

	usersCollection := db.Database("ct").Collection("users")

	// Find user with valid verification code
	var user models.User
	err := usersCollection.FindOne(c.Context(), bson.M{
		"email":             email,
		"verify_token":       code,
		"verify_expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Invalid or expired verification code",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}
	if err != nil {
		logger.Error("Auth", "Database error: "+err.Error())
		c.Set("Content-Type", "text/html")
		return toast.Toast(toast.Props{
			Title:         "Failed to verify email",
			Variant:       toast.VariantError,
			Position:      toast.PositionTopLeft,
			Duration:      5000,
			Dismissible:   true,
			ShowIndicator: true,
			Icon:          true,
		}).Render(c.Context(), c.Response().BodyWriter())
	}

	// Update user as verified and clear token
	usersCollection.UpdateOne(c.Context(), bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"email_verified":    true,
			"verify_token":      "",
			"verify_expires_at": time.Time{},
			"updated_at":        time.Now(),
		},
	})

	logger.Info("Auth", "Email verified with code: "+email)

	// Set session and redirect to dashboard
	sessionManager.SetSession(c, &user)

	// Set HTMX redirect header for client-side redirect
	c.Set("HX-Redirect", "/dashboard")
	return c.SendStatus(fiber.StatusOK)
}

// ResendVerification handles POST /api/auth/resend-verification
func ResendVerification(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		email = strings.TrimSpace(c.FormValue("email"))
	}

	if email == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Email is required")
	}

	usersCollection := db.Database("ct").Collection("users")

	// Find user
	var user models.User
	err := usersCollection.FindOne(c.Context(), bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}
	if err != nil {
		logger.Error("Auth", "Database error: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to resend verification")
	}

	// Check 1-minute cooldown
	if !user.LastEmailSentAt.IsZero() {
		timeSinceLastSent := time.Since(user.LastEmailSentAt)
		if timeSinceLastSent < 1*time.Minute {
			c.Set("Content-Type", "text/html")
			return toast.Toast(toast.Props{
				Title:         "Please wait before resending",
				Variant:       toast.VariantWarning,
				Position:      toast.PositionTopLeft,
				Description:   fmt.Sprintf("You can resend in %d seconds", 60-int(timeSinceLastSent.Seconds())),
				Duration:      3000,
				Dismissible:   true,
				ShowIndicator: true,
				Icon:          true,
			}).Render(c.Context(), c.Response().BodyWriter())
		}
	}

	// Generate new verification code
	newCode := generateVerificationCode()
	expiresAt := time.Now().Add(15 * time.Minute)

	// Update verification code and LastEmailSentAt
	usersCollection.UpdateOne(c.Context(), bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"verify_token":       newCode,
			"verify_expires_at":  expiresAt,
			"last_email_sent_at": time.Now(),
		},
	})

	logger.Info("Auth", "Verification code resent for: "+email)

	// Send verification email
	if sesClient != nil {
		go sendVerificationEmail(email, newCode)
	}

	// Return success response
	c.Set("Content-Type", "text/html")
	return c.SendString(`<span class="text-green-600">Verification code sent!</span>`)
}

// ForgotPassword handles POST /api/auth/forgot-password
func ForgotPassword(c *fiber.Ctx) error {
	email := strings.TrimSpace(c.FormValue("email"))

	if email == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Email is required")
	}

	usersCollection := db.Database("ct").Collection("users")

	// Find user
	var user models.User
	err := usersCollection.FindOne(c.Context(), bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		// Don't reveal if email exists or not
		return c.Redirect("/signin")
	}
	if err != nil {
		logger.Error("Auth", "Database error: "+err.Error())
		return c.Redirect("/signin")
	}

	// Generate reset token
	resetToken, err := auth.GenerateToken()
	if err != nil {
		logger.Error("Auth", "Failed to generate reset token: "+err.Error())
		return c.Redirect("/signin")
	}

	// Update user with reset token
	usersCollection.UpdateOne(c.Context(), bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"reset_token":      resetToken,
			"reset_expires_at": time.Now().Add(1 * time.Hour),
		},
	})

	logger.Info("Auth", "Password reset requested for: "+email)

	// Send reset email
	if sesClient != nil {
		go sendPasswordResetEmail(email, resetToken)
	}

	// Set HTMX redirect header for client-side redirect
	c.Set("HX-Redirect", "/signin")
	return c.SendStatus(fiber.StatusOK)
}

// ResetPassword handles POST /api/auth/reset-password
func ResetPassword(c *fiber.Ctx) error {
	token := c.FormValue("token")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm-password")

	if token == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Token and password are required")
	}

	if password != confirmPassword {
		return c.Status(fiber.StatusBadRequest).SendString("Passwords do not match")
	}

	if len(password) < 8 {
		return c.Status(fiber.StatusBadRequest).SendString("Password must be at least 8 characters")
	}

	usersCollection := db.Database("ct").Collection("users")

	// Find user with valid reset token
	var user models.User
	err := usersCollection.FindOne(c.Context(), bson.M{
		"reset_token":       token,
		"reset_expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid or expired reset token")
	}
	if err != nil {
		logger.Error("Auth", "Database error: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to reset password")
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		logger.Error("Auth", "Failed to hash password: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to reset password")
	}

	// Update password and clear reset token
	usersCollection.UpdateOne(c.Context(), bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"password_hash":    hashedPassword,
			"reset_token":       "",
			"reset_expires_at":  time.Time{},
			"updated_at":        time.Now(),
		},
	})

	logger.Info("Auth", "Password reset for: "+user.Email)

	// Set HTMX redirect header for client-side redirect
	c.Set("HX-Redirect", "/signin")
	return c.SendStatus(fiber.StatusOK)
}

// VerifyEmail handles GET /api/auth/verify-email (old method, keeping for compatibility)
func VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")

	if token == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Token is required")
	}

	usersCollection := db.Database("ct").Collection("users")

	// Find user with valid verification token
	var user models.User
	err := usersCollection.FindOne(c.Context(), bson.M{
		"verify_token":       token,
		"verify_expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid or expired verification token")
	}
	if err != nil {
		logger.Error("Auth", "Database error: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to verify email")
	}

	// Update user as verified and clear token
	usersCollection.UpdateOne(c.Context(), bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"email_verified":    true,
			"verify_token":      "",
			"verify_expires_at": time.Time{},
			"updated_at":        time.Now(),
		},
	})

	logger.Info("Auth", "Email verified: "+user.Email)

	return c.Redirect("/dashboard")
}

// Logout handles GET /logout
func Logout(c *fiber.Ctx) error {
	sessionManager.ClearSession(c)
	return c.Redirect("/signin")
}

// sendVerificationEmail sends verification email with 6-digit code
func sendVerificationEmail(toEmail, code string) {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "CT"
	}

	subject := "Verify your email for " + appName

	// HTML email
	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #5D5CFF;">Verify your email</h2>
			<p>Thank you for signing up for %s!</p>
			<p>Your verification code is:</p>
			<div style="background-color: #f5f5f5; padding: 20px; border-radius: 8px; text-align: center; margin: 30px 0;">
				<span style="font-size: 32px; font-weight: bold; letter-spacing: 4px; color: #5D5CFF;">%s</span>
			</div>
			<p style="color: #666; font-size: 14px;">This code expires in 15 minutes.</p>
			<p style="color: #666; font-size: 14px;">If you didn't create an account, please ignore this email.</p>
		</div>
	`, appName, code)

	// Plain text email
	textBody := fmt.Sprintf(
		"Thank you for signing up for %s!\n\nYour verification code is: %s\n\nThis code expires in 15 minutes.\n\nIf you didn't create an account, please ignore this email.",
		appName, code,
	)

	// Create the email input
	input := &email.SendEmailInput{
		ToEmailAddress:   toEmail,
		FromEmailAddress: sesClient.FormatFromAddress(),
		Content: &email.EmailContent{
			Simple: &email.Message{
				Subject: &email.Content{
					Data:    subject,
					Charset: "UTF-8",
				},
				Body: &email.Body{
					Html: &email.Content{
						Data:    htmlBody,
						Charset: "UTF-8",
					},
					Text: &email.Content{
						Data:    textBody,
						Charset: "UTF-8",
					},
				},
			},
		},
	}

	// Send the email
	logger.Info("Email", "Sending verification email to: "+toEmail)
	result, err := sesClient.SendEmail(input)
	if err != nil {
		logger.Error("Email", "Failed to send verification email via SES: "+err.Error())
		return
	}

	logger.Info("Email", "Verification email sent successfully to: "+toEmail+", MessageID: "+*result.MessageId)
}

// sendPasswordResetEmail sends password reset email
func sendPasswordResetEmail(toEmail, token string) {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "CT"
	}

	baseURL := sesClient.GetBaseURL()
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)
	subject := "Reset your password for " + appName

	// HTML email
	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #5D5CFF;">Reset your password</h2>
			<p>We received a request to reset your password for your %s account.</p>
			<p>Click the button below to reset your password:</p>
			<div style="margin: 30px 0;">
				<a href="%s" style="background-color: #5D5CFF; color: white; padding: 12px 30px; text-decoration: none; border-radius: 8px; display: inline-block; font-weight: bold;">Reset Password</a>
			</div>
			<p style="color: #666; font-size: 14px;">This link expires in 1 hour.</p>
			<p style="color: #666; font-size: 14px;">If you didn't request this, please ignore this email.</p>
			<p style="color: #999; font-size: 12px; margin-top: 30px;">If the button doesn't work, copy and paste this link into your browser:<br>%s</p>
		</div>
	`, appName, resetLink, resetLink)

	// Plain text email
	textBody := fmt.Sprintf(
		"We received a request to reset your password for your %s account.\n\nClick the link below to reset your password:\n\n%s\n\nThis link expires in 1 hour.\n\nIf you didn't request this, please ignore this email.",
		appName, resetLink,
	)

	// Create the email input
	input := &email.SendEmailInput{
		ToEmailAddress:   toEmail,
		FromEmailAddress: sesClient.FormatFromAddress(),
		Content: &email.EmailContent{
			Simple: &email.Message{
				Subject: &email.Content{
					Data:    subject,
					Charset: "UTF-8",
				},
				Body: &email.Body{
					Html: &email.Content{
						Data:    htmlBody,
						Charset: "UTF-8",
					},
					Text: &email.Content{
						Data:    textBody,
						Charset: "UTF-8",
					},
				},
			},
		},
	}

	// Send the email
	logger.Info("Email", "Sending password reset email to: "+toEmail)
	result, err := sesClient.SendEmail(input)
	if err != nil {
		logger.Error("Email", "Failed to send password reset email via SES: "+err.Error())
		return
	}

	logger.Info("Email", "Password reset email sent successfully to: "+toEmail+", MessageID: "+*result.MessageId)
}

// GetSession returns the user ID from session (for use in page handlers)
func GetSession(c *fiber.Ctx) (string, error) {
	userID, err := sessionManager.GetSession(c)
	if err != nil {
		return "", err
	}
	return userID.Hex(), nil
}

// GetDB returns the database instance (for use in page handlers)
func GetDB() *mongo.Client {
	return db
}

// SendVerificationEmail sends a verification email (exported for use in page handlers)
func SendVerificationEmail(toEmail, code string) {
	if sesClient != nil {
		go sendVerificationEmail(toEmail, code)
	}
}
