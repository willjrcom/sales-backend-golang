package mercadopagoservice

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NewSubscriptionExternalRef(companyID string, planType string, frequency int, paymentID string) string {
	now := time.Now().UTC()
	return fmt.Sprintf("SUB:%s:%d:%s:%d:%s:%d:%s", companyID, now.Day(), now.Month().String(), now.Year(), planType, frequency, paymentID)
}

type SubscriptionExternalRef struct {
	CompanyID string
	Day       int
	Month     string
	Year      int
	PlanType  string
	Frequency int
	PaymentID string
}

func ExtractSubscriptionExternalRef(externalRef string) (SubscriptionExternalRef, error) {
	parts := strings.Split(externalRef, ":")

	if len(parts) != 8 {
		return SubscriptionExternalRef{}, errors.New("invalid external ref format")
	}

	if parts[0] != "SUB" {
		return SubscriptionExternalRef{}, errors.New("invalid external ref format")
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return SubscriptionExternalRef{}, errors.New("invalid external ref format")
	}

	year, err := strconv.Atoi(parts[4])
	if err != nil {
		return SubscriptionExternalRef{}, errors.New("invalid external ref format")
	}

	frequency, err := strconv.Atoi(parts[6])
	if err != nil {
		return SubscriptionExternalRef{}, errors.New("invalid external ref format")
	}

	return SubscriptionExternalRef{
		CompanyID: parts[1],
		Day:       day,
		Month:     parts[3],
		Year:      year,
		PlanType:  parts[5],
		Frequency: frequency,
		PaymentID: parts[7],
	}, nil
}

// ----------------------------------------------------------------------------------------------
func NewSubscriptionUpgradeExternalRef(companyID string, planType string, newAmount float64, paymentID string) string {
	now := time.Now().UTC()
	return fmt.Sprintf("SUB_UP:%s:%d:%s:%d:%s:%f:%s", companyID, now.Day(), now.Month().String(), now.Year(), planType, newAmount, paymentID)
}

type SubscriptionUpgradeExternalRef struct {
	CompanyID string
	Day       int
	Month     string
	Year      int
	PlanType  string
	NewAmount float64
	PaymentID string
}

func ExtractSubscriptionUpgradeExternalRef(externalRef string) (SubscriptionUpgradeExternalRef, error) {
	parts := strings.Split(externalRef, ":")

	if len(parts) != 8 {
		return SubscriptionUpgradeExternalRef{}, errors.New("invalid external ref format")
	}

	if parts[0] != "SUB_UP" {
		return SubscriptionUpgradeExternalRef{}, errors.New("invalid external ref format")
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return SubscriptionUpgradeExternalRef{}, errors.New("invalid external ref format")
	}

	year, err := strconv.Atoi(parts[4])
	if err != nil {
		return SubscriptionUpgradeExternalRef{}, errors.New("invalid external ref format")
	}

	newAmount, err := strconv.ParseFloat(parts[6], 64)
	if err != nil {
		return SubscriptionUpgradeExternalRef{}, errors.New("invalid external ref format")
	}

	return SubscriptionUpgradeExternalRef{
		CompanyID: parts[1],
		Day:       day,
		Month:     parts[3],
		Year:      year,
		PlanType:  parts[5],
		NewAmount: newAmount,
		PaymentID: parts[7],
	}, nil
}

func NewCostExternalRef(companyID string, paymentID string) string {
	now := time.Now().UTC()
	return fmt.Sprintf("COST:%s:%d:%s:%d:%s", companyID, now.Day(), now.Month().String(), now.Year(), paymentID)
}

type CostExternalRef struct {
	CompanyID string
	Day       int
	Month     string
	Year      int
	PaymentID string
}

func ExtractCostExternalRef(externalRef string) (CostExternalRef, error) {
	parts := strings.Split(externalRef, ":")

	if len(parts) != 6 {
		return CostExternalRef{}, errors.New("invalid external ref format")
	}

	if parts[0] != "COST" {
		return CostExternalRef{}, errors.New("invalid external ref format")
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return CostExternalRef{}, errors.New("invalid external ref format")
	}

	year, err := strconv.Atoi(parts[4])
	if err != nil {
		return CostExternalRef{}, errors.New("invalid external ref format")
	}

	return CostExternalRef{
		CompanyID: parts[1],
		Day:       day,
		Month:     parts[3],
		Year:      year,
		PaymentID: parts[5],
	}, nil
}
