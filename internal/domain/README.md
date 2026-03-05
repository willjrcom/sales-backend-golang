# `internal/domain`

Contém as entidades e regras puras utilizadas em todo o backend. Cada subpasta possui um README explicando o contexto específico:

| Módulo | O que representa |
|--------|------------------|
| `address/` | Endereços e geocodificação. |
| `advertising/` | Campanhas promocionais. |
| `client/` | Clientes finais e histórico. |
| `company/` | Empresa/tenant, assinatura e billing. |
| `company_category/` | Categorias de empresa e vínculos com patrocinadores. |
| `employee/` | Funcionários e pagamentos. |
| `entity/` | Estruturas base de auditoria. |
| `fiscal_invoice/` | NF-e emitidas. |
| `fiscal_settings/` | Configurações fiscais por empresa. |
| `order/` | Agregado de pedidos, itens e pagamentos. |
| `order_process/` | Workflow de produção/cozinha. |
| `person/` | Dados pessoais reutilizados. |
| `preference/` | Preferências globais. |
| `product/` | Produtos, categorias, tamanhos e regras. |
| `schema/` | Metadados do schema multi-tenant. |
| `shift/` | Turnos operacionais. |
| `sponsor/` | Patrocinadores e incentivos. |
| `stock/` | Controle de estoque (documentação detalhada própria). |
| `table/` | Mesas físicas/virtuais. |

## Boas práticas

- Crie *value objects* específicos em vez de espalhar tipos primitivos.
- Regras de validação devem ser testadas aqui (`*_test.go`).
- Não inclua lógica de persistência ou chamadas HTTP nesta camada.

