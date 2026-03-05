# Repository / Postgres / Delivery Driver

Persistência de entregadores, taxas, status e zonas de atuação.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `List(ctx, filters)` | JOIN em employees e shifts para dashboards. |
| `UpdateStatus(ctx, id, status)` | Atualiza status e order atual com `FOR UPDATE`. |
| `Create(ctx, Driver)` | Insere driver e configura taxas default. |

## 2. Transações e locking
- Mudanças de status ocorrem em tx para evitar race ao redistribuir pedidos.

## 3. Exemplo de SQL
```sql
UPDATE delivery_drivers
SET status=@status, current_order_id=@order, updated_at=NOW()
WHERE id=@id
RETURNING *;
```

## 4. Notas operacionais
- Adicionar índices em status/company_id.
- Sempre limpar referência a `current_order_id` ao finalizar entrega.
