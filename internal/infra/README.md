# `internal/infra`

Adaptadores que conectam usecases ao mundo externo.

| Subpasta | Papel |
|----------|-------|
| `dto/` | Estruturas de request/response expostas pelos handlers. |
| `handler/` | Rotas HTTP (Chi) que validam entrada e chamam usecases. |
| `modules/` | Funções responsáveis por registrar handlers, repositórios e serviços no `ServerChi`. |
| `repository/` | Implementações Bun/PostgreSQL das interfaces de persistência. |
| `scheduler/` | Jobs assíncronos recorrentes (ex.: alertas de estoque). |
| `service/` | Integrações externas (S3, RabbitMQ, FocusNFe, MercadoPago, POS, etc.). |

Cada pasta possui README descrevendo como estender/modificar. Os componentes aqui nunca devem importar pacotes dos handlers diretamente—apenas interfaces declaradas em usecases.

