// Package page defines HTTP handlers for the web application.
package page

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/minhtranin/ct/internal/handler"
	"github.com/minhtranin/ct/internal/models"
	view "github.com/minhtranin/ct/internal/view/components"
	"github.com/minhtranin/ct/internal/render"
	"github.com/minhtranin/ct/internal/view/layouts"
	"go.mongodb.org/mongo-driver/bson"
)

func Home(f *fiber.Ctx) error {
	// Check if user is authenticated
	_, err := handler.GetSession(f)
	if err == nil {
		// User is logged in, fetch user data and show dashboard
		return Dashboard(f)
	}
	return render.HTML(
		f,
		layouts.Base("Sign In", view.SignInPage("")),
	)
}

func SignIn(f *fiber.Ctx) error {
	// Check if user is already authenticated
	_, err := handler.GetSession(f)
	if err == nil {
		// User is logged in, redirect to dashboard
		return f.Redirect("/dashboard")
	}
	return render.HTML(
		f,
		layouts.Base("Sign In", view.SignInPage("")),
	)
}

func SignUp(f *fiber.Ctx) error {
	// Check if user is already authenticated
	_, err := handler.GetSession(f)
	if err == nil {
		// User is logged in, redirect to dashboard
		return f.Redirect("/dashboard")
	}
	return render.HTML(
		f,
		layouts.Base("Sign Up", view.SignUpPage("")),
	)
}

func Dashboard(f *fiber.Ctx) error {
	// Get user from session
	userID, err := handler.GetSession(f)
	if err != nil {
		return f.Redirect("/signin")
	}

	// Fetch user from database
	db := handler.GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(f.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return f.Redirect("/signin")
	}

	return render.HTML(
		f,
		layouts.Base("Dashboard", view.DashboardPage("", &user)),
	)
}

func ForgotPassword(f *fiber.Ctx) error {
	return render.HTML(
		f,
		layouts.Base("Forgot Password", view.ForgotPasswordPage("")),
	)
}

func ResetPassword(f *fiber.Ctx) error {
	token := f.Query("token")
	return render.HTML(
		f,
		layouts.Base("Reset Password", view.ResetPasswordPage("", token)),
	)
}

func Logout(f *fiber.Ctx) error {
	return handler.Logout(f)
}

func VerifyEmailPage(f *fiber.Ctx) error {
	email := f.Query("email")
	if email == "" {
		return f.Redirect("/signup")
	}

	// Fetch user from database
	db := handler.GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	err := usersCollection.FindOne(f.Context(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return f.Redirect("/signup")
	}

	// If already verified, redirect to dashboard
	if user.EmailVerified {
		return f.Redirect("/dashboard")
	}

	// Only send email if:
	// 1. Code is expired (>15 min) OR no code exists
	// 2. AND it's been more than 1 minute since last email
	shouldSendEmail := false
	if user.VerifyToken == "" || user.VerifyExpiresAt.Before(time.Now()) {
		// Code is expired or doesn't exist - check cooldown
		if user.LastEmailSentAt.IsZero() || time.Since(user.LastEmailSentAt) >= 1*time.Minute {
			shouldSendEmail = true
		}
	}

	if shouldSendEmail {
		handler.SendVerificationEmail(email, user.VerifyToken)

		// Update LastEmailSentAt
		usersCollection.UpdateOne(f.Context(), bson.M{"_id": user.ID}, bson.M{
			"$set": bson.M{
				"last_email_sent_at": time.Now(),
			},
		})
	}

	// Render page with LastEmailSentAt
	return render.HTML(
		f,
		layouts.Base("Verify Your Email", view.VerifyEmailPage("", email, user.LastEmailSentAt)),
	)
}
