# DTO / Item

DTOs para itens de pedido e adicionais.

---

## 1. Onde é usado
- handler/item.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| OrderItemRequest | product_id, variation_id, quantity, additions[] | request |
| OrderItemResponse | id, product_id, quantity, unit_price, additions | response |

## 3. Regras de validação
- `quantity` > 0.
- Adicionar `employee_id` quando item for inserido pela cozinha.

## 4. Exemplo de request
```json
{
  "product_id": "prod-1",
  "variation_id": "var-1",
  "quantity": 2,
  "additions": [
    "bacon"
  ]
}
```

## 5. Exemplo de response
```json
{
  "id": "item-1",
  "product_id": "prod-1",
  "quantity": 2,
  "unit_price": 25,
  "additions": [
    "bacon"
  ]
}
```

## 6. Notas e compatibilidade
- `additions` deve guardar snapshot de preço.
