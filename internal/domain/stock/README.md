# Documentação do Sistema de Estoque

## Visão Geral

O sistema de estoque gerencia o controle de entrada/saída de produtos usando a estratégia **FIFO** (First-In, First-Out) com suporte a lotes, reservas, alertas e relatórios.

---

## Estrutura de Dados

```
Stock (estoque por produto/variação)
 ├── CurrentStock     — quantidade disponível para reserva/venda
 ├── ReservedStock    — quantidade reservada por pedidos em aberto  ← salvo no BD
 ├── MinStock         — nível mínimo (gera alerta)
 ├── MaxStock         — nível máximo (gera alerta de over-stock)
 ├── IsActive         — controle ativo ou não
 └── Batches[]        — lotes com validade/custo individuais
      ├── InitialQuantity
      ├── CurrentQuantity
      ├── CostPrice
      └── ExpiresAt (opcional)
```

---

## Fluxos Principais

### 1. Adição de Item ao Pedido (`AddItemOrder`)

```
AddItemOrder
  └─► AddItem (insert no BD)  ← PRIMEIRO (evita orphan se insert falhar)
  └─► DebitStockFromItem
        └── ReserveStock: CurrentStock ↓, ReservedStock ↑
        └── Salva movimento tipo: RESERVE
        └── (falha silenciosa — não bloqueia o pedido)
```

> **Decisão:** Se não há controle de estoque para o produto, o item é adicionado sem erro.

---

### 2. Finalização do Pedido (`FinishOrder`)

```
FinishOrder
  └─► DebitStockFromOrder
        └── Para cada item do pedido:
              └─► DebitStockFIFO
                    ├── SELECT FOR UPDATE (lock pessimista nos lotes)
                    ├── Lotes consumidos em FIFO (mais antigo primeiro)
                    ├── CurrentStock ↓, ReservedStock ↓
                    └── Movimentos tipo: OUT (um por lote consumido)
```

> **Estoque negativo:** movimento residual criado sem lote associado para auditoria.

---

### 3. Cancelamento de Item (`DeleteItemOrder`)

```
DeleteItemOrder
  └─► RestoreStockFromItem → ReservedStock ↓, CurrentStock ↑ (tipo: RESTORE)
  └─► DeleteItem
```

---

### 4. Cancelamento de Grupo (`CancelGroupItem`)

```
CancelGroupItem (chamada isolada)
  └─► (se status != Cancelled) → restoreStockFromGroupItem
        └─► RestoreStockFromItem para cada item
  └─► Atualiza status para Cancelled
```

---

### 5. Cancelamento de Pedido (`CancelOrder`)

```
CancelOrder
  └─► restoreStockFromOrder
        ├── Pedido FINALIZADO → RestoreStockFromOrder (via movimentos OUT no BD, restaura lotes)
        └── Pedido em aberto → restoreStockFromGroupItem (revertendo reservas)
  └─► CancelGroupItemSkipStockRestore (cancela grupos SEM restaurar estoque de novo)
```

> **Importante:** `CancelGroupItemSkipStockRestore` existe para evitar dupla restauração.

---

### 6. Entrada Manual (`AddMovementStock`)

```
AddMovementStock (qtd positiva)
  └─► Cria StockBatch + CurrentStock ↑ + movimento tipo: IN
```

### 7. Saída Manual (`RemoveMovementStock`)

```
RemoveMovementStock
  └─► DebitStockFIFO (orderID = uuid.Nil) → NÃO decrementa ReservedStock
```

### 8. Ajuste de Estoque (`AdjustMovementStock`)

```
AdjustMovementStock
  ├── diferença > 0 → AddMovementStock
  └── diferença < 0 → RemoveMovementStock
```

---

## Tipos de Alertas

### Nível de Estoque (`CheckAlerts`)

| Tipo | Condição |
|------|----------|
| `low_stock` | CurrentStock ≤ MinStock |
| `out_of_stock` | CurrentStock ≤ 0 |
| `over_stock` | CurrentStock > MaxStock |

### Vencimento (`CheckExpirations`)

| Tipo | Condição |
|------|----------|
| `near_expiration` | ExpiresAt < (agora + X dias) |
| `expired` | ExpiresAt < agora |

---

## Endpoints

| Método | Rota | Descrição |
|--------|------|-----------|
| `POST` | `/stock/new` | Criar controle |
| `PUT` | `/stock/update/{id}` | Atualizar |
| `GET` | `/stock/{id}` | Buscar por ID |
| `GET` | `/stock/product/{product_id}` | Buscar por produto |
| `GET` | `/stock/all` | Listar todos |
| `POST` | `/stock/{id}/movement/add` | Entrada manual |
| `POST` | `/stock/{id}/movement/remove` | Saída manual (FIFO) |
| `POST` | `/stock/{id}/movement/adjust` | Ajuste |
| `GET` | `/stock/{stock_id}/movement` | Histórico |
| `GET` | `/stock/alerts` | Todos os alertas |
| `GET` | `/stock/alerts/expiry` | Alertas de vencimento |
| `POST` | `/stock/alerts/expiry/check?days=N` | Checar vencimentos |
| `GET` | `/stock/alerts/{id}` | Alerta por ID |
| `PUT` | `/stock/alerts/{id}/resolve` | Resolver alerta |
| `GET` | `/stock/report` | Relatório |
| `GET` | `/stock/low-stock` | Estoque baixo |
| `GET` | `/stock/out-of-stock` | Sem estoque |

---

## Migrações Necessárias

| Arquivo | Descrição |
|---------|-----------|
| `20260302200000_add_reserved_stock_to_stocks.sql` | Coluna `reserved_stock` em `stocks` |
| `20260302210000_make_variation_id_nullable_in_stock.sql` | Remove NOT NULL de `product_variation_id` em `stock_alerts` e `stock_batches` |

---

## Bugs Corrigidos (16 total)

| # | Severidade | Descrição |
|---|---|---|
| 1 | 🔴 | `DebitStockFromItem` chamava FIFO — corrigido para `ReserveStock` |
| 2 | 🔴 | `RestoreStockFromItem` com receiver por valor |
| 3 | 🔴 | Nil pointer em `*ProductVariationID` em `CheckAlerts`/`CheckExpirations` |
| 4 | 🟡 | `DebitStockFIFO` decrementava `ReservedStock` em remoções manuais |
| 5 | 🟡 | Relatório de custo com valor fixo R$10 |
| 6 | 🟡 | `CreateBatch` com transação própria desconectada da transação pai |
| 7 | 🟡 | Lote manual sem `ProductVariationID` |
| 8 | 🔴 | `restoreStockFromGroupItem` ignorava erros silenciosamente |
| 9 | 🔴 | `CancelGroupItem` com condição invertida para Staging |
| 10 | 🔴 | `CancelOrder` com dupla restauração de estoque |
| 11 | 🔴 | Reserva commitada antes do `AddItem` — orphan em falha |
| 12 | 🟡 | Campo `stockService` morto em `ItemService` |
| 13 | 🔴 | **`ReservedStock` nunca persistido no banco de dados** |
| 14 | 🔴 | Rota `/alerts/expiry` conflitava com `/alerts/{id}` no chi router |
| 15 | 🟡 | `ProductVariationID` `notnull` em `stock_alerts` — falha em produtos sem variação |
| 16 | 🟡 | `ProductVariationID` `notnull` em `stock_batches` — idem |
