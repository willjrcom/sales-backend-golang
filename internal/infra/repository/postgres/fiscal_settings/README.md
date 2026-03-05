# Repository / Postgres / Fiscal Settings

Configurações fiscais (CRT, CNAE, certificados).

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `Upsert(ctx, settings)` | Atualiza jsonb criptografado. |
| `GetByCompany(ctx)` | Recupera settings decodificados para emissão. |

## 2. Transações e locking
- Use `FOR UPDATE` ao atualizar para evitar race com emissão simultânea.

## 3. Exemplo de SQL
```sql
INSERT INTO fiscal_settings (company_id, data)
VALUES (@company, @json::jsonb)
ON CONFLICT (company_id)
DO UPDATE SET data = EXCLUDED.data, updated_at=NOW();
```

## 4. Notas operacionais
- Nunca logar certificados ou senhas.
- Executar teste de conexão após salvar.
