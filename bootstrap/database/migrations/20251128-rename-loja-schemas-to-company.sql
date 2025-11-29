DO $$
DECLARE
    rec record;
    suffix text;
    old_schema text;
    new_schema text;
    target_schema text;
BEGIN
    FOR rec IN
        SELECT schema_name
        FROM information_schema.schemata
        WHERE schema_name LIKE 'loja_%' OR schema_name LIKE 'company_%'
        ORDER BY schema_name
    LOOP
        IF rec.schema_name LIKE 'loja_%' THEN
            suffix := substring(rec.schema_name FROM 6);
        ELSE
            suffix := substring(rec.schema_name FROM 9);
        END IF;

        old_schema := 'loja_' || suffix;
        new_schema := 'company_' || suffix;
        target_schema := rec.schema_name;

        IF rec.schema_name LIKE 'loja_%' THEN
            IF EXISTS (
                SELECT 1
                FROM information_schema.schemata
                WHERE schema_name = new_schema
            ) THEN
                RAISE NOTICE 'Skipping schema % because % already exists', rec.schema_name, new_schema;
                CONTINUE;
            END IF;

            EXECUTE format('ALTER SCHEMA %I RENAME TO %I', rec.schema_name, new_schema);
            target_schema := new_schema;
        END IF;

        EXECUTE format(
            'UPDATE %I.companies SET schema_name = %L WHERE schema_name = %L OR schema_name = %L',
            target_schema,
            new_schema,
            old_schema,
            new_schema
        );

        UPDATE public.companies
        SET schema_name = new_schema
        WHERE schema_name = old_schema OR schema_name = new_schema;
    END LOOP;
END;
$$;
