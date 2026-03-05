# DTO / Employee

DTOs para funcionários e pagamentos.

---

## 1. Onde é usado
- handler/employee.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| EmployeeRequest | name, document, role, create_user, contact | request |
| EmployeeResponse | id, role, status, user_id | response |
| EmployeePaymentRequest | amount, shift_id, type | request |

## 3. Regras de validação
- Documento CPF obrigatório.
- `role` deve ser permitido no RBAC.

## 4. Exemplo de request
```json
{
  "name": "João",
  "document": "12345678900",
  "role": "cashier",
  "create_user": true
}
```

## 5. Exemplo de response
```json
{
  "id": "emp-10",
  "role": "cashier",
  "status": "active",
  "user_id": "usr-55"
}
```

## 6. Notas e compatibilidade
- Pagamentos monetários em decimal string.
