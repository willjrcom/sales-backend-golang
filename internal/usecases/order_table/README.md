# Usecase / Order Table

Gerencia pedidos vinculados a mesas físicas (salão): abertura, transferência e fechamento.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/tables/{id}/open` | handler/order_table.go | Abre mesa e associa atendente. |
| POST | `/tables/{id}/move` | handler/order_table.go | Transfere pedido entre mesas. |
| POST | `/tables/{id}/close` | handler/order_table.go | Fecha mesa e consolida pagamentos. |

## 2. Dependências
- Repositories: order_table, table, order, employee.

## 3. Fluxos e exemplos
### Abrir mesa
Passos:
- Valida que a mesa está livre.
- Cria registro `order_table` com atendente responsável.
- Dispara criação de pedido do tipo `dine_in`.

Exemplo de request:
```json
{
  "employee_id": "emp-7",
  "guests": 4
}
```
Resposta:
```json
{
  "table_id": "tbl-10",
  "order_id": "ord-600",
  "status": "open"
}
```

### Transferir mesa
Passos:
- Valida destino livre.
- Atualiza mesa associada no pedido e notifica cozinha.
- Recalcula mapa do salão.

Exemplo de request:
```json
{
  "target_table_id": "tbl-11"
}
```
Resposta:
```json
{
  "from": "tbl-10",
  "to": "tbl-11",
  "order_id": "ord-600"
}
```

## 4. Falhas conhecidas
- ErrTableBusy
- ErrOrderTableMismatch

## 5. Notas operacionais
- Transferências devem ser logadas para auditoria de gorjeta.
