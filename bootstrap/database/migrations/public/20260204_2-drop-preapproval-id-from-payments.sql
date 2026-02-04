-- Drop PreapprovalID from Company Payments (Redundant with Company Subscriptions)
ALTER TABLE company_payments
DROP COLUMN IF EXISTS preapproval_id;
