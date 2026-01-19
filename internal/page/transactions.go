package page

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/auth"
	"github.com/minhtranin/ct/internal/handler"
	"github.com/minhtranin/ct/internal/models"
	"github.com/minhtranin/ct/internal/render"
	view "github.com/minhtranin/ct/internal/view/components"
	"github.com/minhtranin/ct/internal/view/layouts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const pageSize = 20

// TransactionsPage handles GET /transactions
func TransactionsPage(c *fiber.Ctx) error {
	user, err := getUser(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	db := handler.GetDB()

	// Fetch accounts for dropdown
	var accounts []models.Account
	accountsCursor, _ := db.Database("ct").Collection("accounts").Find(c.Context(), bson.M{
		"company_id": user.CompanyID,
		"is_active":  true,
	})
	if accountsCursor != nil {
		accountsCursor.All(c.Context(), &accounts)
	}

	// Fetch categories for dropdown
	var categories []models.Category
	catCursor, _ := db.Database("ct").Collection("categories").Find(c.Context(), bson.M{
		"company_id": user.CompanyID,
		"is_active":  true,
	})
	if catCursor != nil {
		catCursor.All(c.Context(), &categories)
	}

	// Fetch transactions with filters
	var transactions []models.Transaction
	filter := bson.M{"company_id": user.CompanyID}

	if filterType := c.Query("type"); filterType != "" {
		filter["type"] = filterType
	}
	if filterStatus := c.Query("status"); filterStatus != "" {
		filter["status"] = filterStatus
	}

	// Date range filtering
	if fromDate := c.Query("from"); fromDate != "" {
		if t, err := time.Parse("2006-01-02", fromDate); err == nil {
			filter["transaction_date"] = bson.M{"$gte": t}
		}
	}
	if toDate := c.Query("to"); toDate != "" {
		if t, err := time.Parse("2006-01-02", toDate); err == nil {
			t = t.Add(24 * time.Hour) // Include the entire day
			if existing, ok := filter["transaction_date"].(bson.M); ok {
				existing["$lt"] = t
			} else {
				filter["transaction_date"] = bson.M{"$lt": t}
			}
		}
	}

	// Pagination
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	skip := int64((page - 1) * pageSize)

	// Count total
	totalCount, _ := db.Database("ct").Collection("transactions").CountDocuments(c.Context(), filter)
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(skip).SetLimit(int64(pageSize))
	txnCursor, _ := db.Database("ct").Collection("transactions").Find(c.Context(), filter, opts)
	if txnCursor != nil {
		txnCursor.All(c.Context(), &transactions)
	}

	data := view.TransactionsData{
		Transactions: transactions,
		Accounts:     accounts,
		Categories:   categories,
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalCount:   int(totalCount),
		FilterType:   c.Query("type"),
		FilterStatus: c.Query("status"),
		FilterFrom:   c.Query("from"),
		FilterTo:     c.Query("to"),
		CanCreate:    auth.CanSubmitExpenses(user.Role),
	}

	if isHTMXRequest(c) {
		return render.HTML(c, view.TransactionsPage(data))
	}

	return render.HTML(c, layouts.Dashboard("Transactions", view.TransactionsPage(data), false, user.Email, user.Role, c.Path()))
}

// AccountsPage handles GET /accounts
func AccountsPage(c *fiber.Ctx) error {
	user, err := getUser(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	if !auth.CanManageAccounts(user.Role) {
		return c.Redirect("/dashboard")
	}

	db := handler.GetDB()

	var accounts []models.Account
	cursor, _ := db.Database("ct").Collection("accounts").Find(c.Context(), bson.M{
		"company_id": user.CompanyID,
	})
	if cursor != nil {
		cursor.All(c.Context(), &accounts)
	}

	var totalBalance float64
	for _, acc := range accounts {
		if acc.IsActive {
			totalBalance += acc.Balance
		}
	}

	data := view.AccountsData{
		Accounts:     accounts,
		TotalBalance: totalBalance,
	}

	if isHTMXRequest(c) {
		return render.HTML(c, view.AccountsPage(data))
	}

	return render.HTML(c, layouts.Dashboard("Accounts", view.AccountsPage(data), false, user.Email, user.Role, c.Path()))
}

// CategoriesPage handles GET /categories
func CategoriesPage(c *fiber.Ctx) error {
	user, err := getUser(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	if !auth.CanManageCategories(user.Role) {
		return c.Redirect("/dashboard")
	}

	db := handler.GetDB()

	var incomeCategories []models.Category
	var expenseCategories []models.Category

	cursor, _ := db.Database("ct").Collection("categories").Find(c.Context(), bson.M{
		"company_id": user.CompanyID,
		"is_active":  true,
	})
	if cursor != nil {
		var allCategories []models.Category
		cursor.All(c.Context(), &allCategories)
		for _, cat := range allCategories {
			if cat.Type == models.CategoryTypeIncome {
				incomeCategories = append(incomeCategories, cat)
			} else {
				expenseCategories = append(expenseCategories, cat)
			}
		}
	}

	data := view.CategoriesData{
		IncomeCategories:  incomeCategories,
		ExpenseCategories: expenseCategories,
	}

	if isHTMXRequest(c) {
		return render.HTML(c, view.CategoriesPage(data))
	}

	return render.HTML(c, layouts.Dashboard("Categories", view.CategoriesPage(data), false, user.Email, user.Role, c.Path()))
}

// BudgetsPage handles GET /budgets
func BudgetsPage(c *fiber.Ctx) error {
	user, err := getUser(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	if !auth.CanManageBudgets(user.Role) {
		return c.Redirect("/dashboard")
	}

	db := handler.GetDB()

	var budgets []models.Budget
	cursor, _ := db.Database("ct").Collection("budgets").Find(c.Context(), bson.M{
		"company_id": user.CompanyID,
		"is_active":  true,
	})
	if cursor != nil {
		cursor.All(c.Context(), &budgets)
	}

	var categories []models.Category
	catCursor, _ := db.Database("ct").Collection("categories").Find(c.Context(), bson.M{
		"company_id": user.CompanyID,
		"is_active":  true,
	})
	if catCursor != nil {
		catCursor.All(c.Context(), &categories)
	}

	data := view.BudgetsData{
		Budgets:    budgets,
		Categories: categories,
	}

	if isHTMXRequest(c) {
		return render.HTML(c, view.BudgetsPage(data))
	}

	return render.HTML(c, layouts.Dashboard("Budgets", view.BudgetsPage(data), false, user.Email, user.Role, c.Path()))
}

// ReportsPage handles GET /reports
func ReportsPage(c *fiber.Ctx) error {
	user, err := getUser(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	if !auth.CanGenerateReports(user.Role) {
		return c.Redirect("/dashboard")
	}

	db := handler.GetDB()

	// Fetch approved transactions
	var transactions []models.Transaction
	cursor, _ := db.Database("ct").Collection("transactions").Find(c.Context(), bson.M{
		"company_id": user.CompanyID,
		"status":     models.TransactionStatusApproved,
	})
	if cursor != nil {
		cursor.All(c.Context(), &transactions)
	}

	// Fetch categories for mapping
	var categories []models.Category
	catCursor, _ := db.Database("ct").Collection("categories").Find(c.Context(), bson.M{
		"company_id": user.CompanyID,
	})
	if catCursor != nil {
		catCursor.All(c.Context(), &categories)
	}

	// Build category map
	categoryMap := make(map[string]models.Category)
	for _, cat := range categories {
		categoryMap[cat.ID.Hex()] = cat
	}

	var totalIncome, totalExpense float64
	incomeByCategory := make(map[string]float64)
	expenseByCategory := make(map[string]float64)
	monthlyData := make(map[string]view.MonthlyAmount)

	for _, txn := range transactions {
		monthKey := txn.TransactionDate.Format("Jan 2006")
		monthly := monthlyData[monthKey]
		monthly.Month = monthKey

		if txn.Type == models.TransactionTypeIncome {
			totalIncome += txn.Amount
			monthly.Income += txn.Amount
			catName := "Other"
			if cat, ok := categoryMap[txn.CategoryID.Hex()]; ok {
				catName = cat.Name
			}
			incomeByCategory[catName] += txn.Amount
		} else if txn.Type == models.TransactionTypeExpense {
			totalExpense += txn.Amount
			monthly.Expense += txn.Amount
			catName := "Other"
			if cat, ok := categoryMap[txn.CategoryID.Hex()]; ok {
				catName = cat.Name
			}
			expenseByCategory[catName] += txn.Amount
		}

		monthlyData[monthKey] = monthly
	}

	// Convert maps to slices for charts
	var incomeCatData []view.CategoryAmount
	colors := []string{"#22c55e", "#10b981", "#14b8a6", "#06b6d4", "#0ea5e9", "#3b82f6"}
	i := 0
	for name, amount := range incomeByCategory {
		incomeCatData = append(incomeCatData, view.CategoryAmount{
			Name:   name,
			Amount: amount,
			Color:  colors[i%len(colors)],
		})
		i++
	}

	var expenseCatData []view.CategoryAmount
	expColors := []string{"#ef4444", "#f97316", "#eab308", "#f59e0b", "#ec4899", "#8b5cf6"}
	i = 0
	for name, amount := range expenseByCategory {
		expenseCatData = append(expenseCatData, view.CategoryAmount{
			Name:   name,
			Amount: amount,
			Color:  expColors[i%len(expColors)],
		})
		i++
	}

	var monthlyDataSlice []view.MonthlyAmount
	for _, data := range monthlyData {
		monthlyDataSlice = append(monthlyDataSlice, data)
	}

	data := view.ReportsData{
		TotalIncome:       totalIncome,
		TotalExpense:      totalExpense,
		NetProfit:         totalIncome - totalExpense,
		PeriodLabel:       "All Time",
		IncomeByCategory:  incomeCatData,
		ExpenseByCategory: expenseCatData,
		MonthlyData:       monthlyDataSlice,
	}

	if isHTMXRequest(c) {
		return render.HTML(c, view.ReportsPage(data))
	}

	return render.HTML(c, layouts.Dashboard("Reports", view.ReportsPage(data), false, user.Email, user.Role, c.Path()))
}

// AuditPage handles GET /audit
func AuditPage(c *fiber.Ctx) error {
	user, err := getUser(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	if !auth.CanGenerateReports(user.Role) {
		return c.Redirect("/dashboard")
	}

	db := handler.GetDB()

	var logs []models.AuditLog
	filter := bson.M{"company_id": user.CompanyID}

	if filterEntity := c.Query("entity"); filterEntity != "" {
		filter["entity"] = filterEntity
	}
	if filterAction := c.Query("action"); filterAction != "" {
		filter["action"] = filterAction
	}

	// Date range filtering
	if fromDate := c.Query("from"); fromDate != "" {
		if t, err := time.Parse("2006-01-02", fromDate); err == nil {
			filter["created_at"] = bson.M{"$gte": t}
		}
	}
	if toDate := c.Query("to"); toDate != "" {
		if t, err := time.Parse("2006-01-02", toDate); err == nil {
			t = t.Add(24 * time.Hour)
			if existing, ok := filter["created_at"].(bson.M); ok {
				existing["$lt"] = t
			} else {
				filter["created_at"] = bson.M{"$lt": t}
			}
		}
	}

	// Pagination
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	skip := int64((page - 1) * pageSize)

	// Count total
	totalCount, _ := db.Database("ct").Collection("audit_logs").CountDocuments(c.Context(), filter)
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(skip).SetLimit(int64(pageSize))
	cursor, _ := db.Database("ct").Collection("audit_logs").Find(c.Context(), filter, opts)
	if cursor != nil {
		cursor.All(c.Context(), &logs)
	}

	data := view.AuditData{
		Logs:         logs,
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalCount:   int(totalCount),
		FilterEntity: c.Query("entity"),
		FilterAction: c.Query("action"),
		FilterFrom:   c.Query("from"),
		FilterTo:     c.Query("to"),
	}

	if isHTMXRequest(c) {
		return render.HTML(c, view.AuditPage(data))
	}

	return render.HTML(c, layouts.Dashboard("Audit Log", view.AuditPage(data), false, user.Email, user.Role, c.Path()))
}

// ApprovalsPage handles GET /approvals
func ApprovalsPage(c *fiber.Ctx) error {
	user, err := getUser(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	if !auth.CanApprove(user.Role) {
		return c.Redirect("/dashboard")
	}

	db := handler.GetDB()

	// Fetch pending transactions
	var pendingTransactions []models.Transaction
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, _ := db.Database("ct").Collection("transactions").Find(c.Context(), bson.M{
		"company_id": user.CompanyID,
		"status":     models.TransactionStatusPending,
	}, opts)
	if cursor != nil {
		cursor.All(c.Context(), &pendingTransactions)
	}

	// Fetch accounts
	var accounts []models.Account
	accCursor, _ := db.Database("ct").Collection("accounts").Find(c.Context(), bson.M{
		"company_id": user.CompanyID,
	})
	if accCursor != nil {
		accCursor.All(c.Context(), &accounts)
	}

	data := view.ApprovalsData{
		PendingTransactions: pendingTransactions,
		Accounts:            accounts,
	}

	if isHTMXRequest(c) {
		return render.HTML(c, view.ApprovalsPage(data))
	}

	return render.HTML(c, layouts.Dashboard("Approvals", view.ApprovalsPage(data), false, user.Email, user.Role, c.Path()))
}

// SettingsPage handles GET /settings
func SettingsPage(c *fiber.Ctx) error {
	user, err := getUser(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	if !auth.CanAccessSettings(user.Role) {
		return c.Redirect("/dashboard")
	}

	content := view.PlaceholderPage("Settings", "Configure company and system settings")

	if isHTMXRequest(c) {
		return render.HTML(c, content)
	}

	return render.HTML(c, layouts.Dashboard("Settings", content, false, user.Email, user.Role, c.Path()))
}

// TeamPage handles GET /team
func TeamPage(c *fiber.Ctx) error {
	user, err := getUser(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	if !auth.CanManageTeam(user.Role) {
		return c.Redirect("/dashboard")
	}

	db := handler.GetDB()

	// Fetch all users in company
	var users []models.User
	cursor, _ := db.Database("ct").Collection("users").Find(c.Context(), bson.M{
		"company_id": user.CompanyID,
	})
	if cursor != nil {
		cursor.All(c.Context(), &users)
	}

	data := view.TeamData{
		Users:        users,
		CurrentUser:  user,
		IsSuperAdmin: user.Role == string(auth.RoleSuperAdmin),
	}

	if isHTMXRequest(c) {
		return render.HTML(c, view.TeamPage(data))
	}

	return render.HTML(c, layouts.Dashboard("Team", view.TeamPage(data), false, user.Email, user.Role, c.Path()))
}

// Helper to get user by ID
func getUserByID(c *fiber.Ctx, id string) (*models.User, error) {
	db := handler.GetDB()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = db.Database("ct").Collection("users").FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	return &user, err
}
