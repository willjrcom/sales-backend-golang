# DTO / Checkout

DTOs do cálculo e confirmação do checkout (dinheiro, pix, cartão, POS).

---

## 1. Onde é usado
- handler/checkout.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| CheckoutCalculateRequest | order_id, methods[], tip_amount | request |
| CheckoutCalculateResponse | order_total, change, fees, status | response |
| CheckoutConfirmRequest | order_id, method, metadata | request |
| CheckoutConfirmResponse | order_id, status, payments[] | response |

## 3. Regras de validação
- Soma dos métodos deve igualar total dentro de margem de 0.01.
- `method.type` aceita `cash`, `pix`, `card`, `pos`.
- Quando `method.type=pos`, exigir `terminal_id`.

## 4. Exemplo de request
```json
{
  "order_id": "ord-100",
  "methods": [
    {
      "type": "cash",
      "amount": 50,
      "received": 100
    },
    {
      "type": "pix",
      "amount": 20
    }
  ]
}
```

## 5. Exemplo de response
```json
{
  "order_total": 70,
  "change": 30,
  "fees": {
    "delivery": 5
  },
  "status": "pending_confirmation"
}
```

## 6. Notas e compatibilidade
- Valores monetários em decimal com ponto.
- Respeitar locale no frontend para exibir valores.
