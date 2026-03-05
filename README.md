# Sales Backend Go — Guia de Arquitetura

Este repositório concentra o backend multi-tenant do ecossistema *Sales*. Ele foi escrito em Go 1.21+ adotando **Clean Architecture** com camadas bem definidas (Domain → Usecases → Infra). Este README serve como mapa oficial tanto para humanos quanto para agentes de IA navegarem pelo código sem perderem as regras de negócio.

> Referências complementares: `BACKEND_OVERVIEW.md` ( visão funcional ) e `docs/stock-flows.md` (detalhes de estoque).

---

## 1. Visão de alto nível

- **Entradas**: CLI (`main.go`) expõe comandos Cobra:
  - `httpserver` → sobe API REST (Chi) e injetores de dependência.
  - `emailworker` → worker assíncrono consumindo RabbitMQ.
  - `migrate`, `public-migrate`, `migrate-all`, `public-migrate-all` → executam SQL bruto em todos os schemas.
- **Camadas**:
  1. `internal/domain`: entidades + regras puras.
  2. `internal/usecases`: orquestra regras, validações e transações.
  3. `internal/infra`: adaptadores (handlers HTTP, DTOs, repositórios Bun/PostgreSQL, serviços externos, schedulers).
- **Infra compartilhada**: `bootstrap` monta banco, servidor HTTP, RBAC, módulos e injeta dependências. `pkg` abriga utilidades exportáveis.

```
┌────────────┐   HTTP/Webhooks   ┌─────────────┐    DTO ↔ Entity    ┌──────────────┐
│  Handler   │ ─────────────────▶│  Usecase    │────────────────────▶│    Domain    │
└────────────┘                   └─────────────┘                     └──────────────┘
       │                               │                                   │
       │                               ▼                                   ▼
       │                        Repositórios (infra) ───────▶ PostgreSQL / serviços externos
       ▼
RabbitMQ / Impressão / Kafka / S3 / MercadoPago / FocusNFe
```

---

## 2. Mapa de pastas

| Pasta | Conteúdo | Documento |
|-------|----------|-----------|
| `/bootstrap` | Inicialização de DB, servidor, RBAC e módulos | [bootstrap/README.md](bootstrap/README.md) |
| `/cmd` | Comandos Cobra (API, workers, migrations) | [cmd/README.md](cmd/README.md) |
| `/config` | Manifests e configs compartilhadas (Kafka UI, etc.) | [config/README.md](config/README.md) |
| `/docs` | Documentação de regras de negócio complementares | [docs/README.md](docs/README.md) |
| `/internal` | Código da aplicação dividido em Domain, Usecases e Infra | [internal/README.md](internal/README.md) |
| `/migrations` | Scripts SQL versionados por tenant/public | [migrations/README.md](migrations/README.md) |
| `/pkg` | Pacotes reutilizáveis fora do Clean Architecture | [pkg/README.md](pkg/README.md) |
| `/scripts` | Utilitários CLI/Go auxiliares (ex.: checagem de tabelas) | [scripts/README.md](scripts/README.md) |

---

## 3. Fluxo padrão de requisições

1. **Handler** (`internal/infra/handler`) valida entrada, converte para DTO e chama o usecase.
2. **Usecase** aplica validações de negócio, chama domínio (`internal/domain`) e orquestra transações.
3. **Repositórios** (`internal/infra/repository/postgres`) persistem via Bun. Operações multi-tenant usam schemas específicos configurados no contexto da requisição.
4. **Serviços externos** (RabbitMQ, MercadoPago, FocusNFe, S3, POS) são acessados por `internal/infra/service/*`.
5. **Eventos**: ordens relevantes disparam mensagens (RabbitMQ/Kafka), impressões e notificações email.

Erros são capturados por `bootstrap/server` (middlewares) e convertidos em respostas JSON padronizadas. Estoque e pedidos usam bloqueios pessimistas (`SELECT FOR UPDATE`) para evitar race conditions.

---

## 4. Regras de negócio críticas

### Autenticação multi-tenant
- Login em duas etapas: `id_token` (30 min) + `access_token` atrelado ao schema (2 h).
- `bootstrap/server` injeta middlewares (timeout 5s, refresh token, RBAC).
- Cada request precisa do header `X-Company-Schema`; módulos preparam `context.Context` com schema ativo.

### Pedidos + Estoque
- Fluxos descritos em `docs/stock-flows.md` e `internal/domain/stock/README.md`.
- **Reserva automática**: ao criar item/grupo pendente, chama `StockService.Reserve` e grava movimento `RESERVE`.
- **Consumo FIFO** no fechamento do pedido; estoque negativo permitido, mas gera movimento residual para auditoria.
- Cancelamentos (item, grupo, pedido) sempre restauram estoque antes de remover registros, garantindo consistência.

### Pagamentos e checkout
- `internal/usecases/checkout` orquestra métodos (dinheiro, PIX, cartão, delivery) e consulta `company` para regras de repasse.
- Mercadopago e POS têm serviços dedicados (`internal/infra/service/mercadopago`, `.../pos`) que assinam webhooks.

### Relatórios
- `internal/usecases/report` agrega dados brutos em queries SQL customizadas (ver `internal/infra/repository/postgres/report`).
- Dashboards principais: estoque, top produtos, complementos, uso adicional e performance de funcionários.

### Fiscal e patrocinadores
- `internal/usecases/fiscal_invoice` integra com FocusNFe via `internal/infra/service/focusnfe`.
- Categorias de patrocinadores/ads influenciam sugestões de adicionais e pricing por schema.

---

## 5. Execução e operações

| Objetivo | Comando | Observações |
|----------|---------|-------------|
| Subir API HTTP | `go run main.go httpserver --port :8080 --environment dev` | Requer `DATABASE_URL`, `RABBITMQ_URL`, credenciais S3, MercadoPago e FocusNFe. |
| Worker de email | `go run main.go emailworker` | Consome fila `email.send` do RabbitMQ. |
| Aplicar SQL (schema tenant) | `go run main.go migrate --file bootstrap/database/migrations/XXX.sql` | Executa sequencialmente em todos os schemas ativos. |
| Aplicar todas pendentes | `go run main.go migrate-all` | Usa tabela `schema_migrations`. Ignora arquivos removidos. |
| Apenas schema público | `go run main.go public-migrate --file ...` ou `public-migrate-all` | Útil para colunas compartilhadas. |

Durante o desenvolvimento, execute `docker-compose up -d postgres kafka rabbitmq` para subir dependências. O checklist de saúde inclui:
- `/healthz` (HTTP) retornando 200.
- Conexões S3 e RabbitMQ criadas no startup (`cmd/httpserver.go`).
- Scheduler diário (`internal/infra/scheduler`) ativo para alertas de estoque.

---

## 6. Convenções para colaborar (humanos e IAs)

1. **Leia o README da pasta** que pretende alterar antes de codar.
2. **Nunca** acesse o domínio a partir de handlers; trafegue apenas via usecases.
3. Gere DTOs novos em `internal/infra/dto/<contexto>` e mantenha compatibilidade com o frontend.
4. Operações que mexem em estoque/pedidos precisam atualizar os movimentos e alerts — consulte `internal/domain/stock/README.md`.
5. Adicione novas integrações externas em `internal/infra/service/*` e injete via `internal/infra/modules`.
6. Testes de domínio e usecases devem ficar próximos do código (`*_test.go`). Para repositórios use `internal/infra/repository/...`.
7. Documente qualquer regra nova criando um tópico no README da pasta correspondente ou adicionando novos arquivos em `docs/`.

Seguindo este guia, quem chegar (ou qualquer agente automático) encontra rapidamente onde cada regra vive e como expandir o sistema sem quebrar os fluxos atuais.

