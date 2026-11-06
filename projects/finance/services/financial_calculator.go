package services

import (
	"finance/models"
	"shared"
)

func AnnualizeAmount(amountMinor int64, frequency string) int64 {
	switch frequency {
	case "monthly":
		return amountMinor * 12
	case "quarterly":
		return amountMinor * 4
	case "yearly":
		return amountMinor
	case "one_time":
		return 0
	default:
		return 0
	}
}

func MonthlyEquivalent(amountMinor int64, frequency string) int64 {
	switch frequency {
	case "monthly":
		return amountMinor
	case "quarterly":
		return amountMinor / 3
	case "yearly":
		return amountMinor / 12
	case "one_time":
		return 0
	default:
		return 0
	}
}

func CalculateComponentAmount(calcType string, fixedAmountMinor int64, percentageBasisPoints int64, baseMinor int64) int64 {
	if calcType == "percentage" {
		return shared.CalculatePercentageMinor(baseMinor, percentageBasisPoints)
	}
	return fixedAmountMinor
}

func CalculateProductTotals(pricingItems []models.PricingItem, components []models.ClientProductComponent) (mrr int64, arr int64, oneTime int64, monthlyExp int64) {
	for _, item := range pricingItems {
		if item.Status != "active" {
			continue
		}
		if item.Frequency == "one_time" {
			oneTime += item.AmountMinor
		} else {
			mrr += MonthlyEquivalent(item.AmountMinor, item.Frequency)
			arr += AnnualizeAmount(item.AmountMinor, item.Frequency)
		}
	}

	for _, comp := range components {
		if comp.Status != "active" {
			continue
		}
		amt := comp.AmountMinor
		if comp.CalculationType == "percentage" {
			amt = shared.CalculatePercentageMinor(mrr, comp.Percentage)
		}

		if comp.Classification == "expense" || comp.Classification == "contra_revenue" {
			monthlyExp += MonthlyEquivalent(amt, comp.Frequency)
		}
	}

	return mrr, arr, oneTime, monthlyExp
}

func CalculatePayrollBreakdown(baseSalaryMinor int64, components []models.EmployeeCompensationComponent) (gross int64, deductions int64, net int64, employerContribution int64, totalEmployerCost int64) {
	gross = baseSalaryMinor
	for _, c := range components {
		if c.Status != "active" {
			continue
		}
		amt := c.AmountMinor
		if c.CalculationType == "percentage" {
			amt = shared.CalculatePercentageMinor(baseSalaryMinor, c.Percentage)
		}

		switch c.Classification {
		case "gross_earning":
			gross += amt
		case "payroll_deduction":
			deductions += amt
		case "employer_contribution":
			employerContribution += amt
		}
	}

	net = gross - deductions
	totalEmployerCost = gross + employerContribution
	return gross, deductions, net, employerContribution, totalEmployerCost
}
