DO $$
DECLARE
    rec record;
    new_schema text;
BEGIN
    FOR rec IN
        SELECT schema_name
        FROM information_schema.schemata
        WHERE schema_name LIKE 'loja_%'
        ORDER BY schema_name
    LOOP
        new_schema := 'company_' || substring(rec.schema_name FROM 6);

        IF EXISTS (
            SELECT 1
            FROM information_schema.schemata
            WHERE schema_name = new_schema
        ) THEN
            RAISE NOTICE 'Skipping schema % because % already exists', rec.schema_name, new_schema;
            CONTINUE;
        END IF;

        EXECUTE format('ALTER SCHEMA %I RENAME TO %I', rec.schema_name, new_schema);

        UPDATE public.companies
        SET schema_name = new_schema
        WHERE schema_name = rec.schema_name;
    END LOOP;
END;
$$;
