# Bootstrap / Handler

O pacote `bootstrap/handler` define o contrato mínimo para registrar módulos HTTP no `ServerChi`. Cada handler encapsula:

- `Path`: prefixo base (ex.: `/order`, `/stock`).
- `Handler`: implementação `http.Handler` (geralmente um roteador Chi com middlewares específicos).
- `UnprotectedRoutes`: lista de rotas que ignoram autenticação (ex.: webhooks públicos).

## Como usar

1. Crie um novo handler em `internal/infra/handler/<contexto>.go`, retornando `handler.NewHandler`.
2. No `internal/infra/modules` correspondente, invoque `chi.AddHandler(...)` para que o servidor registre a rota base.
3. Para endpoints públicos, indique o caminho exato em `UnprotectedRoutes` (aceita curingas globais).

Manter este contrato desacoplado facilita testar módulos individualmente e garante que todas as rotas sigam o mesmo padrão de montagem.

