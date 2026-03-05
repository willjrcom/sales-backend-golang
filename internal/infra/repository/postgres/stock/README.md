# Repository / Postgres / Stock

Controle de estoque, lotes e movimentos FIFO.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `GetByProduct(ctx, productID)` | Carrega estoque e batches relacionadas. |
| `LockForUpdate(ctx, stockID)` | SELECT ... FOR UPDATE durante reservas. |
| `InsertMovement(ctx, movement)` | Registra movimento e atualiza saldos. |

## 2. Transações e locking
- Reserva/Débito executam com itens/pedidos na mesma tx.
- Movimentos manuais usam locking pessimista.

## 3. Exemplo de SQL
```sql
SELECT s.id, s.current_stock, b.id AS batch_id, b.current_quantity
FROM stock s
LEFT JOIN stock_batches b ON b.stock_id=s.id
WHERE s.product_id=@product
ORDER BY b.created_at;
```

## 4. Notas operacionais
- Sempre usar service para manter política FIFO.
- Alertas são derivados da diferença Min/Max.
