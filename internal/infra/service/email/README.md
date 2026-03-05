# Service / Email

Produtor/consumidor RabbitMQ responsável por enfileirar e enviar emails transacionais (convites, recibos).

---

## 1. Responsabilidades
- Publicar mensagens `email.send` contendo template + payload.
- Executar worker (vide `cmd/emailworker.go`) que consome fila, renderiza template e envia via SMTP/serviço externo.
- Fornecer API simples para usecases (ex.: user.resetPassword).

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `Send(ctx, payload EmailPayload) error` | Publica mensagem JSON na fila RabbitMQ. |
| `RunConsumer() error` | Inicia loop bloqueante que consome mensagens e chama provider SMTP. |
| `withTemplate(id string) EmailPayload` | Helper para carregar templates do disco. |

## 3. Fluxo típico
- Usecase chama `Send` passando template e parâmetros.
- Mensagem vai para RabbitMQ (`queue=email.send`).
- Worker `cmd/emailworker` roda `RunConsumer`, desserializa payload e envia via SMTP.
- Registro de envio/log é escrito para observabilidade.

## 4. Configuração / Env Vars
- `RABBITMQ_URL`
- `EMAIL_PROVIDER_API_KEY`
- `EMAIL_FROM`

## 5. Exemplo de uso
```go
go
payload := emailservice.EmailPayload{
    Template: "user_invite",
    To:       []string{"admin@loja.com"},
    Data: map[string]any{"invite_link": link},
}
if err := emailService.Send(ctx, payload); err != nil {
    logger.Error("email send failed", err)
}
```

## 6. Falhas comuns
- ErrQueueUnavailable quando RabbitMQ não responde.
- ProviderError quando SMTP/API externo retorna falha.

## 7. Notas operacionais
- Worker deve ser executado separadamente (`go run main.go emailworker`).
- Templates residem em `internal/infra/service/email/templates`.
