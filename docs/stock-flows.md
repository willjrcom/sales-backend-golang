# Fluxos de Estoque — Backend Go

## Visão Geral

O sistema de estoque tem **5 fluxos principais** que cobrem o ciclo completo de um produto: do cadastro até a venda e geração de alertas.

```
          ADMINISTRAÇÃO                    PEDIDOS (automático)
          ─────────────                    ────────────────────
  Criar estoque                     AddItem  →  ReserveStock
  Entrada manual (add)              Finalizar → DebitStockFIFO
  Saída manual (remove)             Cancelar  → RestoreStock
  Ajuste (inventário)
          │
          ▼
     CheckAlerts → Alertas (low/out/over)
     CheckExpirations → Alertas (expired/near)
```

---

## Fluxo 1 — Entrada de Estoque (`POST /stock/{id}/movement/add`)

Usado quando chega mercadoria: NF, compra, devolução, etc.

### Entrada
```json
{
  "reason": "compra fornecedor ABC",
  "quantity": 50,
  "price": 8.50
}
```

### O que acontece internamente
```
1. Busca o Stock no banco (GetStockByID)
2. Chama stock.AddMovementStock(qty, reason, employeeID, price)
   → valida qty > 0
   → valida IsActive
   → CurrentStock += qty
3. Cria lote (StockBatch) com a quantidade e preço de custo
4. Salva o movimento (tipo: "in")
5. Salva o batch
6. Atualiza CurrentStock no banco
7. Roda CheckAlerts → se saiu de low_stock, não gera novo alerta
```

### Saída (movimento criado)
```json
{
  "id": "uuid",
  "stock_id": "uuid",
  "type": "in",
  "quantity": 50,
  "reason": "compra fornecedor ABC",
  "price": "8.50",
  "created_at": "2026-03-02T14:00:00Z"
}
```

**Efeito nos campos do estoque:**
| Campo | Antes | Depois |
|---|---|---|
| `current_stock` | 10 | 60 |
| `reserved_stock` | 0 | 0 |

---

## Fluxo 2 — Saída Manual de Estoque (`POST /stock/{id}/movement/remove`)

Usado para saídas que não são por pedido: quebra, desperdício, uso interno.

### Entrada
```json
{
  "reason": "produto vencido descartado",
  "quantity": 5
}
```

### O que acontece internamente
```
1. Busca o Stock no banco
2. Valida qty > 0 e IsActive
3. Valida CurrentStock >= qty (senão: ErrInsufficientStock)
4. FIFO: itera lotes mais antigos, debita até completar qty
   → Lote 1: 3 un → zera lote
   → Lote 2: 2 un → parcial
5. Salva movimentos por lote (tipo: "out")
6. Salva CurrentStock atualizado
7. CheckAlerts → pode gerar low_stock ou out_of_stock
```

**Efeito:**
| Campo | Antes | Depois |
|---|---|---|
| `current_stock` | 10 | 5 |
| `reserved_stock` | 0 | 0 |

---

## Fluxo 3 — Ajuste de Estoque (`POST /stock/{id}/movement/adjust`)

Usado para acerto de inventário: o sistema calcula a diferença e cria o movimento.

### Entrada
```json
{
  "new_stock": 20,
  "reason": "contagem física do inventário"
}
```

### O que acontece internamente
```
Diferença = new_stock - current_stock

Se diferença > 0:
  → tipo "adjust_in", cria lote com a diferença
Se diferença < 0:
  → tipo "adjust_out", FIFO nos lotes existentes
Se diferença == 0:
  → nenhum movimento, sem efeito
```

**Efeito:**
| Cenário | `current_stock` antes | `new_stock` | Tipo |
|---|---|---|---|
| Acréscimo | 5 | 20 | `adjust_in` (+15) |
| Redução | 20 | 8 | `adjust_out` (-12) |

---

## Fluxo 4 — Reserva por Pedido (automático via `AddItem`)

Quando um item é adicionado a um pedido, o sistema **reserva** o estoque. O produto sai do `current_stock` mas ainda não é consumido definitivamente.

### Gatilho
`POST /orders/{id}/items` → `AddItemOrder` → `DebitStockFromItem`

### O que acontece
```
1. Busca stock por ProductVariationID (ou ProductID como fallback)
2. Chama stock.ReserveStock(qty, orderID, employeeID, price)
   → valida qty > 0 e IsActive
   → valida CurrentStock >= qty
   → CurrentStock -= qty
   → ReservedStock += qty
3. Salva movimento tipo "reserve" (com OrderID)
4. Atualiza stock no banco
```

**Efeito:**
| Campo | Antes | Depois |
|---|---|---|
| `current_stock` | 10 | 7 |
| `reserved_stock` | 0 | 3 |

> O produto saiu da prateleira virtual (não pode ser vendido para outro pedido), mas ainda não foi "consumido" do lote físico.

---

## Fluxo 5 — Finalização do Pedido (FIFO real)

Quando o pedido é finalizado, o estoque reservado é debitado fisicamente dos lotes (FIFO).

### Gatilho
`PUT /orders/{id}/finish` → `FinishOrder` → `DebitStockFromOrder`

### O que acontece
```
Para cada item do pedido:
  1. Busca stock pelo produto/variação
  2. FIFO: itera lotes ordenados por data de entrada (mais antigo primeiro)
     → Lote 1: debita o que puder
     → Lote 2: debita o restante
     → ... até completar a quantidade
  3. Salva movimento tipo "out" com OrderID por lote
  4. ReservedStock -= qty (limpa a reserva)
  5. CheckAlerts → pode gerar low_stock ou out_of_stock
```

**Efeito:**
| Campo | Após reserva | Após finalização |
|---|---|---|
| `current_stock` | 7 | 7 (não muda — já saiu na reserva) |
| `reserved_stock` | 3 | 0 |

---

## Fluxo 6 — Cancelamento de Item/Pedido (restaura reserva)

Quando um item é cancelado ou o pedido é cancelado, o estoque reservado volta.

### Gatilho
`PUT /orders/{id}/items/{item_id}/cancel` → `CancelGroupItem` → `RestoreStockFromItem`

### O que acontece
```
1. Busca stock pelo produto/variação
2. Chama stock.RestoreStock(qty, orderID, employeeID, price, batchID)
   → valida ReservedStock >= qty
   → ReservedStock -= qty
   → CurrentStock += qty
3. Salva movimento tipo "restore" com OrderID
4. Atualiza stock no banco
```

**Efeito:**
| Campo | Durante pedido | Após cancelamento |
|---|---|---|
| `current_stock` | 7 | 10 (volta) |
| `reserved_stock` | 3 | 0 |

---

## Fluxo 7 — Alertas de Estoque

Gerados automaticamente após qualquer movimento. Também pode verificar vencimentos manualmente.

### Alertas de nível (`CheckAlerts`)

Disparado automaticamente após `add`, `remove`, `adjust`, `finalize`.

| Tipo | Condição |
|---|---|
| `low_stock` | `current_stock <= min_stock` e `current_stock > 0` |
| `out_of_stock` | `current_stock <= 0` |
| `over_stock` | `current_stock > max_stock` |

### Alertas de vencimento (`POST /stock/alerts/expiry/check?days=30`)

Varre todos os lotes com `current_quantity > 0` e verifica a data de vencimento.

| Tipo | Condição |
|---|---|
| `near_expiration` | Vence em até `days` dias |
| `expired` | Já venceu, mas ainda tem quantidade |

> **Detalhe importante:** lotes com `current_quantity = 0` são ignorados — produto já consumido não gera alerta.

---

## Mapa de Rotas

| Método | Rota | Fluxo |
|---|---|---|
| `POST` | `/stock/new` | Criar cadastro de estoque |
| `PUT` | `/stock/update/{id}` | Atualizar min/max/unidade |
| `GET` | `/stock/{id}` | Consultar estoque |
| `GET` | `/stock/{id}/with-product` | Consultar com dados do produto |
| `GET` | `/stock/product/{product_id}` | Consultar por produto |
| `GET` | `/stock/all` | Listar todos (paginado) |
| `GET` | `/stock/all/with-product` | Listar com produtos |
| `POST` | `/stock/{id}/movement/add` | **Fluxo 1** — Entrada |
| `POST` | `/stock/{id}/movement/remove` | **Fluxo 2** — Saída manual |
| `POST` | `/stock/{id}/movement/adjust` | **Fluxo 3** — Ajuste |
| `GET` | `/stock/{stock_id}/movement` | Histórico de movimentos |
| `GET` | `/stock/low-stock` | Produtos com estoque baixo |
| `GET` | `/stock/out-of-stock` | Produtos sem estoque |
| `GET` | `/stock/report` | Relatório completo |
| `GET` | `/stock/alerts` | Todos os alertas |
| `GET` | `/stock/alerts/expiry` | Alertas de vencimento |
| `POST` | `/stock/alerts/expiry/check?days=30` | **Fluxo 7** — Checar vencimentos |
| `PUT` | `/stock/alerts/{id}/resolve` | Resolver alerta |
| `DELETE` | `/stock/alerts/{id}` | Deletar alerta |

---

## Tipos de Movimento

| Tipo | Origem | Efeito em `current_stock` | Efeito em `reserved_stock` |
|---|---|---|---|
| `in` | Entrada manual / NF | +qty | — |
| `out` | Saída manual / FIFO finalização | -qty | — |
| `adjust_in` | Ajuste inventário (acréscimo) | +diff | — |
| `adjust_out` | Ajuste inventário (redução) | -diff | — |
| `reserve` | AddItem no pedido | -qty | +qty |
| `restore` | Cancelamento de item/pedido | +qty | -qty |
