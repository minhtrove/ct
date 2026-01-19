// Package page defines HTTP handlers for the web application.
package page

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/minhtranin/ct/internal/handler"
	"github.com/minhtranin/ct/internal/models"
	"github.com/minhtranin/ct/internal/render"
	view "github.com/minhtranin/ct/internal/view/components"
	"github.com/minhtranin/ct/internal/view/layouts"
	"go.mongodb.org/mongo-driver/bson"
)

// getUser fetches user from session and database, returns user and role
func getUser(f *fiber.Ctx) (*models.User, error) {
	userID, err := handler.GetSession(f)
	if err != nil {
		return nil, err
	}

	db := handler.GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(f.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	// Set default role if empty
	if user.Role == "" {
		user.Role = "employee"
	}

	return &user, nil
}

// isHTMXRequest checks if the request is an HTMX request
func isHTMXRequest(f *fiber.Ctx) bool {
	return f.Get("HX-Request") == "true"
}

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
	user, err := getUser(f)
	if err != nil {
		return f.Redirect("/signin")
	}

	db := handler.GetDB()

	// Fetch accounts for the company
	var accounts []models.Account
	accountsCursor, _ := db.Database("ct").Collection("accounts").Find(f.Context(), bson.M{
		"company_id": user.CompanyID,
		"is_active":  true,
	})
	if accountsCursor != nil {
		accountsCursor.All(f.Context(), &accounts)
	}

	// Calculate total balance
	var totalBalance float64
	for _, acc := range accounts {
		totalBalance += acc.Balance
	}

	// Count pending approvals
	pendingCount, _ := db.Database("ct").Collection("transactions").CountDocuments(f.Context(), bson.M{
		"company_id": user.CompanyID,
		"status":     "pending",
	})

	// Fetch budgets for the company
	var budgets []models.Budget
	budgetsCursor, _ := db.Database("ct").Collection("budgets").Find(f.Context(), bson.M{
		"company_id": user.CompanyID,
		"is_active":  true,
	})
	if budgetsCursor != nil {
		budgetsCursor.All(f.Context(), &budgets)
	}

	// Build budget summaries (for now, show budget amount as limit, 0 spent)
	var budgetSummaries []view.BudgetSummary
	for _, budget := range budgets {
		summary := view.BudgetSummary{
			Name:        budget.Name,
			Spent:       budget.Spent,
			Limit:       budget.Amount,
			Utilization: 0,
			Color:       "#6366f1", // indigo
		}
		if budget.Amount > 0 {
			summary.Utilization = (budget.Spent / budget.Amount) * 100
		}
		budgetSummaries = append(budgetSummaries, summary)
	}

	// Fetch approved transactions for last 6 months for chart
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	var transactions []models.Transaction
	txnCursor, _ := db.Database("ct").Collection("transactions").Find(f.Context(), bson.M{
		"company_id":       user.CompanyID,
		"status":           models.TransactionStatusApproved,
		"transaction_date": bson.M{"$gte": sixMonthsAgo},
	})
	if txnCursor != nil {
		txnCursor.All(f.Context(), &transactions)
	}

	// Build monthly chart data
	monthlyData := make(map[string]*view.MonthlyChartPoint)
	for _, txn := range transactions {
		monthKey := txn.TransactionDate.Format("Jan 2006")
		if _, exists := monthlyData[monthKey]; !exists {
			monthlyData[monthKey] = &view.MonthlyChartPoint{Label: txn.TransactionDate.Format("Jan")}
		}
		if txn.Type == models.TransactionTypeIncome {
			monthlyData[monthKey].Income += txn.Amount
		} else if txn.Type == models.TransactionTypeExpense {
			monthlyData[monthKey].Expense += txn.Amount
		}
	}

	// Convert map to sorted slice (last 6 months)
	var monthlyChartData []view.MonthlyChartPoint
	for i := 5; i >= 0; i-- {
		month := time.Now().AddDate(0, -i, 0)
		monthKey := month.Format("Jan 2006")
		label := month.Format("Jan")
		if data, exists := monthlyData[monthKey]; exists {
			monthlyChartData = append(monthlyChartData, *data)
		} else {
			monthlyChartData = append(monthlyChartData, view.MonthlyChartPoint{Label: label, Income: 0, Expense: 0})
		}
	}

	// Build dashboard data
	data := view.DashboardData{
		User:             user,
		TotalBalance:     totalBalance,
		MonthIncome:      0, // TODO: Calculate from transactions
		MonthExpense:     0, // TODO: Calculate from transactions
		PendingCount:     int(pendingCount),
		Accounts:         accounts,
		BudgetSummaries:  budgetSummaries,
		MonthlyChartData: monthlyChartData,
	}

	// Check if HTMX request - return only content
	if isHTMXRequest(f) {
		return render.HTML(f, view.DashboardPage(data))
	}

	// Read sidebar state from cookie
	collapsed := false
	if f.Cookies("sidebar_state") == "false" {
		collapsed = true
	}

	return render.HTML(
		f,
		layouts.Dashboard("Dashboard", view.DashboardPage(data), collapsed, user.Email, user.Role, f.Path()),
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
