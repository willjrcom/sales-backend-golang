# Repository / Postgres / Order

Consultas completas de pedidos, itens, pagamentos e entregas.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `Create(ctx, order)` | INSERT e retorna order_id. |
| `GetByID(ctx, id)` | Carrega pedido completo (JOIN com itens/pagamentos). |
| `List(ctx, filters)` | Filtro por status, canal, período com paginação. |
| `UpdateStatus(ctx, id)` | Atualiza status com optimistic locking (version). |

## 2. Transações e locking
- Fechamento do pedido envolve order, payments, stock e checkout em tx única.
- Cancelamentos precisam restaurar estoque antes de apagar itens.

## 3. Exemplo de SQL
```sql
SELECT o.id, o.status, o.total, json_agg(oi.* ORDER BY oi.created_at) items
FROM orders o
LEFT JOIN order_items oi ON oi.order_id=o.id
WHERE o.id=@order
GROUP BY o.id;
```

## 4. Notas operacionais
- Manter índices compostos em (status, created_at).
- Queries complexas usam CTEs para clareza; documente mudanças.
