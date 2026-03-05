# Repository / Postgres / Table

Mesas físicas/virtuais vinculadas a places.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `ListByPlace(ctx)` | Retorna mesas e status. |
| `Lock(ctx, tableID)` | SELECT FOR UPDATE durante abertura. |
| `UpdateStatus(ctx)` | Atualiza status/ocupação e atendente. |

## 2. Transações e locking
- Abrir/fechar mesa deve ocorrer junto ao order_table repository.
- Transações evitam dupla atribuição.

## 3. Exemplo de SQL
```sql
SELECT id, name, status
FROM tables
WHERE place_id=@place
ORDER BY name;
```

## 4. Notas operacionais
- `status` deriva de pedidos abertos; sincronize com order_table.
- Mesas virtuais devem ter flag `virtual=true`.
