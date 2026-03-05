# Repository / Postgres / Group Item

Snapshots de grupos de itens e status de produção.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `GetWithItems(ctx, groupID)` | JOIN com order_items para retornar composição. |
| `UpdateStatus(ctx, groupID)` | Atualiza status/employee com histórico. |

## 2. Transações e locking
- Criação de itens/grupos acontece na mesma tx do pedido para preservar sequência.

## 3. Exemplo de SQL
```sql
SELECT gi.id, gi.status, json_agg(oi.*) items
FROM group_items gi
JOIN order_items oi ON oi.group_item_id=gi.id
WHERE gi.id=@group
GROUP BY gi.id;
```

## 4. Notas operacionais
- Indexar por order_id para dashboards.
- Status segue enum `status_group_item`.
