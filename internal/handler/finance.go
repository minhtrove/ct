package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// logAudit creates an audit log entry
func logAudit(c *fiber.Ctx, action models.AuditAction, entity models.AuditEntity, entityID primitive.ObjectID, user *models.User, changes map[string]interface{}) {
	db := GetDB()
	auditCollection := db.Database("ct").Collection("audit_logs")

	log := models.NewAuditLog(action, entity, entityID, user.ID, user.CompanyID, user.Name, user.Email)
	log.WithIPAddress(c.IP())
	log.WithUserAgent(c.Get("User-Agent"))
	if changes != nil {
		log.WithChanges(changes)
	}

	auditCollection.InsertOne(c.Context(), log)
}

// CreateAccount handles POST /api/accounts
func CreateAccount(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	// Parse form data
	name := c.FormValue("name")
	accountType := c.FormValue("type")
	currency := c.FormValue("currency")
	accountNumber := c.FormValue("account_number")
	description := c.FormValue("description")

	if name == "" || accountType == "" || currency == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, type, and currency are required"})
	}

	// Create account
	now := time.Now()
	account := models.Account{
		ID:            primitive.NewObjectID(),
		Name:          name,
		Type:          models.AccountType(accountType),
		Currency:      currency,
		AccountNumber: accountNumber,
		Description:   description,
		Balance:       0,
		CompanyID:     user.CompanyID,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	accountsCollection := db.Database("ct").Collection("accounts")
	_, err = accountsCollection.InsertOne(c.Context(), account)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create account"})
	}

	// Audit log
	logAudit(c, models.AuditActionCreate, models.AuditEntityAccount, account.ID, &user, map[string]interface{}{
		"name":     name,
		"type":     accountType,
		"currency": currency,
	})

	// Redirect back to accounts page with success toast
	return c.Redirect("/accounts?success=Account+created")
}

// CreateCategory handles POST /api/categories
func CreateCategory(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	// Parse form data
	name := c.FormValue("name")
	categoryType := c.FormValue("type")
	color := c.FormValue("color")

	if name == "" || categoryType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name and type are required"})
	}

	if color == "" {
		color = "#6366f1" // Default indigo
	}

	// Create category
	now := time.Now()
	category := models.Category{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Type:      models.CategoryType(categoryType),
		Color:     color,
		CompanyID: user.CompanyID,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	categoriesCollection := db.Database("ct").Collection("categories")
	_, err = categoriesCollection.InsertOne(c.Context(), category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create category"})
	}

	// Audit log
	logAudit(c, models.AuditActionCreate, models.AuditEntityCategory, category.ID, &user, map[string]interface{}{
		"name":  name,
		"type":  categoryType,
		"color": color,
	})

	// Redirect back to categories page with success toast
	return c.Redirect("/categories?success=Category+created")
}

// CreateBudget handles POST /api/budgets
func CreateBudget(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	// Parse form data
	name := c.FormValue("name")
	categoryID := c.FormValue("category_id")
	amountStr := c.FormValue("amount")
	period := c.FormValue("period")

	if name == "" || categoryID == "" || amountStr == "" || period == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields are required"})
	}

	catObjID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	var amount float64
	_, err = parseFloat(amountStr, &amount)
	if err != nil || amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid amount"})
	}

	// Get period dates
	now := time.Now()
	startDate, endDate := models.GetBudgetPeriodDates(models.BudgetPeriod(period), now)

	// Create budget
	budget := models.Budget{
		ID:         primitive.NewObjectID(),
		Name:       name,
		CategoryID: catObjID,
		Amount:     amount,
		Spent:      0,
		Currency:   "USD",
		Period:     models.BudgetPeriod(period),
		StartDate:  startDate,
		EndDate:    endDate,
		CompanyID:  user.CompanyID,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	budgetsCollection := db.Database("ct").Collection("budgets")
	_, err = budgetsCollection.InsertOne(c.Context(), budget)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create budget"})
	}

	// Audit log
	logAudit(c, models.AuditActionCreate, models.AuditEntityBudget, budget.ID, &user, map[string]interface{}{
		"name":        name,
		"category_id": categoryID,
		"amount":      amount,
		"period":      period,
	})

	// Redirect back to budgets page with success toast
	return c.Redirect("/budgets?success=Budget+created")
}

// CreateTransaction handles POST /api/transactions
func CreateTransaction(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	// Parse form data
	txnType := c.FormValue("type")
	amountStr := c.FormValue("amount")
	description := c.FormValue("description")
	fromAccountID := c.FormValue("from_account_id")
	toAccountID := c.FormValue("to_account_id")
	categoryID := c.FormValue("category_id")

	if txnType == "" || amountStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Type and amount are required"})
	}

	var amount float64
	_, err = parseFloat(amountStr, &amount)
	if err != nil || amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid amount"})
	}

	// Create transaction
	now := time.Now()
	txn := models.Transaction{
		ID:              primitive.NewObjectID(),
		Type:            models.TransactionType(txnType),
		Amount:          amount,
		Currency:        "USD",
		Description:     description,
		Status:          models.TransactionStatusPending,
		CreatedByID:     user.ID,
		CreatedByName:   user.Name,
		CompanyID:       user.CompanyID,
		TransactionDate: now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Set account IDs based on type
	if fromAccountID != "" {
		fromID, _ := primitive.ObjectIDFromHex(fromAccountID)
		txn.FromAccountID = fromID
	}
	if toAccountID != "" {
		toID, _ := primitive.ObjectIDFromHex(toAccountID)
		txn.ToAccountID = toID
	}
	if categoryID != "" {
		catID, _ := primitive.ObjectIDFromHex(categoryID)
		txn.CategoryID = catID
	}

	transactionsCollection := db.Database("ct").Collection("transactions")
	_, err = transactionsCollection.InsertOne(c.Context(), txn)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create transaction"})
	}

	// Audit log
	logAudit(c, models.AuditActionCreate, models.AuditEntityTransaction, txn.ID, &user, map[string]interface{}{
		"type":        txnType,
		"amount":      amount,
		"description": description,
	})

	// Redirect back to transactions page with success toast
	return c.Redirect("/transactions?success=Transaction+created")
}

// Helper function to parse float
func parseFloat(s string, f *float64) (bool, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false, err
	}
	*f = v
	return true, nil
}

// ApproveTransaction handles POST /api/transactions/:id/approve
func ApproveTransaction(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Redirect("/approvals?error=User+not+found")
	}

	txnID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Redirect("/approvals?error=Invalid+transaction+ID")
	}

	// Get transaction
	var txn models.Transaction
	txnCollection := db.Database("ct").Collection("transactions")
	err = txnCollection.FindOne(c.Context(), bson.M{"_id": txnID}).Decode(&txn)
	if err != nil {
		return c.Redirect("/approvals?error=Transaction+not+found")
	}

	// Update status
	now := time.Now()
	_, err = txnCollection.UpdateOne(c.Context(), bson.M{"_id": txnID}, bson.M{
		"$set": bson.M{
			"status":      models.TransactionStatusApproved,
			"approved_by": user.ID,
			"approved_at": now,
			"updated_at":  now,
		},
	})
	if err != nil {
		return c.Redirect("/approvals?error=Failed+to+approve")
	}

	// Update account balances
	updateAccountBalance(c, &txn)

	// Update budget spent if this is an expense with a category
	if txn.Type == models.TransactionTypeExpense && !txn.CategoryID.IsZero() {
		updateBudgetSpent(c, txn.CategoryID, txn.Amount)
	}

	// Audit log
	logAudit(c, models.AuditActionApprove, models.AuditEntityTransaction, txnID, &user, map[string]interface{}{
		"amount": txn.Amount,
		"type":   string(txn.Type),
	})

	return c.Redirect("/approvals?success=Transaction+approved")
}

// RejectTransaction handles POST /api/transactions/:id/reject
func RejectTransaction(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Redirect("/approvals?error=User+not+found")
	}

	txnID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Redirect("/approvals?error=Invalid+transaction+ID")
	}

	reason := c.FormValue("reason")

	// Get transaction
	var txn models.Transaction
	txnCollection := db.Database("ct").Collection("transactions")
	err = txnCollection.FindOne(c.Context(), bson.M{"_id": txnID}).Decode(&txn)
	if err != nil {
		return c.Redirect("/approvals?error=Transaction+not+found")
	}

	// Update status
	now := time.Now()
	_, err = txnCollection.UpdateOne(c.Context(), bson.M{"_id": txnID}, bson.M{
		"$set": bson.M{
			"status":           models.TransactionStatusRejected,
			"rejection_reason": reason,
			"approved_by":      user.ID,
			"approved_at":      now,
			"updated_at":       now,
		},
	})
	if err != nil {
		return c.Redirect("/approvals?error=Failed+to+reject")
	}

	// Audit log
	logAudit(c, models.AuditActionReject, models.AuditEntityTransaction, txnID, &user, map[string]interface{}{
		"amount": txn.Amount,
		"reason": reason,
	})

	return c.Redirect("/approvals?success=Transaction+rejected")
}

// UpdateUserRole handles POST /api/users/:id/role
func UpdateUserRole(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var currentUser models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&currentUser)
	if err != nil {
		return c.Redirect("/team?error=User+not+found")
	}

	// Only super_admin can change roles
	if currentUser.Role != "super_admin" {
		return c.Redirect("/team?error=Permission+denied")
	}

	targetUserID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Redirect("/team?error=Invalid+user+ID")
	}

	newRole := c.FormValue("role")
	if newRole == "" {
		return c.Redirect("/team?error=Role+is+required")
	}

	// Update role
	_, err = usersCollection.UpdateOne(c.Context(), bson.M{"_id": targetUserID}, bson.M{
		"$set": bson.M{
			"role":       newRole,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return c.Redirect("/team?error=Failed+to+update+role")
	}

	// Audit log
	logAudit(c, models.AuditActionUpdate, models.AuditEntityUser, targetUserID, &currentUser, map[string]interface{}{
		"new_role": newRole,
	})

	return c.Redirect("/team?success=Role+updated")
}

// UpdateAccount handles PUT /api/accounts/:id
func UpdateAccount(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Redirect("/accounts?error=User+not+found")
	}

	accountID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Redirect("/accounts?error=Invalid+account+ID")
	}

	name := c.FormValue("name")
	accountType := c.FormValue("type")
	currency := c.FormValue("currency")

	accountsCollection := db.Database("ct").Collection("accounts")
	_, err = accountsCollection.UpdateOne(c.Context(), bson.M{"_id": accountID, "company_id": user.CompanyID}, bson.M{
		"$set": bson.M{
			"name":       name,
			"type":       accountType,
			"currency":   currency,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return c.Redirect("/accounts?error=Failed+to+update+account")
	}

	logAudit(c, models.AuditActionUpdate, models.AuditEntityAccount, accountID, &user, map[string]interface{}{
		"name": name,
	})

	return c.Redirect("/accounts?success=Account+updated")
}

// DeleteAccount handles DELETE /api/accounts/:id
func DeleteAccount(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Redirect("/accounts?error=User+not+found")
	}

	accountID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Redirect("/accounts?error=Invalid+account+ID")
	}

	accountsCollection := db.Database("ct").Collection("accounts")
	_, err = accountsCollection.UpdateOne(c.Context(), bson.M{"_id": accountID, "company_id": user.CompanyID}, bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return c.Redirect("/accounts?error=Failed+to+delete+account")
	}

	logAudit(c, models.AuditActionDelete, models.AuditEntityAccount, accountID, &user, nil)

	return c.Redirect("/accounts?success=Account+deleted")
}

// UpdateCategory handles PUT /api/categories/:id
func UpdateCategory(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Redirect("/categories?error=User+not+found")
	}

	categoryID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Redirect("/categories?error=Invalid+category+ID")
	}

	name := c.FormValue("name")
	catType := c.FormValue("type")
	color := c.FormValue("color")

	categoriesCollection := db.Database("ct").Collection("categories")
	_, err = categoriesCollection.UpdateOne(c.Context(), bson.M{"_id": categoryID, "company_id": user.CompanyID}, bson.M{
		"$set": bson.M{
			"name":       name,
			"type":       catType,
			"color":      color,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return c.Redirect("/categories?error=Failed+to+update+category")
	}

	logAudit(c, models.AuditActionUpdate, models.AuditEntityCategory, categoryID, &user, map[string]interface{}{
		"name": name,
	})

	return c.Redirect("/categories?success=Category+updated")
}

// DeleteCategory handles DELETE /api/categories/:id
func DeleteCategory(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Redirect("/categories?error=User+not+found")
	}

	categoryID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Redirect("/categories?error=Invalid+category+ID")
	}

	categoriesCollection := db.Database("ct").Collection("categories")
	_, err = categoriesCollection.UpdateOne(c.Context(), bson.M{"_id": categoryID, "company_id": user.CompanyID}, bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return c.Redirect("/categories?error=Failed+to+delete+category")
	}

	logAudit(c, models.AuditActionDelete, models.AuditEntityCategory, categoryID, &user, nil)

	return c.Redirect("/categories?success=Category+deleted")
}

// UpdateBudget handles PUT /api/budgets/:id
func UpdateBudget(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Redirect("/budgets?error=User+not+found")
	}

	budgetID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Redirect("/budgets?error=Invalid+budget+ID")
	}

	name := c.FormValue("name")
	var amount float64
	parseFloat(c.FormValue("amount"), &amount)

	budgetsCollection := db.Database("ct").Collection("budgets")
	_, err = budgetsCollection.UpdateOne(c.Context(), bson.M{"_id": budgetID, "company_id": user.CompanyID}, bson.M{
		"$set": bson.M{
			"name":       name,
			"amount":     amount,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return c.Redirect("/budgets?error=Failed+to+update+budget")
	}

	logAudit(c, models.AuditActionUpdate, models.AuditEntityBudget, budgetID, &user, map[string]interface{}{
		"name":   name,
		"amount": amount,
	})

	return c.Redirect("/budgets?success=Budget+updated")
}

// DeleteBudget handles DELETE /api/budgets/:id
func DeleteBudget(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Redirect("/budgets?error=User+not+found")
	}

	budgetID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Redirect("/budgets?error=Invalid+budget+ID")
	}

	budgetsCollection := db.Database("ct").Collection("budgets")
	_, err = budgetsCollection.UpdateOne(c.Context(), bson.M{"_id": budgetID, "company_id": user.CompanyID}, bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return c.Redirect("/budgets?error=Failed+to+delete+budget")
	}

	logAudit(c, models.AuditActionDelete, models.AuditEntityBudget, budgetID, &user, nil)

	return c.Redirect("/budgets?success=Budget+deleted")
}

// UpdateTransaction handles PUT /api/transactions/:id
func UpdateTransaction(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Redirect("/transactions?error=User+not+found")
	}

	txnID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Redirect("/transactions?error=Invalid+transaction+ID")
	}

	description := c.FormValue("description")
	var amount float64
	parseFloat(c.FormValue("amount"), &amount)

	// Parse transaction date
	var txnDate time.Time
	if dateStr := c.FormValue("transaction_date"); dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			txnDate = parsed
		}
	}

	updateFields := bson.M{
		"description": description,
		"amount":      amount,
		"updated_at":  time.Now(),
	}
	if !txnDate.IsZero() {
		updateFields["transaction_date"] = txnDate
	}

	txnCollection := db.Database("ct").Collection("transactions")
	_, err = txnCollection.UpdateOne(c.Context(), bson.M{"_id": txnID, "company_id": user.CompanyID}, bson.M{
		"$set": updateFields,
	})
	if err != nil {
		return c.Redirect("/transactions?error=Failed+to+update+transaction")
	}

	logAudit(c, models.AuditActionUpdate, models.AuditEntityTransaction, txnID, &user, map[string]interface{}{
		"description": description,
		"amount":      amount,
	})

	return c.Redirect("/transactions?success=Transaction+updated")
}

// DeleteTransaction handles DELETE /api/transactions/:id
func DeleteTransaction(c *fiber.Ctx) error {
	userID, err := GetSession(c)
	if err != nil {
		return c.Redirect("/signin")
	}

	db := GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Redirect("/transactions?error=User+not+found")
	}

	txnID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Redirect("/transactions?error=Invalid+transaction+ID")
	}

	txnCollection := db.Database("ct").Collection("transactions")
	_, err = txnCollection.DeleteOne(c.Context(), bson.M{"_id": txnID, "company_id": user.CompanyID})
	if err != nil {
		return c.Redirect("/transactions?error=Failed+to+delete+transaction")
	}

	logAudit(c, models.AuditActionDelete, models.AuditEntityTransaction, txnID, &user, nil)

	return c.Redirect("/transactions?success=Transaction+deleted")
}

// updateAccountBalance updates account balance based on transaction
func updateAccountBalance(c *fiber.Ctx, txn *models.Transaction) {
	db := GetDB()
	accountsCollection := db.Database("ct").Collection("accounts")

	switch txn.Type {
	case models.TransactionTypeIncome:
		// Add to destination account
		if !txn.ToAccountID.IsZero() {
			accountsCollection.UpdateOne(c.Context(), bson.M{"_id": txn.ToAccountID}, bson.M{
				"$inc": bson.M{"balance": txn.Amount},
			})
		}
	case models.TransactionTypeExpense:
		// Subtract from source account
		if !txn.FromAccountID.IsZero() {
			accountsCollection.UpdateOne(c.Context(), bson.M{"_id": txn.FromAccountID}, bson.M{
				"$inc": bson.M{"balance": -txn.Amount},
			})
		}
	case models.TransactionTypeTransfer:
		// Subtract from source, add to destination
		if !txn.FromAccountID.IsZero() {
			accountsCollection.UpdateOne(c.Context(), bson.M{"_id": txn.FromAccountID}, bson.M{
				"$inc": bson.M{"balance": -txn.Amount},
			})
		}
		if !txn.ToAccountID.IsZero() {
			accountsCollection.UpdateOne(c.Context(), bson.M{"_id": txn.ToAccountID}, bson.M{
				"$inc": bson.M{"balance": txn.Amount},
			})
		}
	}
}

// updateBudgetSpent updates the spent amount for budgets linked to this category
func updateBudgetSpent(c *fiber.Ctx, categoryID primitive.ObjectID, amount float64) {
	db := GetDB()
	budgetsCollection := db.Database("ct").Collection("budgets")

	// Find all active budgets for this category and increment spent
	budgetsCollection.UpdateMany(c.Context(), bson.M{
		"category_id": categoryID,
		"is_active":   true,
	}, bson.M{
		"$inc": bson.M{"spent": amount},
	})
}
