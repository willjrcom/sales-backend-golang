package companyentity

import (
	"os"
	"strconv"
)

type PlanType string
type Frequency string

const (
	PlanFree         PlanType = "free"
	PlanBasic        PlanType = "basic"
	PlanIntermediate PlanType = "intermediate"
	PlanAdvanced     PlanType = "advanced"

	FrequencyMonthly      Frequency = "MONTHLY"
	FrequencySemiannually Frequency = "SEMIANNUALLY"
	FrequencyAnnually     Frequency = "ANNUALLY"
)

type Plan struct {
	Key      PlanType `json:"key"`
	Name     string   `json:"name"`
	Price    float64  `json:"price"`
	Features []string `json:"features"`
	Order    int      `json:"order"`
}

func GetAllPlans() []Plan {
	return []Plan{
		{
			Key:      PlanBasic,
			Name:     "Básico",
			Price:    getEnvFloat("PRICE_BASIC", 99.90),
			Features: []string{"Gestão de Vendas", "Controle de Estoque", "Relatórios Básicos"},
			Order:    1,
		},
		{
			Key:      PlanIntermediate,
			Name:     "Intermediário (+Fiscal)",
			Price:    getEnvFloat("PRICE_INTERMEDIATE", 119.90),
			Features: []string{"Tudo do Básico", "Emissão de NF-e/NFC-e", "Até 300 notas/mês"},
			Order:    2,
		},
		{
			Key:      PlanAdvanced,
			Name:     "Avançado (+Ilimitado)",
			Price:    getEnvFloat("PRICE_ADVANCED", 129.90),
			Features: []string{"Tudo do Intermediário", "Notas Ilimitadas", "Suporte Prioritário"},
			Order:    3,
		},
	}
}

func getEnvFloat(key string, fallback float64) float64 {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return fallback
	}
	return val
}
