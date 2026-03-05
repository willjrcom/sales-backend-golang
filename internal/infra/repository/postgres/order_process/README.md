# Repository / Postgres / Order Process

Processos/filas de produção (etapas, responsáveis).

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `ListActive(ctx)` | Retorna processos em andamento + fila. |
| `UpdateStep(ctx)` | Atualiza status/employee e timestamps. |
| `LinkProduct(ctx)` | Relaciona produtos/categorias às etapas. |

## 2. Transações e locking
- Atualizações em processos e filas devem ocorrer na mesma tx para manter posição consistente.

## 3. Exemplo de SQL
```sql
UPDATE order_process
SET status=@status, employee_id=@emp,
    started_at=COALESCE(started_at, NOW())
WHERE id=@id;
```

## 4. Notas operacionais
- Índices parciais em status IN (pending,in_progress).
