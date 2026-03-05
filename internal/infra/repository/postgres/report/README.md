# Repository / Postgres / Report

Consultas analíticas customizadas (dashboards de vendas, estoque, marketing).

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `SalesSummary(ctx, filter)` | Agrega vendas por período/canal. |
| `AdditionalItemsSold(ctx)` | Top adicionais vendidos. |
| `ComplementsSold(ctx)` | Top complementos (limit 10). |

## 2. Transações e locking
- Somente leitura, mas usar search_path correto e `SET ROLE` se necessário.

## 3. Exemplo de SQL
```sql
SELECT date_trunc(@group_by, created_at) bucket, SUM(total) gross
FROM orders
WHERE created_at BETWEEN @from AND @to
GROUP BY bucket ORDER BY bucket;
```

## 4. Notas operacionais
- Adicionar LIMIT/OFFSET para evitar respostas gigantes.
- Considere materializar consultas pesadas.
