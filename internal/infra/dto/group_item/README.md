# DTO / Group Item

DTOs que representam grupos de itens em pedidos (combo, etapa da cozinha).

---

## 1. Onde é usado
- handler/group_item.go
- handler/order.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| GroupItemRequest | name, employee_id, notes | request |
| GroupItemResponse | id, status, items[] | response |

## 3. Regras de validação
- `name` obrigatório.
- Itens pertencem a um grupo por vez.

## 4. Exemplo de request
```json
{
  "name": "Combo 1",
  "employee_id": "emp-10",
  "notes": "caprichar"
}
```

## 5. Exemplo de response
```json
{
  "id": "grp-1",
  "status": "pending",
  "items": [
    "item-1",
    "item-2"
  ]
}
```

## 6. Notas e compatibilidade
- `status` segue StatusGroup enum.
