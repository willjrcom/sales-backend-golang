-- Add image_path column to companies table
ALTER TABLE public.companies ADD COLUMN IF NOT EXISTS image_path VARCHAR(255);
