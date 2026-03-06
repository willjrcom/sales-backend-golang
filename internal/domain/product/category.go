package productentity

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrInvalidSplitPricingStrategy = errors.New("invalid split pricing strategy")
)

type SplitPricingStrategy string

const (
	SplitPricingStrategyHighestItem SplitPricingStrategy = "highest_item"
	SplitPricingStrategyAverage     SplitPricingStrategy = "average"
	SplitPricingStrategySum         SplitPricingStrategy = "sum"
)

func (s SplitPricingStrategy) IsValid() bool {
	switch s {
	case SplitPricingStrategyHighestItem,
		SplitPricingStrategyAverage,
		SplitPricingStrategySum:
		return true
	default:
		return false
	}
}

func (s SplitPricingStrategy) OrDefault() SplitPricingStrategy {
	if s.IsValid() {
		return s
	}
	return SplitPricingStrategyHighestItem
}

func ParseSplitPricingStrategy(value string) (SplitPricingStrategy, error) {
	if value == "" {
		return SplitPricingStrategyHighestItem, nil
	}

	strategy := SplitPricingStrategy(value)
	if !strategy.IsValid() {
		return "", ErrInvalidSplitPricingStrategy
	}
	return strategy, nil
}

type ProductCategory struct {
	entity.Entity
	ProductCategoryCommonAttributes
}

type ProductCategoryCommonAttributes struct {
	Name                 string
	ImagePath            string
	NeedPrint            bool
	PrinterName          string
	UseProcessRule       bool
	RemovableIngredients []string
	IsActive             bool
	Sizes                []Size
	Products             []Product
	ProcessRules         []ProcessRule
	IsAdditional         bool
	IsComplement         bool
	AdditionalCategories []ProductCategory
	ComplementCategories []ProductCategory
	AllowFractional      bool
	SplitPricingStrategy SplitPricingStrategy
}

func NewProductCategory(categoryCommonAttributes ProductCategoryCommonAttributes) *ProductCategory {
	categoryCommonAttributes.SplitPricingStrategy = categoryCommonAttributes.SplitPricingStrategy.OrDefault()

	return &ProductCategory{
		Entity:                          entity.NewEntity(),
		ProductCategoryCommonAttributes: categoryCommonAttributes,
	}
}

// CalculateSplitPrice applies the configured pricing strategy to the provided prices.
func CalculateSplitPrice(strategy SplitPricingStrategy, prices []decimal.Decimal) decimal.Decimal {
	if len(prices) == 0 {
		return decimal.Zero
	}

	switch strategy {
	case SplitPricingStrategySum:
		total := decimal.Zero
		for _, price := range prices {
			total = total.Add(price)
		}
		return total
	case SplitPricingStrategyAverage:
		total := decimal.Zero
		for _, price := range prices {
			total = total.Add(price)
		}
		return total.Div(decimal.NewFromInt(int64(len(prices)))).Round(2)
	case SplitPricingStrategyHighestItem:
		fallthrough
	default:
		highest := prices[0]
		for _, price := range prices[1:] {
			if price.GreaterThan(highest) {
				highest = price
			}
		}
		return highest
	}
}
