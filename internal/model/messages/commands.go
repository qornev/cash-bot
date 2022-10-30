package messages

const (
	Start       = "Start menu"
	WeekReport  = "Get week report"
	MonthReport = "Get month report"
	YearReport  = "Get year report"
	GetCurrency = "Get currencies menu"
	SetBudget   = "Set month budget"
	ShowBudget  = "Show month budget"
	AddExpense  = "Add expense"
	Unknown     = "Unknown"
)

func commandReportText(command string) string {
	switch command {
	case "/week":
		return WeekReport
	case "/month":
		return MonthReport
	case "/year":
		return YearReport
	default:
		return ""
	}
}
