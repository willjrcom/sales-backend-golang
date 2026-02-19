CREATE TABLE IF NOT EXISTS public.company_to_category (
    company_id UUID NOT NULL,
    category_id UUID NOT NULL,
    PRIMARY KEY (company_id, category_id),
    FOREIGN KEY (company_id) REFERENCES public.companies (id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES public.company_categories (id) ON DELETE CASCADE
);

-- Migrate existing data
INSERT INTO public.company_to_category (company_id, category_id)
SELECT id, category_id FROM public.companies WHERE category_id IS NOT NULL;

-- Drop foreign key constraint first
ALTER TABLE public.companies DROP CONSTRAINT IF EXISTS companies_category_id_fkey;

-- Drop index if exists
DROP INDEX IF EXISTS idx_companies_category_id;

-- Drop column
ALTER TABLE public.companies DROP COLUMN IF EXISTS category_id;
