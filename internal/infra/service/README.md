# Infra / Service

Integrações externas e utilitários (S3, JWT, email, POS, etc.). Cada subpasta implementa um cliente específico, expondo interfaces utilizadas pelos usecases.

## Exemplos

| Serviço | Descrição |
|---------|-----------|
| `bcrypt/` | Hash/salt de senhas. |
| `cnpj/` | Validação de CNPJ. |
| `email/` | Produtor/consumidor RabbitMQ para envio de emails. |
| `focusnfe/` | Cliente da API FocusNFe (NF-e). |
| `geocode/` | Consulta coordenadas para delivery. |
| `jwt/` | Emissão/validação de tokens multi-tenant. |
| `mercadopago/` | Pagamentos e webhooks. |
| `pos/` | Integração com terminais locais. |
| `s3/` | Upload/download em buckets privados. |

## Convenções

1. Clientes devem aceitar `context.Context` e respeitar timeouts curtos.
2. Retorne erros ricos (código/mensagem) para o usecase decidir a ação.
3. Documente endpoints externos e parâmetros sensíveis nos READMEs das subpastas.

