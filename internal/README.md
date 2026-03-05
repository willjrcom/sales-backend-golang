# Camada `internal/`

Abriga todo o código que não deve ser exposto como pacote público. É dividida em três fatias principais que seguem Clean Architecture:

| Pasta | Função | Observações |
|-------|--------|-------------|
| `domain/` | Entidades ricas, regras puras e enums compartilhados. | Não acessa banco ou libs externas além de utilidades determinísticas. |
| `usecases/` | Orquestração de regras de negócio com dependências explícitas. | Só conversa com `domain` e interfaces de repositórios/serviços. |
| `infra/` | Adapters (HTTP handlers, repositórios Bun, DTOs, serviços externos). | Implementa as interfaces utilizadas pelos usecases. |

## Convenção geral

1. **Domain** nunca deve importar `infra` ou `usecases`.
2. **Usecases** recebem dependências através de construtores/métodos `AddDependencies`.
3. **Infra** é responsável por converter requisições HTTP → DTO → Domain e persistir via Bun.
4. Todo novo módulo precisa ter README próprio descrevendo as regras de negócio (veja subpastas).

Consulte os READMEs específicos em `domain/`, `usecases/` e `infra/` para detalhes e listas de módulos.

