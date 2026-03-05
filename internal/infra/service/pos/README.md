# Service / POS

Integração com terminais POS locais (USB/Ethernet) para autorizar transações presenciais.

---

## 1. Responsabilidades
- Descobrir terminais disponíveis.
- Enviar comandos (SALE, CANCEL) e receber respostas.
- Tratar tempo limite e reimpressão.

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `SendSale(ctx, terminalID string, amount decimal.Decimal) (POSResponse, error)` | Dispara venda. |
| `Cancel(ctx, terminalID, authorizationCode string) error` | Solicita cancelamento. |
| `Status(ctx, terminalID string) (POSStatus, error)` | Consulta saúde do terminal. |

## 3. Fluxo típico
- Checkout chama `SendSale` com terminal escolhido.
- Serviço comunica via socket/SDK, aguarda resposta.
- Atualiza order_payment e imprime recibo.

## 4. Configuração / Env Vars
- `POS_BRIDGE_URL`
- `POS_TIMEOUT_MS`

## 5. Exemplo de uso
```go
go
resp, err := pos.SendSale(ctx, "pos-3", decimal.NewFromFloat(45.90))
if err != nil { return err }
orderPayments.AddPOS(resp)
```

## 6. Falhas comuns
- ErrPOSTimeout
- ErrPOSOffline

## 7. Notas operacionais
- Sempre apresente opção de fallback (PIX/dinheiro).
