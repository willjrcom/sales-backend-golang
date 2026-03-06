package productentity

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestParseSplitPricingStrategy(t *testing.T) {
	strategy, err := ParseSplitPricingStrategy("average")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if strategy != SplitPricingStrategyAverage {
		t.Fatalf("expected %s, got %s", SplitPricingStrategyAverage, strategy)
	}

	if _, err := ParseSplitPricingStrategy("invalid"); err == nil {
		t.Fatal("expected error for invalid strategy, got nil")
	}
}

func TestCalculateSplitPrice(t *testing.T) {
	prices := []decimal.Decimal{
		decimal.NewFromInt(40),
		decimal.NewFromInt(50),
	}

	if got := CalculateSplitPrice(SplitPricingStrategyHighestItem, prices); !got.Equal(decimal.NewFromInt(50)) {
		t.Fatalf("expected 50, got %s", got.String())
	}

	if got := CalculateSplitPrice(SplitPricingStrategyAverage, prices); !got.Equal(decimal.NewFromInt(45)) {
		t.Fatalf("expected 45, got %s", got.String())
	}

	if got := CalculateSplitPrice(SplitPricingStrategySum, prices); !got.Equal(decimal.NewFromInt(90)) {
		t.Fatalf("expected 90, got %s", got.String())
	}

	if got := CalculateSplitPrice(SplitPricingStrategy("unsupported"), prices); !got.Equal(decimal.NewFromInt(50)) {
		t.Fatalf("fallback expected 50, got %s", got.String())
	}
}
