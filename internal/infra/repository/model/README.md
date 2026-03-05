# Repository / Model

Structs auxiliares usados pelo Bun para mapear resultados complexos (joins, views, relatórios). Eles não expõem regras de negócio, apenas refletem o formato retornado pelo banco.

## Orientações

- Prefira criar *view models* aqui quando o resultado não corresponder diretamente a uma entidade de domínio.
- Sempre documente campos calculados para facilitar manutenção dos repositórios.
- Use tags `bun:\"\"` para configurar colunas, relations e arrays de forma explícita.

