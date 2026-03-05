# Repository / Postgres / Employee

Funcionários, pagamentos e vínculos com usuários.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `CreateWithPerson(ctx, Employee)` | Cria employee + person + contacts. |
| `ListPayments(ctx, employeeID)` | Consulta tabela employee_payments por período. |
| `Deactivate(ctx, id)` | Marca como inativo, preservando histórico. |

## 2. Transações e locking
- Pagamentos e fechamento de shift precisam da mesma tx para manter saldos.

## 3. Exemplo de SQL
```sql
SELECT e.id, e.role, p.name
FROM employees e
JOIN persons p ON p.id=e.person_id
WHERE e.company_id=@company;
```

## 4. Notas operacionais
- Não remover employees; usar soft-delete.
- Relacionamentos com pedidos exigem consistência histórica.
