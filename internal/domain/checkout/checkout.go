package billingentity

type PlanType string
type Periodicity string

const (
	PlanBasic        PlanType = "BASIC"
	PlanIntermediate PlanType = "INTERMEDIATE"
	PlanAdvanced     PlanType = "ADVANCED"

	PeriodicityMonthly    Periodicity = "MONTHLY"
	PeriodicitySemiannual Periodicity = "SEMIANNUAL"
	PeriodicityAnnual     Periodicity = "ANNUAL"
)
