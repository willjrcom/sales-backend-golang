# Bootstrap / Database

Responsável por expor a conexão Bun com PostgreSQL e centralizar as migrações internas utilizadas pelos comandos `migrate*`.

## Conexão

- `database.NewPostgreSQLConnection()` lê `DATABASE_URL`, configura *connection pool* e aplica `search_path` dinâmico.
- Cada request troca o schema ativo via middleware (`X-Company-Schema`). Para tarefas em lote, os usecases recebem explicitamente o schema que devem usar.
- O pacote também expõe helpers para `context.Context` (armazenando schema) e para abrir transações com *locking* pessimista.

## Migrações

- `bootstrap/database/migrations/` → scripts SQL aplicados em todos os schemas de clientes.
- `bootstrap/database/migrations/public/` → scripts específicos do schema `public`.
- Os comandos `migrate`, `public-migrate`, `migrate-all` e `public-migrate-all` procuram arquivos aqui (vide [cmd/README.md](../../cmd/README.md)).

## Convenções

1. Nomeie os arquivos como `NNN_descricao.sql` para manter ordenação alfabética = ordem cronológica.
2. Scripts precisam ser **idempotentes**: utilize `IF NOT EXISTS` e valide colunas/tabelas antes de alterar.
3. Evite inserir dados sensíveis; use *seeders* específicos em `scripts/` se necessário.
4. Quando surgir uma nova migration, atualize também `scripts/migrations/README.md` para orientar quem executa manualmente.

