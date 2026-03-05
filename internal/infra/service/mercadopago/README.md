# Service / MercadoPago

Cliente REST para criar pagamentos, verificar status e lidar com webhooks do MercadoPago.

---

## 1. Responsabilidades
- Criar preferências/checkout para pedidos delivery/pickup.
- Consultar status de pagamentos.
- Validar webhooks assinados.

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `CreatePayment(ctx, req PaymentRequest) (PaymentResponse, error)` | Cria pagamento e retorna `id`. |
| `GetPayment(ctx, id string) (PaymentResponse, error)` | Consulta status final. |
| `ValidateWebhook(h http.Header, body []byte) error` | Verifica assinatura HMAC. |

## 3. Fluxo típico
- Checkout confirma método online → `CreatePayment`.
- Ao receber webhook → `ValidateWebhook` → `GetPayment` → atualiza order_payment.
- Erros são propagados para reprocessamento manual.

## 4. Configuração / Env Vars
- `MERCADOPAGO_ACCESS_TOKEN`
- `MERCADOPAGO_WEBHOOK_SECRET`

## 5. Exemplo de uso
```go
go
resp, err := mercadopago.CreatePayment(ctx, PaymentRequest{Amount: order.Total, Description: order.ID})
```

## 6. Falhas comuns
- ErrPaymentDeclined
- ErrWebhookSignature

## 7. Notas operacionais
- Sempre armazene `payment_id` para conciliações futuras.
