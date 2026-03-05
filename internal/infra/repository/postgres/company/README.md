# Repository / Postgres / Company

Persiste dados do tenant, preferências, assinaturas e billing.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `CreateWithSchema(ctx, Company)` | Insere empresa e registra schema provisionado. |
| `UpdatePreferences(ctx, id, prefs)` | Atualiza jsonb com versionamento e retorna diff. |
| `ListUsageCost(ctx, companyID, period)` | Retorna custos variáveis agregados por período. |

## 2. Transações e locking
- Onboarding executa `Create` + `schemaRepository.Create` na mesma tx pública.
- Preferências críticas usam `FOR UPDATE` para evitar corrida.

## 3. Exemplo de SQL
```sql
UPDATE companies
SET preferences = preferences || @prefs::jsonb, updated_at = NOW()
WHERE id = @company
RETURNING *;
```

## 4. Notas operacionais
- Campos sensíveis (tokens) não devem ser retornados diretamente.
- Sempre logar `company_id` e usuário responsável.
