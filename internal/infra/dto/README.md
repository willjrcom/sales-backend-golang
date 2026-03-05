# Infra / DTO

Estruturas de dados expostas pelas APIs e consumidas pelo frontend/app. As subpastas agrupam DTOs por contexto (order, stock, report, etc.).

## Convenções

- DTOs nunca importam usecases; eles apenas descrevem payloads.
- Conversões entre DTO ↔ domain acontecem nos handlers antes de chamar o usecase.
- Cada diretório de contexto (ex.: `order/`, `stock/`) deve conter structs e helpers específicos, além de documentação adicional caso o contrato seja complexo.

| Subpasta exemplo | Uso |
|------------------|-----|
| `order/` | Payloads para criação e atualização de pedidos. |
| `stock/` | Requests de movimentação e respostas de alertas. |
| `report/` | Filtros e respostas de dashboards. |

Ao adicionar um novo endpoint, mantenha o contrato neste diretório para facilitar evolução e versionamento.

