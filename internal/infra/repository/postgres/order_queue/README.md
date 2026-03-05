# Repository / Postgres / Order Queue

Persistência de filas com posição e métricas de espera.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `Enqueue(ctx, orderID, processID)` | Calcula posição e insere entrada. |
| `Advance(ctx, queueID)` | Incrementa etapa e registra timestamps. |
| `ListByProcess(ctx, processID)` | Retorna fila ordenada por posição. |

## 2. Transações e locking
- `Advance` usa `FOR UPDATE` para evitar duas estações mexendo no mesmo item.

## 3. Exemplo de SQL
```sql
SELECT id, order_id, position, status
FROM order_queue
WHERE process_id=@proc
ORDER BY position;
```

## 4. Notas operacionais
- Recalcular posições periodicamente para evitar gaps.
- Expor tempo médio via views auxiliares.
