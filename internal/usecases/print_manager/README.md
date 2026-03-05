# Usecase / Print Manager

Centraliza impressões de pedidos, cozinha, fiscais e relatórios de turno.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/print/order` | handler/order_print.go | Imprime pedido completo. |
| POST | `/print/kitchen` | handler/order_print.go | Imprime tickets por estação. |
| POST | `/print/shift` | handler/shift_print.go | Imprime fechamento de turno. |

## 2. Dependências
- Services: printer (ESC/POS), rabbitmq.
- Usecases: order, shift, group_item.

## 3. Fluxos e exemplos
### Ticket de cozinha
Passos:
- Recebe pedido/grupo atualizado.
- Agrupa itens por estação definida em process_rule.
- Renderiza template ESC/POS e envia para impressora cadastrada.

Exemplo de request:
```json
{
  "order_id": "ord-500",
  "station": "grill"
}
```
Resposta:
```json
{
  "printed": true,
  "printer_id": "prn-2"
}
```

## 4. Falhas conhecidas
- ErrPrinterOffline
- ErrTemplateNotFound

## 5. Notas operacionais
- Retentar 3x antes de marcar como falha e alertar operador.
