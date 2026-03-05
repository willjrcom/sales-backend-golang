# DTO / Order Table

DTOs para mesas (abrir, mover, fechar).

---

## 1. Onde é usado
- handler/order_table.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| OpenTableRequest | table_id, employee_id, guests | request |
| OrderTableResponse | table_id, order_id, status, guests | response |

## 3. Regras de validação
- `guests` >=1.

## 4. Exemplo de request
```json
{
  "table_id": "tbl-10",
  "employee_id": "emp-7",
  "guests": 4
}
```

## 5. Exemplo de response
```json
{
  "table_id": "tbl-10",
  "order_id": "ord-600",
  "status": "open",
  "guests": 4
}
```

## 6. Notas e compatibilidade
- `status` sincronizado com domain/table.
