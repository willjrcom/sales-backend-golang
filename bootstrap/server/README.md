# Bootstrap / Server

Implementa `ServerChi`, o wrapper responsável por subir o servidor HTTP (Chi) com todos os middlewares obrigatórios.

## Principais componentes

- `server.go`
  - `ServerChi` encapsula o router, lista de handlers registrados e método `StartServer`.
  - `AddHandler` recebe instâncias criadas em `internal/infra/handler`, aplica middlewares e registra as rotas.
  - `StartServer(port string)` inicia o listener HTTP configurando `Read/Write Timeout` e logs.
- `middlewares.go`
  - `AuthMiddleware` valida `Authorization` + `X-Company-Schema`.
  - `RecoverMiddleware` captura panics e retorna JSON consistente (`error-boundary` para o app).
  - `TimeoutMiddleware` fecha requisições lentas (5 s) para proteger o pool.

## Fluxo de montagem

1. `cmd/httpserver.go` instancia `ServerChi`.
2. `internal/infra/modules.MainModules` registra cada handler via `AddHandler`.
3. Middlewares globais são adicionados automaticamente (ordem: recover → timeout → logging → auth → RBAC).
4. `StartServer` inicia o servidor após todas as dependências (DB, RabbitMQ, S3) estarem prontas.

## Boas práticas

- Não registre rotas diretamente em `chi.Router`; sempre passe por `ServerChi` para manter logs + RBAC.
- Para rotas públicas (webhooks), defina-as em `handler.UnprotectedRoutes`.
- Middlewares específicos de módulo devem ser adicionados dentro do handler, nunca aqui.

