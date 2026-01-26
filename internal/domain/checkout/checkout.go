package billingentity

type PlanType string
type Periodicity string

const (
	PlanBasic        PlanType = "BASIC"
	PlanIntermediate PlanType = "INTERMEDIATE"
	PlanEnterprise   PlanType = "ENTERPRISE"

	PeriodicityMonthly    Periodicity = "MONTHLY"
	PeriodicitySemiannual Periodicity = "SEMIANNUAL"
	PeriodicityAnnual     Periodicity = "ANNUAL"
)
