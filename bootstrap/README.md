# Bootstrap

Esta pasta reúne tudo que é necessário para iniciar o backend antes que o código de negócio seja executado. Os arquivos daqui são chamados pelos comandos definidos em `cmd/` e garantem que banco, servidor HTTP e dependências externas estejam configurados.

## Estrutura

| Subpasta | Função | Documento |
|----------|--------|-----------|
| `database/` | Conexão com PostgreSQL multi-tenant e diretório oficial de migrações internas. | [database/README.md](database/README.md) |
| `handler/` | Registro dos *handlers* HTTP (Chi) e middlewares globais. | [handler/README.md](handler/README.md) |
| `rbac/` | Configuração do controle de acesso baseado em papéis. | [rbac/README.md](rbac/README.md) |
| `server/` | Implementação do `ServerChi`, middlewares padrão e inicialização do router. | [server/README.md](server/README.md) |

## Ordem de inicialização

1. `database.NewPostgreSQLConnection()` abre a conexão Bun e injeta o schema padrão (`public`).
2. `server.NewServerChi()` instancia o roteador com os middlewares de tracing, timeout, auth e logging.
3. `internal/infra/modules.MainModules` injeta repositórios, serviços e handlers no servidor.
4. `cmd/httpserver.go` finaliza registrando health-checks e iniciando o listener HTTP.

## Boas práticas

- Toda nova dependência externa (ex.: serviço de email) deve ser inicializada no `cmd` correspondente e passada para `modules`.
- Middlewares globais pertencem a `bootstrap/server`, não aos handlers individuais.
- Migrações específicas por tenant devem residir em `bootstrap/database/migrations`; scripts auxiliares podem ficar em `scripts/migrations`.

