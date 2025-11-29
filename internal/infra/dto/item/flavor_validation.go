package itemdto

import (
	"errors"
	"strings"
)

var (
	ErrFlavorRequired = errors.New("flavor is required for this product")
	ErrFlavorInvalid  = errors.New("flavor is invalid for this product")
)

// normalizeFlavorSelection validates and normalizes the selected flavor against the product options.
// It returns nil when the product does not expose flavor options.
func normalizeFlavorSelection(selection *string, options []string) (*string, error) {
	if len(options) == 0 {
		return nil, nil
	}

	if selection == nil || strings.TrimSpace(*selection) == "" {
		return nil, ErrFlavorRequired
	}

	value := strings.TrimSpace(*selection)
	for _, option := range options {
		if option == value {
			normalized := value
			return &normalized, nil
		}
	}

	return nil, ErrFlavorInvalid
}

// NormalizeFlavor enforces flavor selection for API payloads.
func NormalizeFlavor(selection *string, options []string) (*string, error) {
	return normalizeFlavorSelection(selection, options)
}
