# Infra / Repository

Implementações das interfaces de persistência utilizadas pelos usecases. Existem dois níveis:

- `postgres/` — Repositórios Bun/PostgreSQL (multi-tenant via `SetSearchPath`).
- `local/` e `model/` — utilitários auxiliares (fakes, structs mapeadas).

Cada subpasta de `postgres/` espelha um módulo de negócio (`order`, `stock`, `report`, etc.) e possui arquivos com consultas específicas.

## Convenções

1. Sempre receba `context.Context` como primeiro parâmetro para garantir o schema correto.
2. Use `bun.Tx` quando operações precisarem de transação; nunca gerencie transações aqui sem coordenação com o usecase.
3. Exporte apenas interfaces necessárias; mantenha structs e consultas encapsuladas.
4. Atualize o README da subpasta ao adicionar consultas complexas (ex.: relatórios customizados).

