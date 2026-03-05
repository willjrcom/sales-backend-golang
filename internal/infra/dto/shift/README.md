# DTO / Shift

DTOs para abertura/fechamento e detalhes do turno.

---

## 1. Onde é usado
- handler/shift.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| ShiftOpenRequest | employee_id, cash_float | request |
| ShiftCloseRequest | shift_id, cash_counted | request |
| ShiftResponse | id, status, totals, cash_diff | response |

## 3. Regras de validação
- `cash_float` >=0.

## 4. Exemplo de request
```json
{
  "employee_id": "emp-1",
  "cash_float": 200
}
```

## 5. Exemplo de response
```json
{
  "id": "shf-1",
  "status": "open",
  "totals": {
    "cash": 0
  }
}
```

## 6. Notas e compatibilidade
- Totais retornam valores em decimal string.
