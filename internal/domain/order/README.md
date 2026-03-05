# Domain / Order

Agregado principal: pedido, grupos, itens, pagamentos, entregas.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Order | Metadados e status. |
| OrderItem | Itens com adicionais e preços travados. |
| OrderGroupItem | Agrupa itens por workflow. |
| OrderPayment | Pagamentos múltiplos. |
| OrderDelivery/Pickup/Table | Modalidades específicas. |
| Coupon | Descontos aplicados. |
| Status enums | StatusOrder, StatusItem etc. |

## 2. Regras de negócio
- Status segue máquina (draft → pending → in_progress → finished/canceled).
- Itens armazenam snapshot de preço/adicionais para auditoria.
- Pedidos delivery vinculam driver/endereço; mesa vincula `order_table`.

## 3. Interações e consumidores
- Usecases: order, checkout, order_table, order_delivery, stock.
- Infra: handler/order.go, repository/postgres/order.

## 4. Exemplo de estrutura
```json
{
  "id": "ord-500",
  "status": "pending",
  "type": "delivery",
  "items": [
    {
      "item_id": "it-1",
      "product_id": "prod-1",
      "quantity": 1,
      "unit_price": 25
    }
  ],
  "payments": [
    {
      "type": "cash",
      "amount": 25
    }
  ]
}
```
