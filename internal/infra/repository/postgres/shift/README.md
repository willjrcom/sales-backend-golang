# Repository / Postgres / Shift

Turnos, taxas de entregador e analytics de produção.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `Open(ctx, shift)` | INSERT com retorno de shift_id. |
| `Close(ctx, shiftID)` | Atualiza totais, status e registra divergências. |
| `List(ctx, filters)` | JOIN orders/pagamentos para relatório. |

## 2. Transações e locking
- Fechamento envolve atualizar orders/pagamentos → rodar dentro de tx.

## 3. Exemplo de SQL
```sql
UPDATE shifts
SET status='closed', cash_diff=@diff, closed_at=NOW()
WHERE id=@id RETURNING *;
```

## 4. Notas operacionais
- Armazena `order_process_analytics` para KPIs.
- Bloqueia novos pedidos para atendente com shift fechado.
