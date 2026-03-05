# Domain / Shift

Turnos de trabalho, métricas e pagamentos agregados.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Shift | Abertura, fechamento, operadores. |
| OrderProcessAnalytics | KPIs de tempo. |
| Redeem | Resgates/gorjetas. |

## 2. Regras de negócio
- Um funcionário só pode ter um turno aberto.
- Fechamento calcula divergência de caixa.

## 3. Interações e consumidores
- Usecases: shift, employee, print_manager.

## 4. Exemplo de estrutura
```json
{
  "id": "shf-1",
  "employee_id": "emp-1",
  "status": "open",
  "opened_at": "2026-03-05T08:00:00Z"
}
```
