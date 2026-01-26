package companyusecases

import (
	"context"
	"log"
	"time"
)

func (s *Service) StartSubscriptionWatcher(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = 24 * time.Hour
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		run := func() {
			if err := s.ProcessSubscriptionExpirations(context.Background()); err != nil {
				log.Printf("subscription watcher error: %v", err)
			}
		}

		run()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				run()
			}
		}
	}()
}

func (s *Service) ProcessSubscriptionExpirations(ctx context.Context) error {
	companies, err := s.r.ListCompaniesForBilling(ctx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, company := range companies {
		if company.SubscriptionExpiresAt == nil {
			continue
		}

		expired := company.SubscriptionExpiresAt.Before(now)
		switch {
		case expired && !company.IsBlocked:
			if err := s.r.UpdateCompanySubscription(ctx, company.ID, company.SchemaName, company.SubscriptionExpiresAt, true); err != nil {
				return err
			}
		case !expired && company.IsBlocked:
			if err := s.r.UpdateCompanySubscription(ctx, company.ID, company.SchemaName, company.SubscriptionExpiresAt, false); err != nil {
				return err
			}
		}
	}

	return nil
}
