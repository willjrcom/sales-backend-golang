# Comandos (`cmd/`)

Os comandos expostos via Cobra (`go run main.go <command>`) controlam como o backend é executado. Use este README como guia rápido.

| Comando | Arquivo | Objetivo | Flags principais |
|---------|---------|----------|------------------|
| `httpserver` | `httpserver.go` | Sobe a API REST com Chi, RabbitMQ, S3 e módulos injetados. | `--port` (default `:8080`), `--environment` (`dev`, `staging`, `prod`). |
| `emailworker` | `emailworker.go` | Inicia o consumidor de emails que lê RabbitMQ (`RABBITMQ_URL`). | Nenhuma; depende das variáveis de ambiente de email. |
| `migrate` | `migrate.go` | Executa um arquivo SQL específico em **todos** os schemas de tenant. | `--file` caminho relativo (default `bootstrap/database/migrations`). |
| `public-migrate` | `migrate.go` | Executa um arquivo SQL apenas no schema `public`. | `--file`. |
| `migrate-all` | `migrate.go` | Aplica todas as migrações pendentes (arquivos existentes) em cada schema. | Sem flags. |
| `public-migrate-all` | `migrate.go` | Versão `all` exclusiva para schema `public`. | Sem flags. |

## Dicas operacionais

1. Os comandos de migration criam/consultam a tabela `schema_migrations` automaticamente; não remova manualmente.
2. `httpserver` precisa de `DATABASE_URL`, `RABBITMQ_URL`, credenciais S3 e chaves das integrações externas.
3. Para homologação, execute `go run main.go httpserver --port :8081 --environment staging`.

