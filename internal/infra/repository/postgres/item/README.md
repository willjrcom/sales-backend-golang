# Repository / Postgres / Item

Gerencia itens do pedido, adicionais e snapshots de preço.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `Create(ctx, item)` | INSERT ... RETURNING para obter IDs. |
| `UpdateQuantity(ctx, id)` | Atualiza quantidade/customizações. |
| `ListByOrder(ctx, orderID)` | JOIN com produtos/adicionais para resposta completa. |

## 2. Transações e locking
- Sempre chamada junto com order/group_item/stock reservation na mesma tx.

## 3. Exemplo de SQL
```sql
SELECT id, product_id, quantity, additions
FROM order_items
WHERE order_id=@order;
```

## 4. Notas operacionais
- `additions` guardado como jsonb para auditoria.
- Garantir locking quando alterar após enviar à cozinha.
