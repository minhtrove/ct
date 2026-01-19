// Package auth provides role-based access control for the application
package auth

// Role represents a user role in the system
type Role string

const (
	RoleEmployee   Role = "employee"
	RoleHolder     Role = "holder"
	RoleAccountant Role = "accountant"
	RoleManager    Role = "manager"
	RoleAdmin      Role = "admin"
	RoleSuperAdmin Role = "super_admin"
	RoleDeveloper  Role = "developer"
)

// RoleLevel defines the hierarchy level for each role
var RoleLevel = map[Role]int{
	RoleEmployee:   1,
	RoleHolder:     2,
	RoleAccountant: 3,
	RoleManager:    4,
	RoleAdmin:      5,
	RoleSuperAdmin: 6,
	RoleDeveloper:  7,
}

// RoleDisplayNames provides human-readable names for roles
var RoleDisplayNames = map[Role]string{
	RoleEmployee:   "Employee",
	RoleHolder:     "Holder",
	RoleAccountant: "Accountant",
	RoleManager:    "Manager",
	RoleAdmin:      "Admin",
	RoleSuperAdmin: "Super Admin",
	RoleDeveloper:  "Developer",
}

// GetRoleLevel returns the hierarchy level for a role string
func GetRoleLevel(role string) int {
	if level, ok := RoleLevel[Role(role)]; ok {
		return level
	}
	return 0
}

// RoleDisplayName returns the display name for a role
func RoleDisplayName(role string) string {
	if name, ok := RoleDisplayNames[Role(role)]; ok {
		return name
	}
	return role
}

// HasPermission checks if a role has at least the required level
func HasPermission(userRole string, requiredLevel int) bool {
	return GetRoleLevel(userRole) >= requiredLevel
}

// CanSubmitExpenses checks if the role can submit expenses
func CanSubmitExpenses(role string) bool {
	r := Role(role)
	return r == RoleEmployee || r == RoleAccountant || r == RoleManager || r == RoleAdmin || r == RoleSuperAdmin || r == RoleDeveloper
}

// CanApprove checks if the role can approve/reject expenses
func CanApprove(role string) bool {
	return HasPermission(role, RoleLevel[RoleHolder])
}

// CanViewAllExpenses checks if the role can view all expenses
func CanViewAllExpenses(role string) bool {
	return HasPermission(role, RoleLevel[RoleHolder])
}

// CanGenerateReports checks if the role can generate reports
func CanGenerateReports(role string) bool {
	return HasPermission(role, RoleLevel[RoleAccountant])
}

// CanManageTeam checks if the role can manage team members
func CanManageTeam(role string) bool {
	return HasPermission(role, RoleLevel[RoleManager])
}

// CanManageAccounts checks if the role can manage financial accounts
func CanManageAccounts(role string) bool {
	return HasPermission(role, RoleLevel[RoleAdmin])
}

// CanManageCategories checks if the role can manage categories
func CanManageCategories(role string) bool {
	return HasPermission(role, RoleLevel[RoleAdmin])
}

// CanManageBudgets checks if the role can manage budgets
func CanManageBudgets(role string) bool {
	return HasPermission(role, RoleLevel[RoleAdmin])
}

// CanAccessSettings checks if the role can access settings
func CanAccessSettings(role string) bool {
	return HasPermission(role, RoleLevel[RoleAdmin])
}

// IsDeveloper checks if the role is developer
func IsDeveloper(role string) bool {
	return Role(role) == RoleDeveloper
}

// Tab represents a navigation tab
type Tab string

const (
	TabDashboard     Tab = "dashboard"
	TabTransactions  Tab = "transactions"
	TabApprovals     Tab = "approvals"
	TabReports       Tab = "reports"
	TabAudit         Tab = "audit"
	TabAccounts      Tab = "accounts"
	TabCategories    Tab = "categories"
	TabBudgets       Tab = "budgets"
	TabSettings      Tab = "settings"
	TabTeam          Tab = "team"
	TabDebug         Tab = "debug"
)

// GetVisibleTabs returns the tabs visible to a role
func GetVisibleTabs(role string) []Tab {
	tabs := []Tab{TabDashboard}

	level := GetRoleLevel(role)

	// Employee+ can see transactions
	if level >= RoleLevel[RoleEmployee] {
		tabs = append(tabs, TabTransactions)
	}

	// Holder+ can see approvals
	if level >= RoleLevel[RoleHolder] {
		tabs = append(tabs, TabApprovals)
	}

	// Accountant+ can see reports and audit
	if level >= RoleLevel[RoleAccountant] {
		tabs = append(tabs, TabReports, TabAudit)
	}

	// Manager+ can see team
	if level >= RoleLevel[RoleManager] {
		tabs = append(tabs, TabTeam)
	}

	// Admin+ can see accounts, categories, budgets, settings
	if level >= RoleLevel[RoleAdmin] {
		tabs = append(tabs, TabAccounts, TabCategories, TabBudgets, TabSettings)
	}

	// Developer can see debug
	if IsDeveloper(role) {
		tabs = append(tabs, TabDebug)
	}

	return tabs
}

// IsTabVisible checks if a tab is visible to a role
func IsTabVisible(role string, tab Tab) bool {
	tabs := GetVisibleTabs(role)
	for _, t := range tabs {
		if t == tab {
			return true
		}
	}
	return false
}

// ValidRoles returns all valid role strings
func ValidRoles() []Role {
	return []Role{
		RoleEmployee,
		RoleHolder,
		RoleAccountant,
		RoleManager,
		RoleAdmin,
		RoleSuperAdmin,
		RoleDeveloper,
	}
}

// IsValidRole checks if a role string is valid
func IsValidRole(role string) bool {
	for _, r := range ValidRoles() {
		if Role(role) == r {
			return true
		}
	}
	return false
}
