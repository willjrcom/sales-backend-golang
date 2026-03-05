# DTO / Order Delivery

DTOs para entrega (endereços, taxas, driver).

---

## 1. Onde é usado
- handler/order_delivery.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| OrderDeliveryRequest | order_id, address_id, driver_id, eta | request |
| OrderDeliveryResponse | order_id, address, driver, status, tracking_code | response |

## 3. Regras de validação
- `address_id` obrigatório.
- `eta` em minutos.

## 4. Exemplo de request
```json
{
  "order_id": "ord-500",
  "address_id": "addr-1",
  "driver_id": "drv-1",
  "eta": 35
}
```

## 5. Exemplo de response
```json
{
  "order_id": "ord-500",
  "status": "on_route",
  "driver": {
    "id": "drv-1",
    "name": "Carlos"
  },
  "tracking_code": "TRK123"
}
```

## 6. Notas e compatibilidade
- Quando driver muda, atualizar `changed_by`.
