# Repository / Postgres / Schema

Metadados do schema multi-tenant e histórico de migrações.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `ListTenants(ctx)` | Consulta tabela pública com schemas ativos. |
| `RecordMigration(ctx)` | Insere registro em schema_migrations. |
| `GetSchemaByCompany(ctx)` | Retorna nome do schema atrelado à empresa. |

## 2. Transações e locking
- Executado no schema public; não altera search_path dos tenants.

## 3. Exemplo de SQL
```sql
SELECT schema_name FROM schemas WHERE company_id=@company;
```

## 4. Notas operacionais
- Usado pelos comandos migrate/migrate-all.
- Nunca remover schemas sem backup.
