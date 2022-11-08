package messages

const greeting = `Бот для учета расходов

Добавить трату: <сумма> <категория> <дата*>
* - необязательный параметр
Пример: 499.99 интернет 2022-01-01

Команды:
/start - запуск бота и инструкция
/week - недельный отчет
/month - месячный отчет
/year - годовой отчет
/currency - изменить валюту
/set_budget 12.3 - установка лимита на месяц
/show_budget - вывод текущего лимита`

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

	// Raw commands
	CommandStart       = "/start"
	CommandWeekReport  = "/week"
	CommandMonthReport = "/month"
	CommandYearReport  = "/year"
	CommandGetCurrency = "/currency"
	CommandSetBudget   = "/set_budget"
	CommandShowBudget  = "/show_budget"
)

func commandReportText(command string) string {
	switch command {
	case CommandWeekReport:
		return WeekReport
	case CommandMonthReport:
		return MonthReport
	case CommandYearReport:
		return YearReport
	default:
		return ""
	}
}
