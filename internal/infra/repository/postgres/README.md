# Repository / Postgres

Implementações Bun de cada módulo do sistema. Cada subpasta reflete um contexto (order, stock, report, etc.) e expõe métodos utilizados pelos usecases.

| Subpasta | Responsabilidade principal |
|----------|---------------------------|
| `address/` | Persistência de endereços e geolocalização. |
| `client/` | CRUD de clientes e histórico. |
| `company/` | Empresas, assinaturas, billing. |
| `item/` | Itens do pedido e adicionais. |
| `order/` | Pedidos e relacionamentos com entregas, mesas, pagamentos. |
| `order_process/` | Filas de produção e analytics. |
| `product/` | Produtos, variações e tamanhos. |
| `report/` | Consultas agregadas (dashboards). |
| `stock/` | Estoque, lotes e movimentos FIFO. |
| ... | (todos os demais módulos seguem o mesmo padrão). |

## Convenções

1. Sempre utilize `ctx` com schema correto (provido pelo middleware).
2. Nomeie métodos seguindo a ação + contexto (`GetOrderByID`, `ListLowStock`).
3. Queries complexas devem ser explicadas no README da subpasta correspondente.
4. Evite lógica de negócio; mantenha apenas consultas e mapeamentos.

