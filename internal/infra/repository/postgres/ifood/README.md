# Repository / Postgres / Ifood

Integração com pedidos importados do iFood.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `UpsertToken(ctx, companyID)` | Armazena tokens refresh/access. |
| `SaveOrderPayload(ctx, payload)` | Persist payload bruto para auditoria. |
| `ListPendingSync(ctx)` | Retorna pedidos aguardando confirmação. |

## 2. Transações e locking
- Executa no schema público; não usar search_path de tenants.

## 3. Exemplo de SQL
```sql
INSERT INTO ifood_orders(order_id, payload)
VALUES (@id, @payload)
ON CONFLICT(order_id) DO UPDATE SET payload=EXCLUDED.payload;
```

## 4. Notas operacionais
- Segregar dados iFood do core.
- Nunca misturar IDs públicos com internos sem mapear.
