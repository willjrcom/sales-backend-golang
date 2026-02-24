-- Add started_at, ended_at and type to advertisements table
ALTER TABLE public.advertisements ADD COLUMN IF NOT EXISTS started_at TIMESTAMP;
ALTER TABLE public.advertisements ADD COLUMN IF NOT EXISTS ended_at TIMESTAMP;
ALTER TABLE public.advertisements ADD COLUMN IF NOT EXISTS type TEXT NOT NULL DEFAULT 'standard';
