package model

import "testing"

func TestProductToDomain_AllowsNilCategoryAndNilVariation(t *testing.T) {
	p := &Product{}
	p.Variations = []*ProductVariation{nil}

	d := p.ToDomain()
	if d == nil {
		t.Fatalf("expected non-nil domain product")
	}
	if d.Category != nil {
		t.Fatalf("expected nil category when model category is nil")
	}
	if len(d.Variations) != 0 {
		t.Fatalf("expected nil variations to be ignored")
	}
}
