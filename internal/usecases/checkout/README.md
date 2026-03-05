# Usecase / Checkout

Orquestra formas de pagamento (dinheiro, PIX, cartão, carteira) e integrações externas (MercadoPago, POS) durante o fechamento do pedido.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/checkout/calculate` | handler/checkout.go | Calcula totais, troco e taxas antes de confirmar. |
| POST | `/checkout/confirm` | handler/checkout.go | Confirma pagamento local/online e gera recibo. |
| POST | `/checkout/webhook/mercadopago` | handler/checkout.go | Recebe notificações de status dos pagamentos online. |

## 2. Dependências
- Repositories: order, company, client, payment.
- Services: mercadopago, pos, rabbitmq (notificações).
- Usecases conectados: order, company, client.

## 3. Fluxos e exemplos
### Pré-cálculo
Passos:
- Recebe itens e métodos selecionados.
- Aplica regras de empresa (taxa de entrega, desconto fidelidade).
- Retorna resumo com total, troco sugerido e parcelas.

Exemplo de request:
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
Resposta:
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

### Confirmação POS
Passos:
- Valida que o pedido está `pending_payment`.
- Dispara comando para POS e aguarda callback.
- Atualiza `order_payment` + envia recibo por email.

Exemplo de request:
```json
{
  "order_id": "ord-100",
  "method": {
    "type": "pos",
    "terminal_id": "pos-3",
    "installments": 2
  }
}
```
Resposta:
```json
{
  "order_id": "ord-100",
  "status": "paid",
  "payments": [
    {
      "type": "pos",
      "authorization_code": "A12345"
    }
  ]
}
```

### Webhook MercadoPago
Passos:
- Valida assinatura do webhook.
- Consulta API MP para obter status definitivo.
- Sincroniza com pedido e emite evento `checkout.status_updated`.

Exemplo de request:
```json
{
  "id": "mp-evt-90",
  "type": "payment",
  "data": {
    "id": "pay-909"
  }
}
```
Resposta:
```json
{
  "received": true
}
```

## 4. Falhas conhecidas
- ErrInvalidSplit: soma das parcelas diferente do total.
- ErrPaymentAlreadyConfirmed: tentativa duplicada.
- MercadoPagoError: falha na API externa (propagar HTTP 502).

## 5. Notas operacionais
- Sempre salvar payload bruto do webhook para auditoria.
- Empresas podem desabilitar métodos via `company.preferences.checkout_methods`.
