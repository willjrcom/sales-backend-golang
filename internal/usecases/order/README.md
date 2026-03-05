# Usecase / Order

Usecase central do fluxo de pedidos: cria grupos/itens, coordena status, integra com estoque, pagamentos, fila e impressão.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/order` | handler/order.go | Cria pedido em modo draft. |
| POST | `/order/{id}/items` | handler/item.go | Adiciona itens e adicionais. |
| POST | `/order/{id}/status` | handler/order.go | Transiciona status. |
| DELETE | `/order/{id}` | handler/order.go | Cancela pedido e restaura recursos. |

## 2. Dependências
- Repositories: order, group_item, item, payment, client.
- Services: stock, rabbitmq (fila), print_manager, checkout.

## 3. Fluxos e exemplos
### Criar pedido
Passos:
- Cria registro base em `order` com status `draft`.
- Cria grupo inicial e itens solicitados.
- Dispara reserva de estoque e envia dados para fila/cozinha.

Exemplo de request:
```json
{
  "client_id": "cli-200",
  "type": "delivery",
  "items": [
    {
      "product_id": "prod-1",
      "quantity": 1
    }
  ]
}
```
Resposta:
```json
{
  "order_id": "ord-500",
  "status": "pending",
  "queue_number": 32
}
```

### Atualizar status
Passos:
- Valida transição (ex.: pending -> in_progress).
- Notifica fila e impressão se aplicável.
- Grava auditoria no histórico.

Exemplo de request:
```json
{
  "next_status": "in_progress",
  "reason": "Cozinha iniciou preparo"
}
```
Resposta:
```json
{
  "order_id": "ord-500",
  "status": "in_progress"
}
```

### Cancelar pedido
Passos:
- Valida se status permite cancelamento.
- Chama estoque para restaurar reservas ou debitos.
- Atualiza pagamentos (estorno) e comunica canais externos.

Exemplo de request:
```json
{
  "reason": "Cliente desistiu"
}
```
Resposta:
```json
{
  "order_id": "ord-500",
  "status": "canceled"
}
```

## 4. Falhas conhecidas
- ErrInvalidStatusTransition
- ErrOrderAlreadyFinished

## 5. Notas operacionais
- Pedidos multi-canal devem carregar `source_channel` para análise.
